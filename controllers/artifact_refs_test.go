package controllers

import (
	"testing"

	"github.com/flux-iac/tofu-controller/runner"
	"github.com/fluxcd/pkg/apis/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"

	infrav1 "github.com/flux-iac/tofu-controller/api/v1alpha2"
)

func newTF(name string) *infrav1.Terraform {
	return &infrav1.Terraform{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "flux-system",
			Name:      name,
		},
	}
}

func TestSetPlanArtifactRef_Primary(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	setPlanArtifactRef(tf, []*runner.SavedObject{
		{ApiVersion: "v1", Kind: "Secret", Namespace: "flux-system", Name: "tfplan-default-hello"},
	}, "plan-abc")

	g.Expect(tf.Status.ArtifactRefs).NotTo(BeNil())
	g.Expect(tf.Status.ArtifactRefs.Plan).NotTo(BeNil())
	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.APIVersion).To(Equal("v1"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.Kind).To(Equal("Secret"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.Namespace).To(Equal("flux-system"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.Name).To(Equal("tfplan-default-hello"))
	g.Expect(tf.Status.ArtifactRefs.Plan.PlanID).To(Equal("plan-abc"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Chunks).To(Equal(int32(1)))
}

func TestSetPlanArtifactRef_Chunked(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	setPlanArtifactRef(tf, []*runner.SavedObject{
		{ApiVersion: "v1", Kind: "Secret", Namespace: "flux-system", Name: "tfplan-default-hello-0"},
		{ApiVersion: "v1", Kind: "Secret", Namespace: "flux-system", Name: "tfplan-default-hello-1"},
		{ApiVersion: "v1", Kind: "Secret", Namespace: "flux-system", Name: "tfplan-default-hello-2"},
	}, "plan-xyz")

	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.Name).To(Equal("tfplan-default-hello-0"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Chunks).To(Equal(int32(3)))
}

func TestSetPlanArtifactRef_ReplacesPreviousRef(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Status.ArtifactRefs = &infrav1.ArtifactReferences{
		Plan: &infrav1.PlanArtifactRef{
			Ref:    meta.NamespacedObjectKindReference{Kind: "Secret", Name: "tfplan-default-hello-0"},
			PlanID: "plan-old",
			Chunks: 3,
		},
	}

	setPlanArtifactRef(tf, []*runner.SavedObject{
		{ApiVersion: "v1", Kind: "Secret", Namespace: "flux-system", Name: "tfplan-default-hello"},
	}, "plan-new")

	g.Expect(tf.Status.ArtifactRefs.Plan.Ref.Name).To(Equal("tfplan-default-hello"))
	g.Expect(tf.Status.ArtifactRefs.Plan.PlanID).To(Equal("plan-new"))
	g.Expect(tf.Status.ArtifactRefs.Plan.Chunks).To(Equal(int32(1)), "stale chunk count must not survive")
}

func TestSetPlanArtifactRef_EmptyClears(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Status.ArtifactRefs = &infrav1.ArtifactReferences{
		Plan: &infrav1.PlanArtifactRef{PlanID: "stale"},
	}

	setPlanArtifactRef(tf, nil, "ignored")

	g.Expect(tf.Status.ArtifactRefs.Plan).To(BeNil())
}

func TestClearPlanArtifactRef_NoStatus(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	g.Expect(func() { clearPlanArtifactRef(tf) }).NotTo(Panic())
	g.Expect(tf.Status.ArtifactRefs).To(BeNil())
}

func TestSetOutputsArtifactRef(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	setOutputsArtifactRef(tf, "hello-outputs")

	g.Expect(tf.Status.ArtifactRefs.Outputs).NotTo(BeNil())
	g.Expect(tf.Status.ArtifactRefs.Outputs.APIVersion).To(Equal("v1"))
	g.Expect(tf.Status.ArtifactRefs.Outputs.Kind).To(Equal("Secret"))
	g.Expect(tf.Status.ArtifactRefs.Outputs.Namespace).To(Equal("flux-system"))
	g.Expect(tf.Status.ArtifactRefs.Outputs.Name).To(Equal("hello-outputs"))
}

func TestClearOutputsArtifactRef(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Status.ArtifactRefs = &infrav1.ArtifactReferences{
		Outputs: &meta.NamespacedObjectKindReference{Name: "stale"},
	}

	clearOutputsArtifactRef(tf)

	g.Expect(tf.Status.ArtifactRefs.Outputs).To(BeNil())
}

func TestSetStateArtifactRef_DefaultKubernetesBackend(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	setStateArtifactRef(tf, false)

	state := tf.Status.ArtifactRefs.State
	g.Expect(state).NotTo(BeNil())
	g.Expect(state.Backend).To(Equal("kubernetes"))
	g.Expect(state.Ref).NotTo(BeNil())
	g.Expect(state.Ref.APIVersion).To(Equal("v1"))
	g.Expect(state.Ref.Kind).To(Equal("Secret"))
	g.Expect(state.Ref.Namespace).To(Equal("flux-system"))
	g.Expect(state.Ref.Name).To(Equal("tfstate-default-hello"))
	g.Expect(state.ExternalLocation).To(BeEmpty())
}

func TestSetStateArtifactRef_ExplicitSecretSuffix(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		SecretSuffix: "customsuffix",
	}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("kubernetes"))
	g.Expect(tf.Status.ArtifactRefs.State.Ref.Name).To(Equal("tfstate-default-customsuffix"))
}

func TestSetStateArtifactRef_ExplicitBackendConfigEmptySuffix(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		InClusterConfig: true,
	}

	setStateArtifactRef(tf, false)

	// setupTerraform renders secret_suffix verbatim, so an empty suffix
	// really does produce a "tfstate-{workspace}-" Secret. The reference
	// mirrors that rather than inventing a default the backend never uses.
	g.Expect(tf.Status.ArtifactRefs.State.Ref.Name).To(Equal("tfstate-default-"))
}

