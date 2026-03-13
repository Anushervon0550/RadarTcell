package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRateLimit_AllowsWithinLimit(t *testing.T) {
	mw := RateLimit(RateLimitConfig{Limit: 2, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected status %d, got %d", http.StatusNoContent, rr.Code)
		}
	}
}

func TestRateLimit_BlocksAfterLimit(t *testing.T) {
	mw := RateLimit(RateLimitConfig{Limit: 1, Window: time.Minute, Message: "limited"})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	first := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	first.RemoteAddr = "10.0.0.2:5555"
	firstRR := httptest.NewRecorder()
	h.ServeHTTP(firstRR, first)
	if firstRR.Code != http.StatusOK {
		t.Fatalf("expected first status %d, got %d", http.StatusOK, firstRR.Code)
	}

	second := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	second.RemoteAddr = "10.0.0.2:6666"
	secondRR := httptest.NewRecorder()
	h.ServeHTTP(secondRR, second)
	if secondRR.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusTooManyRequests, secondRR.Code, secondRR.Body.String())
	}
	if secondRR.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}
	if !strings.Contains(strings.ToLower(secondRR.Body.String()), "limited") {
		t.Fatalf("expected body to contain custom message, got: %s", secondRR.Body.String())
	}
}

func TestRateLimit_UsesDifferentKeysIndependently(t *testing.T) {
	mw := RateLimit(RateLimitConfig{Limit: 1, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	reqA := httptest.NewRequest(http.MethodPost, "/api/preferences", nil)
	reqA.RemoteAddr = "10.0.0.3:1000"
	rrA := httptest.NewRecorder()
	h.ServeHTTP(rrA, reqA)
	if rrA.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rrA.Code)
	}

	reqB := httptest.NewRequest(http.MethodPost, "/api/preferences", nil)
	reqB.RemoteAddr = "10.0.0.4:1000"
	rrB := httptest.NewRecorder()
	h.ServeHTTP(rrB, reqB)
	if rrB.Code != http.StatusOK {
		t.Fatalf("expected status %d for different key, got %d", http.StatusOK, rrB.Code)
	}
}

func TestRateLimit_IgnoresForwardedHeadersForKey(t *testing.T) {
	mw := RateLimit(RateLimitConfig{Limit: 1, Window: time.Minute})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	first := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	first.RemoteAddr = "8.8.8.8:1234"
	first.Header.Set("X-Forwarded-For", "10.0.0.10")
	first.Header.Set("X-Real-IP", "10.0.0.11")
	firstRR := httptest.NewRecorder()
	h.ServeHTTP(firstRR, first)
	if firstRR.Code != http.StatusOK {
		t.Fatalf("expected first status %d, got %d", http.StatusOK, firstRR.Code)
	}

	second := httptest.NewRequest(http.MethodPost, "/api/admin/login", nil)
	second.RemoteAddr = "8.8.8.8:5678"
	second.Header.Set("X-Forwarded-For", "127.0.0.1")
	second.Header.Set("X-Real-IP", "127.0.0.1")
	secondRR := httptest.NewRecorder()
	h.ServeHTTP(secondRR, second)
	if secondRR.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, secondRR.Code)
	}
}
