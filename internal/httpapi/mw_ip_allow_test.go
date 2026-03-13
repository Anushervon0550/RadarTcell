package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowPrivateNetworks_AllowsPrivateIP(t *testing.T) {
	mw := AllowPrivateNetworks()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.RemoteAddr = "10.10.10.10:1234"
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected %d, got %d", http.StatusNoContent, rr.Code)
	}
}

func TestAllowPrivateNetworks_BlocksPublicIP(t *testing.T) {
	mw := AllowPrivateNetworks()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.RemoteAddr = "8.8.8.8:1234"
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d", http.StatusForbidden, rr.Code)
	}
}

