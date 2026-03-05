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

type AdminOrganizationRepo struct {
	db *pgxpool.Pool
}

func NewAdminOrganizationRepo(db *pgxpool.Pool) *AdminOrganizationRepo {
	return &AdminOrganizationRepo{db: db}
}

var _ ports.AdminOrganizationRepository = (*AdminOrganizationRepo)(nil)

func (r *AdminOrganizationRepo) Create(ctx context.Context, cmd domain.OrganizationUpsert) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO organizations (slug, name, logo_url, description, website, headquarters)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text
	`,
		strings.TrimSpace(cmd.Slug),
		strings.TrimSpace(cmd.Name),
		cmd.LogoURL,
		cmd.Description,
		cmd.Website,
		cmd.Headquarters,
	).Scan(&id)

	if err != nil {
		return "", mapOrgPGErr(err, "organization slug already exists")
	}
	return id, nil
}

func (r *AdminOrganizationRepo) Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM organizations WHERE slug=$1`, slug).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find organization: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		UPDATE organizations
		SET name=$2, logo_url=$3, description=$4, website=$5, headquarters=$6, updated_at=now()
		WHERE id=$1::uuid
	`,
		id,
		strings.TrimSpace(cmd.Name),
		cmd.LogoURL,
		cmd.Description,
		cmd.Website,
		cmd.Headquarters,
	)
	if err != nil {
		return "", false, fmt.Errorf("update organization: %w", err)
	}
	return id, true, nil
}

func (r *AdminOrganizationRepo) Delete(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM organizations WHERE slug=$1`, slug)
	if err != nil {
		return false, mapOrgPGErr(err, "organization is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminOrganizationRepo) List(ctx context.Context) ([]domain.Organization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			o.id::text,
			o.slug,
			o.name,
			o.logo_url,
			o.description,
			o.website,
			o.headquarters,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		GROUP BY o.id, o.slug, o.name, o.logo_url, o.description, o.website, o.headquarters
		ORDER BY o.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	defer rows.Close()

	var out []domain.Organization
	for rows.Next() {
		var it domain.Organization
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.LogoURL, &it.Description, &it.Website, &it.Headquarters, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan organization: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *AdminOrganizationRepo) Get(ctx context.Context, slug string) (domain.Organization, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			o.id::text,
			o.slug,
			o.name,
			o.logo_url,
			o.description,
			o.website,
			o.headquarters,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		WHERE o.slug = $1
		GROUP BY o.id, o.slug, o.name, o.logo_url, o.description, o.website, o.headquarters
	`, strings.TrimSpace(slug))

	var it domain.Organization
	if err := row.Scan(&it.ID, &it.Slug, &it.Name, &it.LogoURL, &it.Description, &it.Website, &it.Headquarters, &it.TechnologiesCount); err != nil {
		if err == pgx.ErrNoRows {
			return domain.Organization{}, false, nil
		}
		return domain.Organization{}, false, fmt.Errorf("get organization: %w", err)
	}
	return it, true, nil
}

func mapOrgPGErr(err error, msg string) error {
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
