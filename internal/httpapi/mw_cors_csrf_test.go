package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_OptionsAllowedOrigin_ReturnsNoContent(t *testing.T) {
	mw := CORS(CORSConfig{AllowedOrigins: []string{"https://app.example.com"}})

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/technologies", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "https://app.example.com" {
		t.Fatalf("expected access-control-allow-origin header, got %q", rr.Header().Get("Access-Control-Allow-Origin"))
	}
	if nextCalled {
		t.Fatal("next must not be called for OPTIONS preflight")
	}
}

func TestCORS_DisallowedOrigin_DoesNotSetHeaders(t *testing.T) {
	mw := CORS(CORSConfig{AllowedOrigins: []string{"https://app.example.com"}})

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/technologies", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Fatalf("did not expect CORS allow origin header, got %q", rr.Header().Get("Access-Control-Allow-Origin"))
	}
	if !nextCalled {
		t.Fatal("expected next to be called")
	}
}

func TestCSRF_StateChanging_DisallowedOrigin_Forbidden(t *testing.T) {
	mw := CSRF(CSRFConfig{TrustedOrigins: []string{"https://app.example.com"}})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/preferences", nil)
	req.Header.Set("Origin", "https://evil.example.com")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusForbidden, rr.Code, rr.Body.String())
	}
}

func TestCSRF_StateChanging_AllowedOrigin_Passes(t *testing.T) {
	mw := CSRF(CSRFConfig{TrustedOrigins: []string{"https://app.example.com"}})

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/preferences", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if !nextCalled {
		t.Fatal("expected next to be called")
	}
}

func TestCSRF_StateChanging_DisallowedReferer_Forbidden(t *testing.T) {
	mw := CSRF(CSRFConfig{TrustedOrigins: []string{"https://app.example.com"}})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPut, "/api/preferences", nil)
	req.Header.Set("Referer", "https://evil.example.com/path")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusForbidden, rr.Code, rr.Body.String())
	}
}
