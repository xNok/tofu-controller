package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestResourceHandler_Success(t *testing.T) {
	obj := makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "OCIRepository")
	h := newTestHandler(t, obj)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resource?namespace=default&name=my-app", nil)
	rec := httptest.NewRecorder()
	h.ResourceHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var detail ResourceDetail
	if err := json.NewDecoder(rec.Body).Decode(&detail); err != nil {
		t.Fatal(err)
	}
	if detail.Name != "my-app" {
		t.Errorf("want name=my-app, got %s", detail.Name)
	}
	if detail.Namespace != "default" {
		t.Errorf("want namespace=default, got %s", detail.Namespace)
	}
	if detail.Status != "Ready" {
		t.Errorf("want status=Ready, got %s", detail.Status)
	}
	if detail.SourceKind != "OCIRepository" {
		t.Errorf("want sourceKind=OCIRepository, got %s", detail.SourceKind)
	}
}

func TestResourceHandler_MissingParams(t *testing.T) {
	h := newTestHandler(t)

	tests := []struct {
		url  string
		desc string
	}{
		{"/api/v1/resource", "no params"},
		{"/api/v1/resource?namespace=default", "missing name"},
		{"/api/v1/resource?name=foo", "missing namespace"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			h.ResourceHandler(rec, req)
			if rec.Code != http.StatusBadRequest {
				t.Errorf("want 400, got %d", rec.Code)
			}
		})
	}
}

func TestResourceHandler_NotFound(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resource?namespace=default&name=missing", nil)
	rec := httptest.NewRecorder()
	h.ResourceHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", rec.Code)
	}
}

func TestResourceHandler_IncludesSpecAndStatus(t *testing.T) {
	tf := makeTerraformFull("infra", "prod", false, "plan-001", metav1.ConditionUnknown, "GitRepository")
	h := newTestHandler(t, (client.Object)(tf))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resource?namespace=prod&name=infra", nil)
	rec := httptest.NewRecorder()
	h.ResourceHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var detail ResourceDetail
	if err := json.NewDecoder(rec.Body).Decode(&detail); err != nil {
		t.Fatal(err)
	}
	if detail.PendingPlan != "plan-001" {
		t.Errorf("want pendingPlan=plan-001, got %q", detail.PendingPlan)
	}
}
