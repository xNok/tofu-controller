package web

import (
	"encoding/json"
	"net/http"
	"time"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/fluxcd/pkg/apis/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ActionRequest is the request body for POST /api/v1/resource/action.
type ActionRequest struct {
	Action    string `json:"action"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// ActionHandler performs an action on a Terraform resource.
// POST /api/v1/resource/action
// Supported actions: approve, suspend, resume, force-reconcile
func (h *Handler) ActionHandler(w http.ResponseWriter, r *http.Request) {
	var req ActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Namespace == "" || req.Name == "" {
		http.Error(w, "namespace and name are required", http.StatusBadRequest)
		return
	}

	var tf infrav1.Terraform
	key := client.ObjectKey{Namespace: req.Namespace, Name: req.Name}
	if err := h.kubeClient.Get(r.Context(), key, &tf); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	patch := client.MergeFrom(tf.DeepCopy())

	switch req.Action {
	case "approve":
		if tf.Status.Plan.Pending == "" {
			http.Error(w, "no pending plan to approve", http.StatusBadRequest)
			return
		}
		tf.Spec.ApprovePlan = tf.Status.Plan.Pending

	case "suspend":
		tf.Spec.Suspend = true

	case "resume":
		tf.Spec.Suspend = false

	case "force-reconcile":
		ann := tf.GetAnnotations()
		if ann == nil {
			ann = make(map[string]string)
		}
		ann[meta.ReconcileRequestAnnotation] = time.Now().Format(time.RFC3339Nano)
		tf.SetAnnotations(ann)

	default:
		http.Error(w, "unknown action: "+req.Action, http.StatusBadRequest)
		return
	}

	if err := h.kubeClient.Patch(r.Context(), &tf, patch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
