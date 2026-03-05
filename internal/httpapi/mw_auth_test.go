package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type authServiceStub struct {
	verifyFn func(ctx context.Context, token string) (string, bool, error)
	loginFn  func(ctx context.Context, username, password string) (string, bool, error)

	gotVerifyToken string
	verifyCalls    int
}

func (s *authServiceStub) Login(ctx context.Context, username, password string) (string, bool, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, username, password)
	}
	return "", false, nil
}

func (s *authServiceStub) Verify(ctx context.Context, token string) (string, bool, error) {
	s.gotVerifyToken = token
	s.verifyCalls++
	if s.verifyFn != nil {
		return s.verifyFn(ctx, token)
	}
	return "", false, nil
}

func TestAuthRequired_MissingBearerToken(t *testing.T) {
	stub := &authServiceStub{}
	mw := AuthRequired(stub)

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
	if stub.verifyCalls != 0 {
		t.Fatalf("expected verify not called, got %d", stub.verifyCalls)
	}
	if nextCalled {
		t.Fatal("next must not be called")
	}
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	stub := &authServiceStub{
		verifyFn: func(ctx context.Context, token string) (string, bool, error) {
			return "", false, nil
		},
	}
	mw := AuthRequired(stub)

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), "invalid token") {
		t.Fatalf("expected invalid token message, got: %s", rr.Body.String())
	}
	if nextCalled {
		t.Fatal("next must not be called")
	}
}

func TestAuthRequired_VerifyError_InternalServerError(t *testing.T) {
	stub := &authServiceStub{
		verifyFn: func(ctx context.Context, token string) (string, bool, error) {
			return "", false, errors.New("db connection failed")
		},
	}
	mw := AuthRequired(stub)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), internalErrorMessage) {
		t.Fatalf("expected safe internal error message, got: %s", rr.Body.String())
	}
}

func TestAuthRequired_Success_SetsAdminSubjectAndTrimsToken(t *testing.T) {
	stub := &authServiceStub{
		verifyFn: func(ctx context.Context, token string) (string, bool, error) {
			if token != "good-token" {
				t.Fatalf("expected token good-token, got %q", token)
			}
			return "admin-user", true, nil
		},
	}
	mw := AuthRequired(stub)

	nextCalled := false
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		if got := AdminSubject(r); got != "admin-user" {
			t.Fatalf("expected admin subject admin-user, got %q", got)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/admin/me", nil)
	req.Header.Set("Authorization", "Bearer   good-token   ")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusNoContent, rr.Code, rr.Body.String())
	}
	if !nextCalled {
		t.Fatal("expected next to be called")
	}
}
