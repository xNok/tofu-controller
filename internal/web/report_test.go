package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/flux-iac/tofu-controller/internal/web/kubeclient"
)

func newTestHandler(t *testing.T, objs ...client.Object) *Handler {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	if err := infrav1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithStatusSubresource(&infrav1.Terraform{}, &corev1.ConfigMap{}, &corev1.Secret{}).
		Build()
	kc := kubeclient.New(fakeClient, fakeClient, scheme)
	return &Handler{kubeClient: kc}
}

func makeTerraform(name, namespace string, suspend bool, pendingPlan string, readyStatus metav1.ConditionStatus) *infrav1.Terraform {
	tf := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
		Spec:       infrav1.TerraformSpec{Suspend: suspend},
	}
	tf.Status.Plan.Pending = pendingPlan
	if readyStatus != "" {
		tf.Status.Conditions = []metav1.Condition{
			{Type: "Ready", Status: readyStatus, Reason: "Test", LastTransitionTime: metav1.Now()},
		}
	}
	return tf
}

func TestReportHandler(t *testing.T) {
	objs := []client.Object{
		makeTerraform("ready1", "default", false, "", metav1.ConditionTrue),
		makeTerraform("ready2", "default", false, "", metav1.ConditionTrue),
		makeTerraform("failed1", "default", false, "", metav1.ConditionFalse),
		makeTerraform("suspended1", "default", true, "", metav1.ConditionTrue),
		makeTerraform("pending1", "default", false, "plan-abc123", metav1.ConditionUnknown),
		makeTerraform("progressing1", "default", false, "", metav1.ConditionUnknown),
	}

	h := newTestHandler(t, objs...)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	rec := httptest.NewRecorder()
	h.ReportHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var summary ReportSummary
	if err := json.NewDecoder(rec.Body).Decode(&summary); err != nil {
		t.Fatal(err)
	}

	if summary.Total != 6 {
		t.Errorf("total: want 6, got %d", summary.Total)
	}
	if summary.Ready != 2 {
		t.Errorf("ready: want 2, got %d", summary.Ready)
	}
	if summary.Failed != 1 {
		t.Errorf("failed: want 1, got %d", summary.Failed)
	}
	if summary.Suspended != 1 {
		t.Errorf("suspended: want 1, got %d", summary.Suspended)
	}
	if summary.PendingApproval != 1 {
		t.Errorf("pendingApproval: want 1, got %d", summary.PendingApproval)
	}
	if summary.Progressing != 1 {
		t.Errorf("progressing: want 1, got %d", summary.Progressing)
	}
}

func TestBuildReportSummary_Empty(t *testing.T) {
	s := buildReportSummary(nil)
	if s.Total != 0 {
		t.Errorf("expected total=0, got %d", s.Total)
	}
}
