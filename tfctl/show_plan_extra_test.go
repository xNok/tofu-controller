package tfctl

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/fluxcd/pkg/apis/meta"
	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
)

func TestPlanExists_TypedRefPreferred(t *testing.T) {
	g := NewWithT(t)
	tf := &infrav1.Terraform{
		Status: infrav1.TerraformStatus{
			ArtifactRefs: &infrav1.ArtifactReferences{
				Plan: &infrav1.PlanArtifactRef{
					Ref: meta.NamespacedObjectKindReference{
						APIVersion: "v1",
						Kind:       "ConfigMap",
						Namespace:  "shared",
						Name:       "tfplan-default-hello",
					},
				},
			},
		},
	}

	g.Expect(planExists(tf)).To(BeTrue())
	g.Expect(planRefNamespace(tf, "fallback-ns")).To(Equal("shared"))
}

func TestPlanExists_FallbackToLegacyField(t *testing.T) {
	g := NewWithT(t)
	tf := &infrav1.Terraform{}
	g.Expect(planExists(tf)).To(BeFalse())

	tf.Status.Plan.Pending = "plan-id-123"
	g.Expect(planExists(tf)).To(BeTrue())
	g.Expect(planRefNamespace(tf, "fallback-ns")).To(Equal("fallback-ns"))
}
