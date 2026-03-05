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

type AdminTrendRepo struct {
	db *pgxpool.Pool
}

func NewAdminTrendRepo(db *pgxpool.Pool) *AdminTrendRepo {
	return &AdminTrendRepo{db: db}
}

var _ ports.AdminTrendRepository = (*AdminTrendRepo)(nil)

func (r *AdminTrendRepo) List(ctx context.Context) ([]domain.AdminTrend, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, slug, name, description, image_url, order_index
		FROM trends
		ORDER BY order_index ASC, name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list trends: %w", err)
	}
	defer rows.Close()

	var out []domain.AdminTrend
	for rows.Next() {
		var it domain.AdminTrend
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.Description, &it.ImageURL, &it.OrderIndex); err != nil {
			return nil, fmt.Errorf("scan trend: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *AdminTrendRepo) Get(ctx context.Context, slug string) (domain.AdminTrend, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, slug, name, description, image_url, order_index
		FROM trends
		WHERE slug = $1
	`, strings.TrimSpace(slug))

	var it domain.AdminTrend
	if err := row.Scan(&it.ID, &it.Slug, &it.Name, &it.Description, &it.ImageURL, &it.OrderIndex); err != nil {
		if err == pgx.ErrNoRows {
			return domain.AdminTrend{}, false, nil
		}
		return domain.AdminTrend{}, false, fmt.Errorf("get trend: %w", err)
	}
	return it, true, nil
}

func (r *AdminTrendRepo) Create(ctx context.Context, cmd domain.TrendUpsert) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO trends (slug, name, description, image_url, order_index)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text
	`, strings.TrimSpace(cmd.Slug), strings.TrimSpace(cmd.Name), cmd.Description, cmd.ImageURL, cmd.Order).Scan(&id)

	if err != nil {
		return "", mapPGErr(err, "trend slug already exists")
	}
	return id, nil
}

func (r *AdminTrendRepo) Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM trends WHERE slug=$1`, slug).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find trend: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		UPDATE trends
		SET name=$2, description=$3, image_url=$4, order_index=$5, updated_at=now()
		WHERE id=$1::uuid
	`, id, strings.TrimSpace(cmd.Name), cmd.Description, cmd.ImageURL, cmd.Order)
	if err != nil {
		return "", false, fmt.Errorf("update trend: %w", err)
	}

	return id, true, nil
}

func (r *AdminTrendRepo) Delete(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM trends WHERE slug=$1`, slug)
	if err != nil {
		return false, mapPGErr(err, "trend is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func mapPGErr(err error, msg string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "23503":
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		}
	}
	return fmt.Errorf("db error: %w", err)
}
