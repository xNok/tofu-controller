package controllers

import (
	"crypto/rand"
	"testing"

	"github.com/flux-iac/tofu-controller/api/plan"
	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This file is the registry of every object the controller and its runners
// create. A new creation site that does not carry the common labels should
// fail here rather than be noticed in a cluster.
//
// Covered below (pure constructors):
//   - runner Pod                    controllers/tf_controller_runner.go
//   - plan Secret and its chunks    api/plan/plan.go
//   - plan ConfigMap                api/plan/plan.go
//
// Covered elsewhere, because the object is only produced by a running
// component:
//   - outputs Secret                controllers/tc000041_all_outputs_test.go
//   - runner TLS Secret             controllers/tc009990_mtls_generate_creds_test.go
//
// The tfstate Secret is written by Terraform's Kubernetes backend, not by us,
// and is deliberately absent.

const labelTestRevision = "v2.6@sha256:c7fd0cc69b924aa5f9a6928477311737e439ca1b9e444855b0377e8a8ec65bb5"

func labelTestTerraform(name, namespace string) *infrav1.Terraform {
	return &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
	}
}

func labelTestPlan(t *testing.T, name, namespace string) *plan.Plan {
	t.Helper()

	p, err := plan.NewFromBytes(name, namespace, "default", "uuid-1", "plan-abc", []byte("plan body"))
	if err != nil {
		t.Fatalf("building plan: %s", err)
	}

	return p
}

func TestCreatedObjectsCarryOwnershipLabels(t *testing.T) {
	const (
		name      = "my-stack"
		namespace = "flux-system"
	)

	cases := []struct {
		object string
		labels func(t *testing.T) map[string]string
	}{
		{
			object: "runner Pod",
			labels: func(t *testing.T) map[string]string {
				pod, err := runnerPodTemplate(labelTestTerraform(name, namespace), "runner.tls-123", labelTestRevision)
				if err != nil {
					t.Fatalf("building runner pod: %s", err)
				}
				return pod.Labels
			},
		},
		{
			object: "plan Secret",
			labels: func(t *testing.T) map[string]string {
				secrets, err := labelTestPlan(t, name, namespace).ToSecret("")
				if err != nil {
					t.Fatalf("building plan secret: %s", err)
				}
				return secrets[0].Labels
			},
		},
		{
			object: "plan ConfigMap",
			labels: func(t *testing.T) map[string]string {
				configMaps, err := labelTestPlan(t, name, namespace).ToConfigMap("")
				if err != nil {
					t.Fatalf("building plan configmap: %s", err)
				}
				return configMaps[0].Labels
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.object, func(t *testing.T) {
			g := NewWithT(t)
			labels := tc.labels(t)

			g.Expect(labels).To(HaveKeyWithValue(infrav1.ManagedByLabel, infrav1.ManagedByValue))
			g.Expect(labels).To(HaveKeyWithValue(infrav1.TerraformNameLabel, name))
			g.Expect(labels).To(HaveKeyWithValue(infrav1.TerraformNamespaceLabel, namespace))
		})
	}
}

// The plan identity labels are what every consumer lists chunks by, so the
// ownership labels must be additive rather than replacing them.
func TestPlanObjectsKeepPlanIdentityLabels(t *testing.T) {
	g := NewWithT(t)

	secrets, err := labelTestPlan(t, "my-stack", "flux-system").ToSecret("")
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(secrets[0].Labels).To(HaveKeyWithValue(plan.TFPlanNameLabel, "my-stack"))
	g.Expect(secrets[0].Labels).To(HaveKeyWithValue(plan.TFPlanWorkspaceLabel, "default"))
}

// Every chunk of a split plan must be labelled, not just the first.
func TestPlanChunksAreAllLabelled(t *testing.T) {
	g := NewWithT(t)

	// Genuinely random, so it does not compress below the chunk threshold.
	// A patterned filler would gzip to a single chunk and quietly skip the
	// case this test exists for.
	body := make([]byte, 4*1024*1024)
	_, err := rand.Read(body)
	g.Expect(err).NotTo(HaveOccurred())

	p, err := plan.NewFromBytes("my-stack", "flux-system", "default", "uuid-1", "plan-abc", body)
	g.Expect(err).NotTo(HaveOccurred())

	secrets, err := p.ToSecret("")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(secrets)).To(BeNumerically(">", 1), "test needs a chunked plan")

	for _, secret := range secrets {
		g.Expect(secret.Labels).To(HaveKeyWithValue(infrav1.ManagedByLabel, infrav1.ManagedByValue), secret.Name)
		g.Expect(secret.Labels).To(HaveKeyWithValue(infrav1.TerraformNameLabel, "my-stack"), secret.Name)
	}
}

// A custom runner pod label must not be able to overwrite the keys the
// controller selects and attributes on.
func TestRunnerPodCustomLabelsCannotOverrideOwnership(t *testing.T) {
	g := NewWithT(t)

	terraform := labelTestTerraform("my-stack", "flux-system")
	terraform.Spec.RunnerPodTemplate.Metadata.Labels = map[string]string{
		"team":                     "platform",
		infrav1.ManagedByLabel:     "someone-else",
		infrav1.TerraformNameLabel: "lies",
		"company.com/custom":       "kept",
	}

	pod, err := runnerPodTemplate(terraform, "runner.tls-123", labelTestRevision)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(pod.Labels).To(HaveKeyWithValue("team", "platform"))
	g.Expect(pod.Labels).To(HaveKeyWithValue("company.com/custom", "kept"))
	g.Expect(pod.Labels).To(HaveKeyWithValue(infrav1.ManagedByLabel, infrav1.ManagedByValue))
	g.Expect(pod.Labels).To(HaveKeyWithValue(infrav1.TerraformNameLabel, "my-stack"))
}

// A name past the label length limit must still produce a usable label, with
// the full value recoverable from the annotation.
func TestLongTerraformNameStaysLabelled(t *testing.T) {
	g := NewWithT(t)

	long := "a-very-long-terraform-object-name-that-goes-well-past-the-sixty-three-character-label-limit"
	terraform := labelTestTerraform(long, "flux-system")

	pod, err := runnerPodTemplate(terraform, "runner.tls-123", labelTestRevision)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(pod.Labels[infrav1.TerraformNameLabel]).To(HaveLen(63))
	g.Expect(pod.Annotations).To(HaveKeyWithValue(infrav1.TerraformFullNameAnnotation, long))
}
