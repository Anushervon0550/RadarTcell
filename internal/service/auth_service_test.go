package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type fakeAuthRepo struct {
	getHashFn func(ctx context.Context, username string) (string, bool, error)
}

func (f *fakeAuthRepo) GetAdminPasswordHash(ctx context.Context, username string) (string, bool, error) {
	if f.getHashFn != nil {
		return f.getHashFn(ctx, username)
	}
	return "", false, nil
}

func TestAuthService_Login_DBHash_Success(t *testing.T) {
	h, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc, err := NewAuthService(&fakeAuthRepo{
		getHashFn: func(ctx context.Context, username string) (string, bool, error) {
			if username != "alice" {
				t.Fatalf("expected username alice, got %q", username)
			}
			return string(h), true, nil
		},
	}, "", "", "jwt-secret", time.Hour)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	tok, ok, err := svc.Login(context.Background(), "alice", "secret123")
	if err != nil {
		t.Fatalf("login err: %v", err)
	}
	if !ok || tok == "" {
		t.Fatalf("expected successful login, got ok=%v token=%q", ok, tok)
	}
}

func TestAuthService_Login_DBHash_BadPassword(t *testing.T) {
	h, err := bcrypt.GenerateFromPassword([]byte("good-pass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc, err := NewAuthService(&fakeAuthRepo{
		getHashFn: func(ctx context.Context, username string) (string, bool, error) {
			return string(h), true, nil
		},
	}, "", "", "jwt-secret", time.Hour)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	tok, ok, err := svc.Login(context.Background(), "alice", "bad-pass")
	if err != nil {
		t.Fatalf("login err: %v", err)
	}
	if ok || tok != "" {
		t.Fatalf("expected failed login, got ok=%v token=%q", ok, tok)
	}
}

func TestAuthService_Login_DBRepoError(t *testing.T) {
	expected := errors.New("db down")
	svc, err := NewAuthService(&fakeAuthRepo{
		getHashFn: func(ctx context.Context, username string) (string, bool, error) {
			return "", false, expected
		},
	}, "", "", "jwt-secret", time.Hour)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	_, ok, err := svc.Login(context.Background(), "alice", "x")
	if ok {
		t.Fatal("expected ok=false")
	}
	if !errors.Is(err, expected) {
		t.Fatalf("expected err %v, got %v", expected, err)
	}
}

func TestAuthService_Login_ENVFallback(t *testing.T) {
	svc, err := NewAuthService(&fakeAuthRepo{
		getHashFn: func(ctx context.Context, username string) (string, bool, error) {
			return "", false, nil
		},
	}, "root", "root-pass", "jwt-secret", time.Hour)
	if err != nil {
		t.Fatalf("new auth service: %v", err)
	}

	tok, ok, err := svc.Login(context.Background(), "root", "root-pass")
	if err != nil {
		t.Fatalf("login err: %v", err)
	}
	if !ok || tok == "" {
		t.Fatalf("expected fallback login success, got ok=%v token=%q", ok, tok)
	}
}

