package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	tfplan "github.com/flux-iac/tofu-controller/api/plan"
	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/flux-iac/tofu-controller/internal/web/kubeclient"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PlanResponse is the response body for GET /api/v1/plan.
type PlanResponse struct {
	PlanID               string `json:"planId"`
	Content              string `json:"content"`
	IsDestroyPlan        bool   `json:"isDestroyPlan"`
	IsDriftDetectionPlan bool   `json:"isDriftDetectionPlan"`
}

// PlanHandler returns the decoded plan content for a Terraform resource.
// GET /api/v1/plan?namespace=...&name=...
func (h *Handler) PlanHandler(w http.ResponseWriter, r *http.Request) {
	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" || name == "" {
		http.Error(w, "namespace and name query parameters are required", http.StatusBadRequest)
		return
	}

	var tf infrav1.Terraform
	if err := h.kubeClient.Get(r.Context(), client.ObjectKey{Namespace: namespace, Name: name}, &tf); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if tf.Status.Plan.Pending == "" {
		http.Error(w, "no pending plan", http.StatusNotFound)
		return
	}

	content, err := fetchPlanContent(r.Context(), h.kubeClient, &tf)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve plan: %s", err), http.StatusInternalServerError)
		return
	}

	resp := PlanResponse{
		PlanID:               tf.Status.Plan.Pending,
		Content:              content,
		IsDestroyPlan:        tf.Status.Plan.IsDestroyPlan,
		IsDriftDetectionPlan: tf.Status.Plan.IsDriftDetectionPlan,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fetchPlanContent retrieves plan output from ConfigMaps, falling back to Secrets.
// Mirrors the pattern from internal/informer/branch-planner/informer.go.
func fetchPlanContent(ctx context.Context, kc *kubeclient.Client, tf *infrav1.Terraform) (string, error) {
	labels := map[string]string{
		tfplan.TFPlanNameLabel:      tfplan.SafeLabelValue(tf.Name),
		tfplan.TFPlanWorkspaceLabel: tfplan.SafeLabelValue(tf.WorkspaceName()),
	}
	listOpts := []client.ListOption{
		client.InNamespace(tf.Namespace),
		client.MatchingLabels(labels),
	}

	var configMaps corev1.ConfigMapList
	if err := kc.List(ctx, &configMaps, listOpts...); err != nil {
		return "", fmt.Errorf("unable to list plan configmaps: %w", err)
	}

	if len(configMaps.Items) > 0 {
		p, err := tfplan.NewFromConfigMaps(tf.Name, tf.Namespace, string(tf.UID), configMaps.Items)
		if err != nil {
			return "", fmt.Errorf("unable to reconstruct plan from configmaps: %w", err)
		}
		return p.ToString(), nil
	}

	// Fallback: try Secrets (older controller versions stored plans in Secrets).
	var secrets corev1.SecretList
	if err := kc.List(ctx, &secrets, listOpts...); err != nil {
		return "", fmt.Errorf("unable to list plan secrets: %w", err)
	}

	if len(secrets.Items) > 0 {
		p, err := tfplan.NewFromSecrets(tf.Name, tf.Namespace, string(tf.UID), secrets.Items)
		if err != nil {
			return "", fmt.Errorf("unable to reconstruct plan from secrets: %w", err)
		}
		return p.ToString(), nil
	}

	return "Please set `spec.storeReadablePlan: human` to view the plan", nil
}
