package web

import (
	"net/http"

	"github.com/go-logr/logr"

	"github.com/flux-iac/tofu-controller/internal/web/kubeclient"
	webfs "github.com/flux-iac/tofu-controller/web"
)

// Handler provides HTTP handlers for the API endpoints and SPA static files.
type Handler struct {
	kubeClient *kubeclient.Client
	logger     logr.Logger
}

// NewHandler creates a new HTTP handler for the web server, wiring all routes
// and wrapping the mux with the standard middleware chain.
func NewHandler(kubeClient *kubeclient.Client, l logr.Logger) http.Handler {
	h := &Handler{
		kubeClient: kubeClient,
		logger:     l,
	}

	// Build SPA handler from embedded filesystem.
	spaHandler := http.FileServer(http.FS(NewFileSystem(webfs.GetFS())))

	// Create HTTP request multiplexer.
	mux := http.NewServeMux()

	// SPA — serves index.html for all non-API paths.
	mux.Handle("/", spaHandler)

	// API routes.
	mux.HandleFunc("GET /api/v1/report", h.ReportHandler)
	mux.HandleFunc("GET /api/v1/resources", h.ResourcesHandler)
	mux.HandleFunc("GET /api/v1/resource", h.ResourceHandler)
	mux.HandleFunc("GET /api/v1/plan", h.PlanHandler)
	mux.HandleFunc("POST /api/v1/resource/action", h.ActionHandler)

	// Wrap the mux with middlewares (innermost → outermost):
	// Logging → SecurityHeaders → Gzip → CacheControl → MaxBodySize → mux
	return LoggingMiddleware(l, SecurityHeadersMiddleware(
		GzipMiddleware(CacheControlMiddleware(MaxBodySizeMiddleware(1<<20)(mux)))))
}

// notImplemented is a temporary placeholder returning 501 for unbuilt handlers.
func notImplemented(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

// ResourceHandler returns a single Terraform resource by namespace and name.
// GET /api/v1/resource?namespace=...&name=...
func (h *Handler) ResourceHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

// PlanHandler returns the decoded plan content for a Terraform resource.
// GET /api/v1/plan?namespace=...&name=...
func (h *Handler) PlanHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}

// ActionHandler performs an action on a Terraform resource.
// POST /api/v1/resource/action
// Supported actions: approve, suspend, resume, force-reconcile
func (h *Handler) ActionHandler(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, r)
}
