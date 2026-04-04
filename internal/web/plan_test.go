package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tfplan "github.com/flux-iac/tofu-controller/api/plan"
	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func makePlanConfigMap(tfName, namespace, workspace, planID, content string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tfplan-" + workspace + "-" + tfName,
			Namespace: namespace,
			Labels: map[string]string{
				tfplan.TFPlanNameLabel:      tfplan.SafeLabelValue(tfName),
				tfplan.TFPlanWorkspaceLabel: tfplan.SafeLabelValue(workspace),
			},
			Annotations: map[string]string{
				tfplan.TFPlanFullWorkspaceAnnotation: workspace,
				tfplan.TFPlanSavedAnnotation:         planID,
			},
		},
		Data: map[string]string{
			"tfplan": content,
		},
	}
}

func TestPlanHandler_NoPendingPlan(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "OCIRepository")
	h := newTestHandler(t, tf)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan?namespace=default&name=my-app", nil)
	rec := httptest.NewRecorder()
	h.PlanHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404 when no pending plan, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPlanHandler_MissingParams(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan", nil)
	rec := httptest.NewRecorder()
	h.PlanHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestPlanHandler_WithConfigMap(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "plan-abc", metav1.ConditionUnknown, "OCIRepository")
	tf.Status.Plan.IsDestroyPlan = true
	cm := makePlanConfigMap("my-app", "default", "default", "plan-abc", "terraform plan output here")

	objs := []client.Object{tf, cm}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan?namespace=default&name=my-app", nil)
	rec := httptest.NewRecorder()
	h.PlanHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp PlanResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.PlanID != "plan-abc" {
		t.Errorf("want planId=plan-abc, got %s", resp.PlanID)
	}
	if resp.Content != "terraform plan output here" {
		t.Errorf("unexpected content: %q", resp.Content)
	}
	if !resp.IsDestroyPlan {
		t.Error("want isDestroyPlan=true")
	}
}

func TestPlanHandler_NoConfigMapFallbackMessage(t *testing.T) {
	tf := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Name: "my-app", Namespace: "default"},
		Spec: infrav1.TerraformSpec{
			SourceRef: infrav1.CrossNamespaceSourceReference{Kind: "OCIRepository", Name: "src"},
		},
	}
	tf.Status.Plan.Pending = "plan-xyz"
	h := newTestHandler(t, tf)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/plan?namespace=default&name=my-app", nil)
	rec := httptest.NewRecorder()
	h.PlanHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var resp PlanResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Content == "" {
		t.Error("expected fallback message, got empty content")
	}
}
