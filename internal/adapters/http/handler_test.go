package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/popcoondev/gcp-service-learning/internal/application"
)

func TestHealthEndpoint(t *testing.T) {
	handler := NewHandler(application.OrderService{}).Routes()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestHealthEndpointRejectsNonGet(t *testing.T) {
	handler := NewHandler(application.OrderService{}).Routes()

	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}
