package postgres

import (
	"context"
	"fmt"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepo struct {
	db *pgxpool.Pool
}

func NewCatalogRepo(db *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{db: db}
}

func (r *CatalogRepo) ListTrends(ctx context.Context) ([]domain.Trend, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			t.id::text,
			t.slug,
			t.name,
			COUNT(tech.id)::int AS technologies_count
		FROM trends t
		LEFT JOIN technologies tech ON tech.trend_id = t.id
		GROUP BY t.id, t.slug, t.name
		ORDER BY t.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list trends: %w", err)
	}
	defer rows.Close()

	var out []domain.Trend
	for rows.Next() {
		var it domain.Trend
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan trend: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *CatalogRepo) ListSDGs(ctx context.Context) ([]domain.SDG, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			s.id::text,
			s.code,
			s.title,
			COUNT(ts.technology_id)::int AS technologies_count
		FROM sdgs s
		LEFT JOIN technology_sdgs ts ON ts.sdg_id = s.id
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

func (r *CatalogRepo) ListTags(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id::text,
			slug,
			title,
			category,
			description
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

func (r *CatalogRepo) ListOrganizations(ctx context.Context) ([]domain.Organization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			o.id::text,
			o.slug,
			o.name,
			o.logo_url,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		GROUP BY o.id, o.slug, o.name, o.logo_url
		ORDER BY o.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	defer rows.Close()

	var out []domain.Organization
	for rows.Next() {
		var it domain.Organization
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.LogoURL, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan org: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *CatalogRepo) ListMetrics(ctx context.Context) ([]domain.MetricDefinition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id::text,
			name,
			type,
			description,
			orderable
		FROM metrics_definitions
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}
	defer rows.Close()

	var out []domain.MetricDefinition
	for rows.Next() {
		var it domain.MetricDefinition
		if err := rows.Scan(&it.ID, &it.Name, &it.Type, &it.Description, &it.Orderable); err != nil {
			return nil, fmt.Errorf("scan metric: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
func (r *CatalogRepo) GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			o.id::text,
			o.slug,
			o.name,
			o.logo_url,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		WHERE o.slug = $1
		GROUP BY o.id, o.slug, o.name, o.logo_url
	`, slug)

	var it domain.Organization
	if err := row.Scan(&it.ID, &it.Slug, &it.Name, &it.LogoURL, &it.TechnologiesCount); err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Organization{}, false, nil
		}
		return domain.Organization{}, false, fmt.Errorf("get organization: %w", err)
	}
	return it, true, nil
}
