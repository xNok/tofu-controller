package v1alpha2

import (
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	validation "k8s.io/apimachinery/pkg/util/validation"
)

func TestSafeLabelValue(t *testing.T) {
	g := NewWithT(t)

	g.Expect(SafeLabelValue("short")).To(Equal("short"))

	exactly63 := strings.Repeat("a", 63)
	g.Expect(SafeLabelValue(exactly63)).To(Equal(exactly63), "63 characters is within the limit")

	tooLong := strings.Repeat("a", 200)
	truncated := SafeLabelValue(tooLong)
	g.Expect(truncated).To(HaveLen(63))
	g.Expect(validation.IsValidLabelValue(truncated)).To(BeEmpty())

	// Distinct inputs that share a prefix must stay distinct.
	other := strings.Repeat("a", 199) + "b"
	g.Expect(SafeLabelValue(other)).NotTo(Equal(truncated))
}

func TestOwnedLabels(t *testing.T) {
	g := NewWithT(t)

	tf := &Terraform{ObjectMeta: metav1.ObjectMeta{Name: "my-stack", Namespace: "flux-system"}}
	labels := OwnedLabels(tf)

	g.Expect(labels).To(HaveKeyWithValue(ManagedByLabel, ManagedByValue))
	g.Expect(labels).To(HaveKeyWithValue(TerraformNameLabel, "my-stack"))
	g.Expect(labels).To(HaveKeyWithValue(TerraformNamespaceLabel, "flux-system"))

	// No annotation when nothing was truncated.
	g.Expect(OwnedAnnotations(tf)).To(BeEmpty())
}

func TestOwnedLabels_LongNameStaysValidAndRecoverable(t *testing.T) {
	g := NewWithT(t)

	long := strings.Repeat("n", 120)
	tf := &Terraform{ObjectMeta: metav1.ObjectMeta{Name: long, Namespace: "flux-system"}}

	labels := OwnedLabels(tf)
	g.Expect(validation.IsValidLabelValue(labels[TerraformNameLabel])).To(BeEmpty(),
		"a name past the limit must still produce a usable label value")
	g.Expect(labels[TerraformNameLabel]).NotTo(Equal(long))

	// The full name is still recoverable from the annotation.
	g.Expect(OwnedAnnotations(tf)).To(HaveKeyWithValue(TerraformFullNameAnnotation, long))
}

func TestOwnedLabels_NilTerraformDegradesToManaged(t *testing.T) {
	g := NewWithT(t)

	g.Expect(OwnedLabels(nil)).To(Equal(ManagedLabels()))
	g.Expect(OwnedAnnotations(nil)).To(BeEmpty())
}

// Attribution must not be breakable from spec, so our keys win over a
// user-supplied label map.
func TestWithOwnedLabels_OurKeysWin(t *testing.T) {
	g := NewWithT(t)

	userLabels := map[string]string{
		"team":                  "platform",
		ManagedByLabel:          "something-else",
		TerraformNameLabel:      "lies",
		TerraformNamespaceLabel: "lies",
	}

	merged := WithOwnedLabels(userLabels, "my-stack", "flux-system")

	g.Expect(merged).To(HaveKeyWithValue("team", "platform"), "unrelated user labels survive")
	g.Expect(merged).To(HaveKeyWithValue(ManagedByLabel, ManagedByValue))
	g.Expect(merged).To(HaveKeyWithValue(TerraformNameLabel, "my-stack"))
	g.Expect(merged).To(HaveKeyWithValue(TerraformNamespaceLabel, "flux-system"))

	// The caller's map must not be mutated.
	g.Expect(userLabels).To(HaveKeyWithValue(ManagedByLabel, "something-else"))
}

func TestWithOwnedLabels_NilExisting(t *testing.T) {
	g := NewWithT(t)

	merged := WithOwnedLabels(nil, "my-stack", "flux-system")
	g.Expect(merged).To(HaveKeyWithValue(TerraformNameLabel, "my-stack"))
}

func TestWithOwnedAnnotations_UnchangedWhenNothingToAdd(t *testing.T) {
	g := NewWithT(t)

	existing := map[string]string{"note": "keep me"}

	g.Expect(WithOwnedAnnotations(existing, "short-name")).To(Equal(existing))
	g.Expect(WithOwnedAnnotations(nil, "short-name")).To(BeEmpty())

	long := strings.Repeat("n", 120)
	merged := WithOwnedAnnotations(existing, long)
	g.Expect(merged).To(HaveKeyWithValue("note", "keep me"))
	g.Expect(merged).To(HaveKeyWithValue(TerraformFullNameAnnotation, long))
	g.Expect(existing).NotTo(HaveKey(TerraformFullNameAnnotation), "caller's map must not be mutated")
}

func TestWithManagedLabels(t *testing.T) {
	g := NewWithT(t)

	merged := WithManagedLabels(map[string]string{"existing": "value"})

	g.Expect(merged).To(HaveKeyWithValue("existing", "value"))
	g.Expect(merged).To(HaveKeyWithValue(ManagedByLabel, ManagedByValue))
	g.Expect(merged).NotTo(HaveKey(TerraformNameLabel),
		"managed-only objects must not claim an owning Terraform")
}
