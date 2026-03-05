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

type AdminTagRepo struct {
	db *pgxpool.Pool
}

func NewAdminTagRepo(db *pgxpool.Pool) *AdminTagRepo {
	return &AdminTagRepo{db: db}
}

var _ ports.AdminTagRepository = (*AdminTagRepo)(nil)

func (r *AdminTagRepo) Create(ctx context.Context, cmd domain.TagUpsert) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO tags (slug, title, category, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`, strings.TrimSpace(cmd.Slug), strings.TrimSpace(cmd.Title), strings.TrimSpace(cmd.Category), cmd.Description).Scan(&id)

	if err != nil {
		return "", mapPGErrTag(err, "tag slug already exists")
	}
	return id, nil
}

func (r *AdminTagRepo) Update(ctx context.Context, slug string, cmd domain.TagUpsert) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM tags WHERE slug=$1`, slug).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find tag: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		UPDATE tags
		SET title=$2, category=$3, description=$4, updated_at=now()
		WHERE id=$1::uuid
	`, id, strings.TrimSpace(cmd.Title), strings.TrimSpace(cmd.Category), cmd.Description)
	if err != nil {
		return "", false, fmt.Errorf("update tag: %w", err)
	}

	return id, true, nil
}

func (r *AdminTagRepo) Delete(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM tags WHERE slug=$1`, slug)
	if err != nil {
		return false, mapPGErrTag(err, "tag is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminTagRepo) List(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, slug, title, category, description
		FROM tags
		ORDER BY COALESCE(category,''), title ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer rows.Close()

	var out []domain.Tag
	for rows.Next() {
		var it domain.Tag
		if err := rows.Scan(&it.ID, &it.Slug, &it.Title, &it.Category, &it.Description); err != nil {
			return nil, fmt.Errorf("scan tag: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *AdminTagRepo) Get(ctx context.Context, slug string) (domain.Tag, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, slug, title, category, description
		FROM tags
		WHERE slug = $1
	`, strings.TrimSpace(slug))

	var it domain.Tag
	if err := row.Scan(&it.ID, &it.Slug, &it.Title, &it.Category, &it.Description); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Tag{}, false, nil
		}
		return domain.Tag{}, false, fmt.Errorf("get tag: %w", err)
	}
	return it, true, nil
}

func mapPGErrTag(err error, msg string) error {
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