func TestSetStateArtifactRef_BackendDisabled(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{Disable: true}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("disabled"))
	g.Expect(tf.Status.ArtifactRefs.State.Ref).To(BeNil())
	g.Expect(tf.Status.ArtifactRefs.State.ExternalLocation).To(BeEmpty())
}

func TestSetStateArtifactRef_S3CustomConfig(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		CustomConfiguration: `backend "s3" {
  bucket = "my-state-bucket"
  key    = "envs/prod/hello.tfstate"
  region = "us-east-1"
}`,
	}

	setStateArtifactRef(tf, false)

	state := tf.Status.ArtifactRefs.State
	g.Expect(state.Backend).To(Equal("s3"))
	g.Expect(state.Ref).To(BeNil())
	g.Expect(state.ExternalLocation).To(Equal("s3://my-state-bucket/envs/prod/hello.tfstate"))
}

func TestSetStateArtifactRef_GCSCustomConfig(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		CustomConfiguration: `backend "gcs" {
  bucket = "my-bucket"
  prefix = "tf/state"
}`,
	}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("gcs"))
	g.Expect(tf.Status.ArtifactRefs.State.ExternalLocation).To(Equal("gs://my-bucket/tf/state"))
}

func TestSetStateArtifactRef_UnknownCustomBackend(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		CustomConfiguration: `backend "consul" {
  path    = "tf/state"
  address = "consul.example.com:8500"
}`,
	}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("consul"))
	g.Expect(tf.Status.ArtifactRefs.State.Ref).To(BeNil())
	g.Expect(tf.Status.ArtifactRefs.State.ExternalLocation).To(BeEmpty())
}

func TestSetStateArtifactRef_CloudWorkspace(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.Cloud = &infrav1.CloudSpec{
		Organization: "my-org",
		Workspaces: &infrav1.CloudWorkspacesSpec{
			Name: "prod",
		},
	}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("cloud"))
	g.Expect(tf.Status.ArtifactRefs.State.Ref).To(BeNil())
	g.Expect(tf.Status.ArtifactRefs.State.ExternalLocation).To(Equal("cloud://my-org/prod"))
}

func TestSetStateArtifactRef_CustomWorkspace(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.Workspace = "staging"

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.Ref.Name).To(Equal("tfstate-staging-hello"))
}

func TestSetStateArtifactRef_DisableTFK8SBackend(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")

	setStateArtifactRef(tf, true)

	g.Expect(tf.Status.ArtifactRefs.State.Backend).To(Equal("local"))
	g.Expect(tf.Status.ArtifactRefs.State.Ref).To(BeNil())
}

func TestParseBackendType(t *testing.T) {
	g := NewWithT(t)
	g.Expect(parseBackendType(`backend "s3" { bucket = "x" }`)).To(Equal("s3"))
	g.Expect(parseBackendType(`backend "remote" {}`)).To(Equal("remote"))
	g.Expect(parseBackendType(`no backend block here`)).To(Equal("unknown"))
	g.Expect(parseBackendType(`backend "" {}`)).To(Equal("unknown"))
}

func TestExtractHCLString(t *testing.T) {
	g := NewWithT(t)
	hcl := `backend "s3" {
  bucket = "my-bucket"
  region = "us-east-1"
}`
	g.Expect(extractHCLString(hcl, "bucket")).To(Equal("my-bucket"))
	g.Expect(extractHCLString(hcl, "region")).To(Equal("us-east-1"))
	g.Expect(extractHCLString(hcl, "missing")).To(BeEmpty())
}

func TestExtractHCLString_WholeKeyOnly(t *testing.T) {
	g := NewWithT(t)
	hcl := `backend "s3" {
  bucket_prefix = "not-the-bucket"
  bucket        = "the-bucket"
}`
	// "bucket" must not be satisfied by the bucket_prefix line that precedes it.
	g.Expect(extractHCLString(hcl, "bucket")).To(Equal("the-bucket"))
	g.Expect(extractHCLString(hcl, "bucket_prefix")).To(Equal("not-the-bucket"))
	// A key that only exists as a substring of another key is absent.
	g.Expect(extractHCLString(hcl, "prefix")).To(BeEmpty())
}

func TestExtractHCLString_SingleLineBlock(t *testing.T) {
	g := NewWithT(t)
	g.Expect(extractHCLString(`backend "s3" { bucket = "inline" }`, "bucket")).To(Equal("inline"))
}

func TestSetStateArtifactRef_S3BucketPrefix(t *testing.T) {
	g := NewWithT(t)
	tf := newTF("hello")
	tf.Spec.BackendConfig = &infrav1.BackendConfigSpec{
		CustomConfiguration: `backend "s3" {
  bucket_prefix = "ignore-me"
  bucket        = "real-bucket"
  key           = "state.tfstate"
}`,
	}

	setStateArtifactRef(tf, false)

	g.Expect(tf.Status.ArtifactRefs.State.ExternalLocation).To(Equal("s3://real-bucket/state.tfstate"))
}
