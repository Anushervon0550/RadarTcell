package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type adminUsersRepoStub struct {
	ports.AdminUserRepository

	setActiveFn func(ctx context.Context, username string, active bool) (bool, error)
	createFn    func(ctx context.Context, username, passwordHash string) (string, error)

	gotSetActiveUsername string
	gotSetActiveValue    bool
}

func (s *adminUsersRepoStub) SetActive(ctx context.Context, username string, active bool) (bool, error) {
	s.gotSetActiveUsername = username
	s.gotSetActiveValue = active
	if s.setActiveFn != nil {
		return s.setActiveFn(ctx, username, active)
	}
	return false, nil
}

func (s *adminUsersRepoStub) Create(ctx context.Context, username, passwordHash string) (string, error) {
	if s.createFn != nil {
		return s.createFn(ctx, username, passwordHash)
	}
	return "", nil
}

func TestAdminUsersService_SetActive_EmptyUsername(t *testing.T) {
	svc := NewAdminUsersService(&adminUsersRepoStub{})

	ok, err := svc.SetActive(context.Background(), "   ", false)
	if ok {
		t.Fatal("expected ok=false")
	}
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected domain.ErrInvalid, got %v", err)
	}
}

func TestAdminUsersService_SetActive_ConflictPassthrough(t *testing.T) {
	expected := errors.New("sentinel")
	repo := &adminUsersRepoStub{
		setActiveFn: func(ctx context.Context, username string, active bool) (bool, error) {
			return false, errors.Join(domain.ErrConflict, expected)
		},
	}
	svc := NewAdminUsersService(repo)

	ok, err := svc.SetActive(context.Background(), "alice", false)
	if ok {
		t.Fatal("expected ok=false")
	}
	if !errors.Is(err, domain.ErrConflict) {
		t.Fatalf("expected domain.ErrConflict, got %v", err)
	}
	if !errors.Is(err, expected) {
		t.Fatalf("expected wrapped sentinel error, got %v", err)
	}
}

func TestAdminUsersService_SetActive_CallsRepo(t *testing.T) {
	repo := &adminUsersRepoStub{
		setActiveFn: func(ctx context.Context, username string, active bool) (bool, error) {
			return true, nil
		},
	}
	svc := NewAdminUsersService(repo)

	ok, err := svc.SetActive(context.Background(), "alice", true)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true")
	}
	if repo.gotSetActiveUsername != "alice" || !repo.gotSetActiveValue {
		t.Fatalf("unexpected call args: username=%q active=%v", repo.gotSetActiveUsername, repo.gotSetActiveValue)
	}
}

