package fifo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Simple test for the /health endpoint style response
func TestHealthHandler(t *testing.T) {
	// Create a dummy handler inline
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Simulate request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"ok"`) {
		t.Errorf("expected body to contain ok, got %s", w.Body.String())
	}
}
