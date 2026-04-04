package kubeclient

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Client wraps controller-runtime client for use in the web handlers.
type Client struct {
	reader client.Reader
	writer client.Client
	scheme *runtime.Scheme
}

// New creates a new Client.
func New(reader client.Reader, writer client.Client, scheme *runtime.Scheme) *Client {
	return &Client{
		reader: reader,
		writer: writer,
		scheme: scheme,
	}
}

// Get retrieves an object from the Kubernetes API.
func (c *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	return c.reader.Get(ctx, key, obj)
}

// List retrieves a list of objects from the Kubernetes API.
func (c *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return c.reader.List(ctx, list, opts...)
}

// Patch applies a patch to an object in the Kubernetes API.
func (c *Client) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return c.writer.Patch(ctx, obj, patch, opts...)
}

// Scheme returns the runtime scheme.
func (c *Client) Scheme() *runtime.Scheme {
	return c.scheme
}
