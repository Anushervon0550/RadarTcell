package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{db: db}
}

var _ ports.AuthRepository = (*AuthRepo)(nil)

func (r *AuthRepo) GetAdminPasswordHash(ctx context.Context, username string) (string, bool, error) {
	username = strings.TrimSpace(username)
	var hash string
	err := r.db.QueryRow(ctx, `
		SELECT password_hash
		FROM admin_users
		WHERE username = $1 AND is_active = TRUE
	`, username).Scan(&hash)
	if err == nil {
		return hash, true, nil
	}
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	return "", false, fmt.Errorf("get admin password hash: %w", err)
}

