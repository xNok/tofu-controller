package web

import (
	"encoding/json"
	"net/http"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReportSummary contains aggregate counts of Terraform resource statuses.
type ReportSummary struct {
	Total           int `json:"total"`
	Ready           int `json:"ready"`
	Failed          int `json:"failed"`
	Suspended       int `json:"suspended"`
	PendingApproval int `json:"pendingApproval"`
	Progressing     int `json:"progressing"`
}

// ReportHandler returns a summary of all Terraform resources.
// GET /api/v1/report
func (h *Handler) ReportHandler(w http.ResponseWriter, r *http.Request) {
	var list infrav1.TerraformList
	if err := h.kubeClient.List(r.Context(), &list, &client.ListOptions{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	summary := buildReportSummary(list.Items)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(summary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func buildReportSummary(items []infrav1.Terraform) ReportSummary {
	s := ReportSummary{Total: len(items)}
	for i := range items {
		tf := &items[i]
		if tf.Spec.Suspend {
			s.Suspended++
			continue
		}
		if tf.Status.Plan.Pending != "" {
			s.PendingApproval++
			continue
		}
		switch readyConditionStatus(tf) {
		case metav1.ConditionTrue:
			s.Ready++
		case metav1.ConditionFalse:
			s.Failed++
		default:
			s.Progressing++
		}
	}
	return s
}

// readyConditionStatus returns the status of the Ready condition, or "" if not found.
func readyConditionStatus(tf *infrav1.Terraform) metav1.ConditionStatus {
	for _, c := range tf.Status.Conditions {
		if c.Type == "Ready" {
			return metav1.ConditionStatus(c.Status)
		}
	}
	return ""
}
