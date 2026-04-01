package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/flux-iac/tofu-controller/internal/web/kubeclient"
)

const (
	serverTimeout = time.Minute
)

// RunServer starts the web server and blocks until the provided context is canceled.
// It serves the embedded SPA frontend and the /api/v1/* REST endpoints.
func RunServer(ctx context.Context, c cluster.Cluster, port int) error {
	l := ctrl.Log.WithName("web-server").WithValues("port", port)

	// Build Kubernetes client wrapper.
	kubeClient := kubeclient.New(c.GetAPIReader(), c.GetClient(), c.GetScheme())

	// Build the HTTP handler.
	handler := NewHandler(kubeClient, l)

	// Start HTTP server.
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  serverTimeout,
		WriteTimeout: serverTimeout,
		IdleTimeout:  serverTimeout,
		Handler:      handler,
	}

	serverStopped := make(chan error, 1)
	go func() {
		l.Info("starting web server")
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Error(err, "web server stopped unexpectedly")
			os.Exit(1)
		}
		close(serverStopped)
	}()

	// Wait for context cancellation, then shut down gracefully.
	<-ctx.Done()
	l.Info("shutting down web server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("web server shutdown error: %w", err)
	}

	select {
	case <-shutdownCtx.Done():
		return errors.New("web server graceful shutdown timed out")
	case <-serverStopped:
		return nil
	}
}
