package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

// ---- RequireRole ----

func TestRequireRole_AllowsMatchingRole(t *testing.T) {
	called := false
	h := RequireRole(domain.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxPrincipal, domain.Principal{Subject: "a", Role: domain.RoleAdmin}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent || !called {
		t.Fatalf("expected pass-through 204, got %d (called=%v)", rr.Code, called)
	}
}

func TestRequireRole_DeniesMismatchedRole(t *testing.T) {
	called := false
	h := RequireRole(domain.RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), ctxPrincipal, domain.Principal{Subject: "a", Role: "viewer"}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
	if called {
		t.Fatal("next must not be called for mismatched role")
	}
}

// ---- preferences ownership (IDOR) ----

func TestPreferencesHandler_Get_Forbidden_OtherUser(t *testing.T) {
	stub := &preferencesServiceStub{}
	h := NewPreferencesHandler(stub)

	req := httptest.NewRequest(http.MethodGet, "/api/preferences/victim", nil)
	req = withChiUserID(req, "victim")
	req = withPrincipal(req, "attacker")
	rr := httptest.NewRecorder()
	h.Get(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for cross-user access, got %d body=%s", rr.Code, rr.Body.String())
	}
	if stub.gotGet != "" {
		t.Fatalf("service must not be called on forbidden access, got %q", stub.gotGet)
	}
}

func TestPreferencesHandler_Save_Forbidden_OtherUser(t *testing.T) {
	stub := &preferencesServiceStub{}
	h := NewPreferencesHandler(stub)

	body := `{"user_id":"victim","settings":{"x":1}}`
	req := httptest.NewRequest(http.MethodPost, "/api/preferences", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withPrincipal(req, "attacker")
	rr := httptest.NewRecorder()
	h.Save(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when saving for another user, got %d body=%s", rr.Code, rr.Body.String())
	}
}

// ---- sanitizeQuery ----

func TestSanitizeQuery_RedactsSensitiveKeys(t *testing.T) {
	v := url.Values{}
	v.Set("search", "ai")
	v.Set("token", "super-secret")
	v.Set("password", "p@ss")
	out := sanitizeQuery(v)

	if !strings.Contains(out, "search=ai") {
		t.Fatalf("expected non-sensitive param preserved, got %q", out)
	}
	if strings.Contains(out, "super-secret") || strings.Contains(out, "p%40ss") || strings.Contains(out, "p@ss") {
		t.Fatalf("sensitive values must be redacted, got %q", out)
	}
	if !strings.Contains(out, "REDACTED") {
		t.Fatalf("expected REDACTED marker, got %q", out)
	}
}

// ---- domain error mapping ----

func TestWriteDomainErr_Mapping(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"invalid", fmt.Errorf("%w: bad", domain.ErrInvalid), http.StatusBadRequest},
		{"not_found", fmt.Errorf("%w: missing", domain.ErrNotFound), http.StatusNotFound},
		{"conflict", fmt.Errorf("%w: dup", domain.ErrConflict), http.StatusConflict},
		{"internal", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			writeDomainErr(rr, c.err)
			if rr.Code != c.want {
				t.Fatalf("err=%v: expected %d, got %d (body=%s)", c.err, c.want, rr.Code, rr.Body.String())
			}
		})
	}
}

// ---- admin login handler ----

func TestAdminHandler_Login_OK(t *testing.T) {
	stub := &authServiceStub{
		loginFn: func(ctx context.Context, username, password string) (string, bool, error) {
			if username == "admin" && password == "secret" {
				return "jwt-token", true, nil
			}
			return "", false, nil
		},
	}
	h := NewAdminHandler(stub)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(`{"username":"admin","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "jwt-token") {
		t.Fatalf("expected token in response, got %s", rr.Body.String())
	}
}

func TestAdminHandler_Login_BadCredentials(t *testing.T) {
	stub := &authServiceStub{
		loginFn: func(ctx context.Context, username, password string) (string, bool, error) {
			return "", false, nil
		},
	}
	h := NewAdminHandler(stub)

	req := httptest.NewRequest(http.MethodPost, "/api/admin/login", strings.NewReader(`{"username":"x","password":"y"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rr.Code, rr.Body.String())
	}
}

// ---- distributed rate limit store ----

type fakeRateLimiter struct {
	counts map[string]int64
}

func (f *fakeRateLimiter) Incr(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	if f.counts == nil {
		f.counts = map[string]int64{}
	}
	f.counts[key]++
	return f.counts[key], window, nil
}

func TestRateLimit_UsesDistributedStore(t *testing.T) {
	store := &fakeRateLimiter{}
	mw := RateLimit(RateLimitConfig{Name: "test", Limit: 2, Window: time.Minute, Store: store})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	do := func() int {
		req := httptest.NewRequest(http.MethodGet, "/x", nil)
		req.RemoteAddr = "203.0.113.5:1111"
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		return rr.Code
	}

	if c := do(); c != http.StatusOK {
		t.Fatalf("req1 expected 200, got %d", c)
	}
	if c := do(); c != http.StatusOK {
		t.Fatalf("req2 expected 200, got %d", c)
	}
	if c := do(); c != http.StatusTooManyRequests {
		t.Fatalf("req3 expected 429, got %d", c)
	}
	if got := store.counts["rl:test:203.0.113.5"]; got != 3 {
		t.Fatalf("expected store to be used with 3 increments, got %d", got)
	}
}
