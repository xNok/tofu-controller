package v1alpha2

import (
	"crypto/sha256"
	"fmt"
	"maps"
)

const (
	// ManagedByLabel marks an object as created by this controller. It is the
	// Kubernetes-recommended key for the tool managing an object, and is set on
	// everything the controller and its runners create.
	ManagedByLabel = "app.kubernetes.io/managed-by"

	// ManagedByValue is the value of ManagedByLabel on objects we create.
	ManagedByValue = "tofu-controller"

	// TerraformNameLabel and TerraformNamespaceLabel attribute an object to the
	// Terraform that produced it, so `kubectl get -l` can answer "what did this
	// Terraform create?" without reading status.
	//
	// The name is passed through SafeLabelValue and may therefore be truncated;
	// TerraformFullNameAnnotation always carries it intact.
	TerraformNameLabel      = "infra.contrib.fluxcd.io/terraform-name"
	TerraformNamespaceLabel = "infra.contrib.fluxcd.io/terraform-namespace"

	// TerraformFullNameAnnotation carries the untruncated Terraform name,
	// mirroring the plan-full-name annotation that api/plan already sets for
	// the same reason.
	TerraformFullNameAnnotation = "infra.contrib.fluxcd.io/terraform-full-name"
)

// maxLabelValueLength is the Kubernetes limit on a label value.
const maxLabelValueLength = 63

// SafeLabelValue returns a string that is safe to use as a Kubernetes label
// value, truncating and appending a checksum when it exceeds the 63-character
// limit so that distinct inputs stay distinct.
func SafeLabelValue(value string) string {
	// Values that are equal to or less than 63 characters are already good
	if len(value) <= maxLabelValueLength {
		return value
	}

	// Create hash
	checksum := sha256.Sum256([]byte(value))

	// Build a prefix to append to end of truncated value
	checksumPrefix := fmt.Sprintf("-%x", checksum[:8])

	prefix := value[:maxLabelValueLength-len(checksumPrefix)]

	return prefix + checksumPrefix
}

// ManagedLabels marks an object as created by this controller.
//
// Use it for objects that belong to the installation rather than to a single
// Terraform — the per-namespace runner TLS Secret, for example, is shared by
// every runner in its namespace and has no one owner.
func ManagedLabels() map[string]string {
	return map[string]string{
		ManagedByLabel: ManagedByValue,
	}
}

// OwnedLabelsFor marks an object as produced by the named Terraform. Takes
// strings rather than the object because api/plan builds its Secrets from a
// name and namespace it already carries. Includes ManagedLabels, so callers
// set one or the other, never both.
func OwnedLabelsFor(name, namespace string) map[string]string {
	return map[string]string{
		ManagedByLabel:          ManagedByValue,
		TerraformNameLabel:      SafeLabelValue(name),
		TerraformNamespaceLabel: SafeLabelValue(namespace),
	}
}

// OwnedLabels marks an object as produced by a specific Terraform.
func OwnedLabels(terraform *Terraform) map[string]string {
	if terraform == nil {
		return ManagedLabels()
	}

	return OwnedLabelsFor(terraform.Name, terraform.Namespace)
}

// OwnedAnnotationsFor carries the values that do not survive the label length
// limit. Returns nil when nothing was truncated, so no annotation is written
// for the common case.
func OwnedAnnotationsFor(name string) map[string]string {
	if name == SafeLabelValue(name) {
		return nil
	}

	return map[string]string{
		TerraformFullNameAnnotation: name,
	}
}

// OwnedAnnotations carries the untruncated name for a specific Terraform.
func OwnedAnnotations(terraform *Terraform) map[string]string {
	if terraform == nil {
		return nil
	}

	return OwnedAnnotationsFor(terraform.Name)
}

// WithOwnedLabels merges the ownership labels into an existing map, which may
// be nil, and returns the result.
//
// Our keys are written last and deliberately win: these labels are what
// tooling and cleanup select on, so a user-supplied label map — from
// spec.writeOutputsToSecret.labels or spec.runnerPodTemplate.metadata.labels —
// must not be able to silently break attribution.
func WithOwnedLabels(existing map[string]string, name, namespace string) map[string]string {
	merged := make(map[string]string, len(existing)+3)
	maps.Copy(merged, existing)
	maps.Copy(merged, OwnedLabelsFor(name, namespace))

	return merged
}

// WithOwnedAnnotations merges the ownership annotations into an existing map,
// which may be nil. Returns existing unchanged when there is nothing to add.
func WithOwnedAnnotations(existing map[string]string, name string) map[string]string {
	owned := OwnedAnnotationsFor(name)
	if len(owned) == 0 {
		return existing
	}

	merged := make(map[string]string, len(existing)+len(owned))
	maps.Copy(merged, existing)
	maps.Copy(merged, owned)

	return merged
}

// WithManagedLabels merges the managed-by label into an existing map, which
// may be nil, and returns the result.
func WithManagedLabels(existing map[string]string) map[string]string {
	merged := make(map[string]string, len(existing)+1)
	maps.Copy(merged, existing)
	maps.Copy(merged, ManagedLabels())

	return merged
}
