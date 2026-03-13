package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminUserService interface {
	List(ctx context.Context) ([]domain.AdminUser, error)
	Create(ctx context.Context, cmd domain.AdminUserCreate) (id string, err error)
	SetActive(ctx context.Context, username string, active bool) (ok bool, err error)
}

type AdminUserRepository interface {
	List(ctx context.Context) ([]domain.AdminUser, error)
	Create(ctx context.Context, username, passwordHash string) (id string, err error)
	SetActive(ctx context.Context, username string, active bool) (ok bool, err error)
}

