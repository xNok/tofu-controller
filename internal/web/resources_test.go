package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func makeTerraformFull(name, namespace string, suspend bool, pendingPlan string, readyStatus metav1.ConditionStatus, sourceKind string) *infrav1.Terraform {
	tf := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec: infrav1.TerraformSpec{
			Suspend: suspend,
			SourceRef: infrav1.CrossNamespaceSourceReference{
				Kind: sourceKind,
				Name: "my-source",
			},
		},
	}
	tf.Status.Plan.Pending = pendingPlan
	if readyStatus != "" {
		tf.Status.Conditions = []metav1.Condition{
			{
				Type:               "Ready",
				Status:             readyStatus,
				Reason:             "Test",
				Message:            fmt.Sprintf("status=%s", readyStatus),
				LastTransitionTime: metav1.Now(),
			},
		}
	}
	return tf
}

func TestResourcesHandler_NoFilter(t *testing.T) {
	objs := []client.Object{
		makeTerraformFull("app1", "default", false, "", metav1.ConditionTrue, "GitRepository"),
		makeTerraformFull("app2", "staging", false, "plan-123", metav1.ConditionUnknown, "OCIRepository"),
		makeTerraformFull("app3", "default", true, "", metav1.ConditionTrue, "Bucket"),
	}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources", nil)
	rec := httptest.NewRecorder()
	h.ResourcesHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var results []ResourceSummary
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 3 {
		t.Errorf("want 3 results, got %d", len(results))
	}
}

func TestResourcesHandler_FilterNamespace(t *testing.T) {
	objs := []client.Object{
		makeTerraformFull("app1", "default", false, "", metav1.ConditionTrue, "GitRepository"),
		makeTerraformFull("app2", "staging", false, "", metav1.ConditionTrue, "GitRepository"),
	}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources?namespace=default", nil)
	rec := httptest.NewRecorder()
	h.ResourcesHandler(rec, req)

	var results []ResourceSummary
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("want 1 result, got %d", len(results))
	}
	if results[0].Name != "app1" {
		t.Errorf("want app1, got %s", results[0].Name)
	}
}

func TestResourcesHandler_FilterName(t *testing.T) {
	objs := []client.Object{
		makeTerraformFull("my-app", "default", false, "", metav1.ConditionTrue, "GitRepository"),
		makeTerraformFull("other", "default", false, "", metav1.ConditionTrue, "GitRepository"),
	}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources?name=my", nil)
	rec := httptest.NewRecorder()
	h.ResourcesHandler(rec, req)

	var results []ResourceSummary
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Errorf("want 1 result, got %d", len(results))
	}
}

func TestResourcesHandler_FilterStatus(t *testing.T) {
	objs := []client.Object{
		makeTerraformFull("ready1", "default", false, "", metav1.ConditionTrue, "GitRepository"),
		makeTerraformFull("failed1", "default", false, "", metav1.ConditionFalse, "GitRepository"),
		makeTerraformFull("pending1", "default", false, "plan-xyz", metav1.ConditionUnknown, "GitRepository"),
	}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources?status=PendingApproval", nil)
	rec := httptest.NewRecorder()
	h.ResourcesHandler(rec, req)

	var results []ResourceSummary
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Name != "pending1" {
		t.Errorf("want [pending1], got %v", results)
	}
}

func TestResourcesHandler_FilterSourceKind(t *testing.T) {
	objs := []client.Object{
		makeTerraformFull("app1", "default", false, "", metav1.ConditionTrue, "GitRepository"),
		makeTerraformFull("app2", "default", false, "", metav1.ConditionTrue, "OCIRepository"),
	}
	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources?sourceKind=OCIRepository", nil)
	rec := httptest.NewRecorder()
	h.ResourcesHandler(rec, req)

	var results []ResourceSummary
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].SourceKind != "OCIRepository" {
		t.Errorf("want 1 OCIRepository result, got %v", results)
	}
}

func TestResourceStatus(t *testing.T) {
	tests := []struct {
		name    string
		tf      *infrav1.Terraform
		want    string
	}{
		{
			name: "suspended",
			tf:   makeTerraformFull("x", "ns", true, "", metav1.ConditionTrue, "GitRepository"),
			want: "Suspended",
		},
		{
			name: "pendingApproval",
			tf:   makeTerraformFull("x", "ns", false, "plan-1", metav1.ConditionUnknown, "GitRepository"),
			want: "PendingApproval",
		},
		{
			name: "ready",
			tf:   makeTerraformFull("x", "ns", false, "", metav1.ConditionTrue, "GitRepository"),
			want: "Ready",
		},
		{
			name: "failed",
			tf:   makeTerraformFull("x", "ns", false, "", metav1.ConditionFalse, "GitRepository"),
			want: "Failed",
		},
		{
			name: "progressing",
			tf:   makeTerraformFull("x", "ns", false, "", metav1.ConditionUnknown, "GitRepository"),
			want: "Progressing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resourceStatus(tt.tf)
			if got != tt.want {
				t.Errorf("want %q, got %q", tt.want, got)
			}
		})
	}
}
