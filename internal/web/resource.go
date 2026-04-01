package web

import (
	"encoding/json"
	"net/http"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceDetail is the full response for GET /api/v1/resource.
type ResourceDetail struct {
	ResourceSummary
	Spec          infrav1.TerraformSpec   `json:"spec"`
	TerraformStatus infrav1.TerraformStatus `json:"terraformStatus"`
}

// ResourceHandler returns a single Terraform resource by namespace and name.
// GET /api/v1/resource?namespace=...&name=...
func (h *Handler) ResourceHandler(w http.ResponseWriter, r *http.Request) {
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

	detail := ResourceDetail{
		ResourceSummary: toResourceSummary(&tf),
		Spec:            tf.Spec,
		TerraformStatus: tf.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(detail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
