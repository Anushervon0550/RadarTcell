package ports

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (token string, ok bool, err error)
	Verify(ctx context.Context, token string) (subject string, ok bool, err error)
}

type AuthRepository interface {
	GetAdminPasswordHash(ctx context.Context, username string) (passwordHash string, ok bool, err error)
}

