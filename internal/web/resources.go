package web

import (
	"encoding/json"
	"net/http"
	"strings"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceSummary is a trimmed-down view of a Terraform resource for list display.
type ResourceSummary struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Suspended   bool   `json:"suspended"`
	SourceKind  string `json:"sourceKind"`
	SourceName  string `json:"sourceName"`
	PendingPlan string `json:"pendingPlan,omitempty"`
	Message     string `json:"message,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// ResourcesHandler returns a filtered list of Terraform resources.
// GET /api/v1/resources?namespace=...&name=...&status=...&sourceKind=...
func (h *Handler) ResourcesHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filterNamespace := q.Get("namespace")
	filterName := q.Get("name")
	filterStatus := q.Get("status")
	filterSourceKind := q.Get("sourceKind")

	listOpts := &client.ListOptions{}
	if filterNamespace != "" {
		listOpts.Namespace = filterNamespace
	}

	var list infrav1.TerraformList
	if err := h.kubeClient.List(r.Context(), &list, listOpts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := make([]ResourceSummary, 0, len(list.Items))
	for i := range list.Items {
		tf := &list.Items[i]
		summary := toResourceSummary(tf)

		if filterName != "" && !strings.Contains(strings.ToLower(tf.Name), strings.ToLower(filterName)) {
			continue
		}
		if filterStatus != "" && summary.Status != filterStatus {
			continue
		}
		if filterSourceKind != "" && tf.Spec.SourceRef.Kind != filterSourceKind {
			continue
		}

		results = append(results, summary)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func toResourceSummary(tf *infrav1.Terraform) ResourceSummary {
	s := ResourceSummary{
		Namespace:  tf.Namespace,
		Name:       tf.Name,
		Suspended:  tf.Spec.Suspend,
		SourceKind: tf.Spec.SourceRef.Kind,
		SourceName: tf.Spec.SourceRef.Name,
	}
	s.Status = resourceStatus(tf)
	if tf.Status.Plan.Pending != "" {
		s.PendingPlan = tf.Status.Plan.Pending
	}
	for _, c := range tf.Status.Conditions {
		if c.Type == "Ready" {
			s.Message = c.Message
			s.LastUpdated = c.LastTransitionTime.UTC().Format("2006-01-02T15:04:05Z")
			break
		}
	}
	return s
}

// resourceStatus maps Terraform state to a UI status string.
func resourceStatus(tf *infrav1.Terraform) string {
	if tf.Spec.Suspend {
		return "Suspended"
	}
	if tf.Status.Plan.Pending != "" {
		return "PendingApproval"
	}
	switch readyConditionStatus(tf) {
	case metav1.ConditionTrue:
		return "Ready"
	case metav1.ConditionFalse:
		return "Failed"
	default:
		// Unknown or no condition yet — look for more specific info
		for _, c := range tf.Status.Conditions {
			if c.Type == "Apply" && c.Status == metav1.ConditionStatus(metav1.ConditionUnknown) {
				return "Applying"
			}
			if c.Type == "Plan" && c.Status == metav1.ConditionStatus(metav1.ConditionUnknown) {
				return "Planning"
			}
		}
		return "Progressing"
	}
}
