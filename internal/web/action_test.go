package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func postAction(t *testing.T, h *Handler, req ActionRequest) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(req)
	r := httptest.NewRequest(http.MethodPost, "/api/v1/resource/action", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ActionHandler(w, r)
	return w
}

func TestActionHandler_Approve(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "plan-abc", metav1.ConditionUnknown, "OCIRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{
		Action:    "approve",
		Namespace: "default",
		Name:      "my-app",
	})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d: %s", rec.Code, rec.Body.String())
	}

	// Verify spec.approvePlan was set
	var updated infrav1.Terraform
	if err := h.kubeClient.Get(t.Context(), client.ObjectKey{Namespace: "default", Name: "my-app"}, &updated); err != nil {
		t.Fatal(err)
	}
	if updated.Spec.ApprovePlan != "plan-abc" {
		t.Errorf("want approvePlan=plan-abc, got %q", updated.Spec.ApprovePlan)
	}
}

func TestActionHandler_ApproveMissingPlan(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "OCIRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{
		Action:    "approve",
		Namespace: "default",
		Name:      "my-app",
	})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestActionHandler_Suspend(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "GitRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{Action: "suspend", Namespace: "default", Name: "my-app"})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated infrav1.Terraform
	if err := h.kubeClient.Get(t.Context(), client.ObjectKey{Namespace: "default", Name: "my-app"}, &updated); err != nil {
		t.Fatal(err)
	}
	if !updated.Spec.Suspend {
		t.Error("want suspend=true")
	}
}

func TestActionHandler_Resume(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", true, "", metav1.ConditionTrue, "GitRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{Action: "resume", Namespace: "default", Name: "my-app"})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated infrav1.Terraform
	if err := h.kubeClient.Get(t.Context(), client.ObjectKey{Namespace: "default", Name: "my-app"}, &updated); err != nil {
		t.Fatal(err)
	}
	if updated.Spec.Suspend {
		t.Error("want suspend=false")
	}
}

func TestActionHandler_ForceReconcile(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "OCIRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{Action: "force-reconcile", Namespace: "default", Name: "my-app"})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d: %s", rec.Code, rec.Body.String())
	}

	var updated infrav1.Terraform
	if err := h.kubeClient.Get(t.Context(), client.ObjectKey{Namespace: "default", Name: "my-app"}, &updated); err != nil {
		t.Fatal(err)
	}
	if updated.GetAnnotations()[meta.ReconcileRequestAnnotation] == "" {
		t.Error("expected ReconcileRequestAnnotation to be set")
	}
}

func TestActionHandler_UnknownAction(t *testing.T) {
	tf := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "GitRepository")
	h := newTestHandler(t, (client.Object)(tf))

	rec := postAction(t, h, ActionRequest{Action: "delete-everything", Namespace: "default", Name: "my-app"})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

func TestActionHandler_MissingFields(t *testing.T) {
	h := newTestHandler(t)

	tests := []ActionRequest{
		{Action: "suspend"},
		{Action: "suspend", Namespace: "default"},
		{Action: "suspend", Name: "foo"},
	}
	for _, req := range tests {
		rec := postAction(t, h, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("want 400 for %+v, got %d", req, rec.Code)
		}
	}
}
