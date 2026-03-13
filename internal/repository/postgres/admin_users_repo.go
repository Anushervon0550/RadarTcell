package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminUsersRepo struct {
	db *pgxpool.Pool
}

func NewAdminUsersRepo(db *pgxpool.Pool) *AdminUsersRepo {
	return &AdminUsersRepo{db: db}
}

var _ ports.AdminUserRepository = (*AdminUsersRepo)(nil)

func (r *AdminUsersRepo) List(ctx context.Context) ([]domain.AdminUser, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, username, is_active, created_at
		FROM admin_users
		ORDER BY username ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list admin users: %w", err)
	}
	defer rows.Close()

	var out []domain.AdminUser
	for rows.Next() {
		var it domain.AdminUser
		if err := rows.Scan(&it.ID, &it.Username, &it.IsActive, &it.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan admin user: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *AdminUsersRepo) Create(ctx context.Context, username, passwordHash string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO admin_users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id::text
	`, strings.TrimSpace(username), passwordHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%w: admin username already exists", domain.ErrConflict)
		}
		return "", fmt.Errorf("create admin user: %w", err)
	}
	return id, nil
}

func (r *AdminUsersRepo) SetActive(ctx context.Context, username string, active bool) (bool, error) {
	username = strings.TrimSpace(username)

	if active {
		ct, err := r.db.Exec(ctx, `
			UPDATE admin_users
			SET is_active = TRUE, updated_at = now()
			WHERE username = $1
		`, username)
		if err != nil {
			return false, fmt.Errorf("set admin user active: %w", err)
		}
		return ct.RowsAffected() > 0, nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `LOCK TABLE admin_users IN SHARE ROW EXCLUSIVE MODE`); err != nil {
		return false, fmt.Errorf("lock admin users: %w", err)
	}

	var isActive bool
	err = tx.QueryRow(ctx, `SELECT is_active FROM admin_users WHERE username = $1`, username).Scan(&isActive)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("find admin user: %w", err)
	}

	if isActive {
		var activeCount int
		if err := tx.QueryRow(ctx, `SELECT COUNT(*)::int FROM admin_users WHERE is_active = TRUE`).Scan(&activeCount); err != nil {
			return false, fmt.Errorf("count active admins: %w", err)
		}
		if activeCount <= 1 {
			return false, fmt.Errorf("%w: cannot deactivate last active admin", domain.ErrConflict)
		}
	}

	if _, err := tx.Exec(ctx, `
		UPDATE admin_users
		SET is_active = FALSE, updated_at = now()
		WHERE username = $1
	`, username); err != nil {
		return false, fmt.Errorf("deactivate admin user: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("commit: %w", err)
	}
	return true, nil
}

