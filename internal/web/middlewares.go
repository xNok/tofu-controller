package web

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// gzipResponseWriter wraps http.ResponseWriter to compress responses with gzip.
type gzipResponseWriter struct {
	http.ResponseWriter
	writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	if gzipWriter, ok := w.writer.(*gzip.Writer); ok {
		_ = gzipWriter.Flush()
	}
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// GzipMiddleware adds gzip compression to responses when the client accepts it.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length")

		next.ServeHTTP(&gzipResponseWriter{
			ResponseWriter: w,
			writer:         gzipWriter,
		}, r)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// CacheControlMiddleware sets appropriate cache headers for static assets.
// Hashed assets in /assets/ are cached indefinitely; all other files are not cached.
func CacheControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeadersMiddleware adds security headers to all responses.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("X-Robots-Tag", "noindex, nofollow")
		next.ServeHTTP(w, r)
	})
}

// MaxBodySizeMiddleware limits the size of request bodies to prevent abuse.
// Only applies to POST, PUT, and PATCH requests.
func MaxBodySizeMiddleware(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs HTTP requests and responses with their duration.
func LoggingMiddleware(logger logr.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		reqFields := map[string]any{
			"method":     r.Method,
			"path":       r.URL.RequestURI(),
			"remote":     r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}
		ctx := log.IntoContext(r.Context(), logger.WithValues("httpRequest", reqFields))
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		duration := time.Since(start)
		logger.V(1).Info("HTTP request completed", "httpRequest", reqFields, "httpResponse", map[string]any{
			"status":     wrapped.statusCode,
			"latency_ms": duration.Round(time.Millisecond).Milliseconds(),
		})
	})
}
