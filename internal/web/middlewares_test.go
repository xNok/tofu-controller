package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-logr/logr"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func TestSecurityHeadersMiddleware(t *testing.T) {
	handler := SecurityHeadersMiddleware(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	headers := map[string]string{
		"X-Frame-Options":        "DENY",
		"X-Content-Type-Options": "nosniff",
		"X-XSS-Protection":       "1; mode=block",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"X-Robots-Tag":           "noindex, nofollow",
	}
	for header, expected := range headers {
		got := rr.Header().Get(header)
		if got != expected {
			t.Errorf("expected %s: %q, got: %q", header, expected, got)
		}
	}
}

func TestCacheControlMiddleware_Assets(t *testing.T) {
	handler := CacheControlMiddleware(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/assets/main-abc123.js", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if !strings.Contains(cc, "immutable") {
		t.Errorf("expected immutable cache for /assets/ path, got: %q", cc)
	}
}

func TestCacheControlMiddleware_Index(t *testing.T) {
	handler := CacheControlMiddleware(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	cc := rr.Header().Get("Cache-Control")
	if !strings.Contains(cc, "no-cache") {
		t.Errorf("expected no-cache for non-asset path, got: %q", cc)
	}
}

func TestGzipMiddleware_WithAcceptEncoding(t *testing.T) {
	handler := GzipMiddleware(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Encoding") != "gzip" {
		t.Error("expected Content-Encoding: gzip when client accepts gzip")
	}
}

func TestGzipMiddleware_WithoutAcceptEncoding(t *testing.T) {
	handler := GzipMiddleware(okHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Content-Encoding") == "gzip" {
		t.Error("expected no gzip when client does not accept gzip")
	}
}

func TestMaxBodySizeMiddleware_POST(t *testing.T) {
	// 10 bytes limit
	handler := MaxBodySizeMiddleware(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read entire body — should be limited
		buf := make([]byte, 100)
		n, _ := r.Body.Read(buf)
		_, _ = w.Write(buf[:n])
	}))

	body := strings.NewReader("this body is longer than 10 bytes")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/resource/action", body)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// The handler should not have been able to read more than 10 bytes
	if len(rr.Body.String()) > 10 {
		t.Errorf("expected body to be limited to 10 bytes, got %d bytes", len(rr.Body.String()))
	}
}

func TestMaxBodySizeMiddleware_GET(t *testing.T) {
	// GET requests should not have their body limited
	called := false
	handler := MaxBodySizeMiddleware(1)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/resources", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected handler to be called for GET request")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := LoggingMiddleware(logr.Discard(), okHandler())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/report", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}
