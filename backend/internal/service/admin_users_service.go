package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"golang.org/x/crypto/bcrypt"
)

type AdminUsersService struct {
	repo ports.AdminUserRepository
}

func NewAdminUsersService(repo ports.AdminUserRepository) *AdminUsersService {
	return &AdminUsersService{repo: repo}
}

func (s *AdminUsersService) List(ctx context.Context) ([]domain.AdminUser, error) {
	return s.repo.List(ctx)
}

func (s *AdminUsersService) Create(ctx context.Context, cmd domain.AdminUserCreate) (string, error) {
	username := strings.TrimSpace(cmd.Username)
	if username == "" {
		return "", fmt.Errorf("%w: username is required", domain.ErrInvalid)
	}
	if len(cmd.Password) < 8 {
		return "", fmt.Errorf("%w: password must be at least 8 characters", domain.ErrInvalid)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash admin password: %w", err)
	}
	return s.repo.Create(ctx, username, string(hash))
}

func (s *AdminUsersService) SetActive(ctx context.Context, username string, active bool) (bool, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return false, fmt.Errorf("%w: username is required", domain.ErrInvalid)
	}
	return s.repo.SetActive(ctx, username, active)
}

