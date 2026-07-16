package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// =============================== PUBLIC =====================================

func (r *CatalogRepo) ListSDGs(ctx context.Context, locale string) ([]domain.SDG, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			s.id::text,
			s.code,
			s.title,
			COUNT(tech.id)::int AS technologies_count
		FROM sdgs s
		LEFT JOIN technology_sdgs ts ON ts.sdg_id = s.id
		LEFT JOIN technologies tech ON tech.id = ts.technology_id AND tech.deleted_at IS NULL
		WHERE s.deleted_at IS NULL
		GROUP BY s.id, s.code, s.title
		ORDER BY s.code ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list sdgs: %w", err)
	}
	defer rows.Close()

	var out []domain.SDG
	for rows.Next() {
		var it domain.SDG
		if err := rows.Scan(&it.ID, &it.Code, &it.Title, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan sdg: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// =============================== ADMIN ======================================

type AdminSDGRepo struct {
	db *pgxpool.Pool
}

func NewAdminSDGRepo(db *pgxpool.Pool) *AdminSDGRepo {
	return &AdminSDGRepo{db: db}
}

var _ ports.AdminSDGRepository = (*AdminSDGRepo)(nil)

func (r *AdminSDGRepo) Create(ctx context.Context, cmd domain.SDGUpsert) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO sdgs (code, title, description, icon)
		VALUES ($1,$2,$3,$4)
		RETURNING id
	`,
		strings.TrimSpace(cmd.Code),
		strings.TrimSpace(cmd.Title),
		cmd.Description,
		cmd.Icon,
	).Scan(&id)
	if err != nil {
		return "", mapPGErr(err, "sdg already exists")
	}
	return id, nil
}

func (r *AdminSDGRepo) Update(ctx context.Context, code string, cmd domain.SDGUpsert) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE sdgs
		SET title=$2, description=$3, icon=$4, updated_at=now()
		WHERE code=$1 AND deleted_at IS NULL
	`,
		strings.TrimSpace(code),
		strings.TrimSpace(cmd.Title),
		cmd.Description,
		cmd.Icon,
	)
	if err != nil {
		return false, mapPGErr(err, "sdg conflict")
	}
	return tag.RowsAffected() > 0, nil
}

func (r *AdminSDGRepo) Delete(ctx context.Context, code string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		UPDATE sdgs
		SET deleted_at = now(), updated_at = now()
		WHERE code=$1 AND deleted_at IS NULL
	`, strings.TrimSpace(code))
	if err != nil {
		return false, mapPGErr(err, "sdg is referenced")
	}
	return tag.RowsAffected() > 0, nil
}

func (r *AdminSDGRepo) List(ctx context.Context) ([]domain.AdminSDG, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, code, title, description, icon
		FROM sdgs
		WHERE deleted_at IS NULL
		ORDER BY code ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list sdgs: %w", err)
	}
	defer rows.Close()

	var out []domain.AdminSDG
	for rows.Next() {
		var it domain.AdminSDG
		if err := rows.Scan(&it.ID, &it.Code, &it.Title, &it.Description, &it.Icon); err != nil {
			return nil, fmt.Errorf("scan sdg: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *AdminSDGRepo) Get(ctx context.Context, code string) (domain.AdminSDG, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, code, title, description, icon
		FROM sdgs
		WHERE code = $1 AND deleted_at IS NULL
	`, strings.TrimSpace(code))

	var it domain.AdminSDG
	if err := row.Scan(&it.ID, &it.Code, &it.Title, &it.Description, &it.Icon); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AdminSDG{}, false, nil
		}
		return domain.AdminSDG{}, false, fmt.Errorf("get sdg: %w", err)
	}
	return it, true, nil
}
