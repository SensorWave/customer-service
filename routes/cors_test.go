package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	company "customer-service/internal/company"
	user "customer-service/internal/user"

	"github.com/gin-gonic/gin"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(CORSMiddleware())
	RegisterRoutes(r, company.NewHandler(nil), user.NewHandler(nil))

	return r
}

func TestCORSMiddleware_PreflightOptions(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/companies", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	rec := httptest.NewRecorder()
	newTestRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected wildcard allow-origin header, got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Fatal("expected allow-methods header to be set")
	}
	if rec.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Fatal("expected allow-headers header to be set")
	}
}

func TestCORSMiddleware_GetRequestIncludesHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	rec := httptest.NewRecorder()
	newTestRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("expected wildcard allow-origin header, got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}
