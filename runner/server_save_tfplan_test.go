package runner

import (
	"crypto/rand"
	"testing"

	"github.com/flux-iac/tofu-controller/api/plan"
	"github.com/flux-iac/tofu-controller/api/planid"
	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	"github.com/go-logr/logr"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestSaveTFPlanWithMultipleChunks(t *testing.T) {
	g := NewGomegaWithT(t)

	terraform := &infrav1.Terraform{}

	server := &TerraformRunnerServer{
		InstanceID: "51b32416-d76d-4720-b2ef-1c13996d3c4a",
		terraform:  terraform,
	}

	req := &SaveTFPlanRequest{
		Name:       "a-really-large-plan",
		Namespace:  "terraform",
		TfInstance: server.InstanceID,
		Revision:   "master@sha1:453f0eeb",
	}

	// create 4MB of plan data.
	// intentionally not gzipping as this will make the data a lot smaller
	planData := make([]byte, 4*1024*1024) // 4MB of plans!

	planId := planid.GetPlanID(req.Revision)

	plan, err := plan.NewFromBytes(req.Name, req.Namespace, server.terraform.WorkspaceName(), "plan-uuid-1", planId, planData)
	g.Expect(err).NotTo(HaveOccurred())

	secrets, err := plan.ToSecret("")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(len(secrets)).To(Equal(4), "plan data should have been chunked into four secrets")
}

func TestWritePlanAsSecretReportsSavedObjects(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := t.Context()
	log := logr.Discard()

	terraform := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Namespace: "terraform", Name: "a-plan"},
	}

	server := &TerraformRunnerServer{
		InstanceID: "51b32416-d76d-4720-b2ef-1c13996d3c4a",
		terraform:  terraform,
		Client:     fake.NewClientBuilder().Build(),
	}

	planId := planid.GetPlanID("master@sha1:453f0eeb")

	tfPlan, err := plan.NewFromBytes("a-plan", "terraform", terraform.WorkspaceName(), "plan-uuid-1", planId, []byte("plan data"))
	g.Expect(err).NotTo(HaveOccurred())

	saved, err := server.writePlanAsSecret(ctx, "a-plan", "terraform", log, planId, tfPlan, "", "plan-uuid-1")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(saved).To(HaveLen(1))
	g.Expect(saved[0].ApiVersion).To(Equal("v1"))
	g.Expect(saved[0].Kind).To(Equal("Secret"))
	g.Expect(saved[0].Namespace).To(Equal("terraform"))
	g.Expect(saved[0].Name).To(Equal("tfplan-" + terraform.WorkspaceName() + "-a-plan"))

	// api/plan is the single place that owns the ownerReference; the runner
	// must not add a second one pointing at the same Terraform.
	written := &corev1.Secret{}
	g.Expect(server.Get(ctx, types.NamespacedName{Namespace: "terraform", Name: saved[0].Name}, written)).To(Succeed())
	g.Expect(written.OwnerReferences).To(HaveLen(1))
	g.Expect(written.OwnerReferences[0].Kind).To(Equal(infrav1.TerraformKind))
	g.Expect(string(written.OwnerReferences[0].UID)).To(Equal("plan-uuid-1"))
}

func TestWritePlanAsSecretReportsEveryChunk(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := t.Context()
	log := logr.Discard()

	terraform := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Namespace: "terraform", Name: "a-really-large-plan"},
	}

	server := &TerraformRunnerServer{
		InstanceID: "51b32416-d76d-4720-b2ef-1c13996d3c4a",
		terraform:  terraform,
		Client:     fake.NewClientBuilder().Build(),
	}

	planId := planid.GetPlanID("master@sha1:453f0eeb")

	// 4MB of incompressible data so the plan is chunked.
	planData := make([]byte, 4*1024*1024)
	_, err := rand.Read(planData)
	g.Expect(err).NotTo(HaveOccurred())

	tfPlan, err := plan.NewFromBytes("a-really-large-plan", "terraform", terraform.WorkspaceName(), "plan-uuid-1", planId, planData)
	g.Expect(err).NotTo(HaveOccurred())

	expected, err := tfPlan.ToSecret("")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(expected)).To(BeNumerically(">", 1), "test needs a chunked plan")

	saved, err := server.writePlanAsSecret(ctx, "a-really-large-plan", "terraform", log, planId, tfPlan, "", "plan-uuid-1")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(saved).To(HaveLen(len(expected)), "every chunk must be reported so status.chunks is accurate")
	g.Expect(saved[0].Name).To(Equal(expected[0].Name), "the first entry must be chunk 0")
}

func TestWritePlanAsConfigMapReportsSavedObjects(t *testing.T) {
	g := NewGomegaWithT(t)
	ctx := t.Context()
	log := logr.Discard()

	terraform := &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Namespace: "terraform", Name: "a-plan"},
	}

	server := &TerraformRunnerServer{
		InstanceID: "51b32416-d76d-4720-b2ef-1c13996d3c4a",
		terraform:  terraform,
		Client:     fake.NewClientBuilder().Build(),
	}

	planId := planid.GetPlanID("master@sha1:453f0eeb")

	tfPlan, err := plan.NewFromBytes("a-plan", "terraform", terraform.WorkspaceName(), "plan-uuid-1", planId, []byte("human readable plan"))
	g.Expect(err).NotTo(HaveOccurred())

	saved, err := server.writePlanAsConfigMap(ctx, "a-plan", "terraform", log, planId, tfPlan, "", "plan-uuid-1")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(saved).To(HaveLen(1))
	g.Expect(saved[0].Kind).To(Equal("ConfigMap"))

	written := &corev1.ConfigMap{}
	g.Expect(server.Get(ctx, types.NamespacedName{Namespace: "terraform", Name: saved[0].Name}, written)).To(Succeed())
	g.Expect(written.OwnerReferences).To(HaveLen(1))
}
