package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepo struct {
	db *pgxpool.Pool
}

func NewCatalogRepo(db *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{db: db}
}

func (r *CatalogRepo) ListTrends(ctx context.Context, locale string) ([]domain.Trend, error) {
	locale = strings.TrimSpace(locale)
	args := []any{}
	nameExpr := "t.name"
	join := ""
	if locale != "" {
		args = append(args, locale)
		nameExpr = "COALESCE(ti.name, t.name)"
		join = "LEFT JOIN trend_i18n ti ON ti.trend_id = t.id AND ti.locale = $1"
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT
			t.id::text,
			t.slug,
			%s,
			COUNT(tech.id)::int AS technologies_count
		FROM trends t
		LEFT JOIN technologies tech ON tech.trend_id = t.id
		%s
		GROUP BY t.id, t.slug, %s
		ORDER BY %s ASC
	`, nameExpr, join, nameExpr, nameExpr), args...)
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

func (r *CatalogRepo) ListSDGs(ctx context.Context, locale string) ([]domain.SDG, error) {
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

func (r *CatalogRepo) ListTags(ctx context.Context, locale string) ([]domain.Tag, error) {
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

func (r *CatalogRepo) ListOrganizations(ctx context.Context, locale string) ([]domain.Organization, error) {
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

func (r *CatalogRepo) ListMetrics(ctx context.Context, locale string) ([]domain.MetricDefinition, error) {
	locale = strings.TrimSpace(locale)
	args := []any{}
	nameExpr := "m.name"
	descExpr := "m.description"
	join := ""
	if locale != "" {
		args = append(args, locale)
		nameExpr = "COALESCE(mi.name, m.name)"
		descExpr = "COALESCE(mi.description, m.description)"
		join = "LEFT JOIN metric_definition_i18n mi ON mi.metric_id = m.id AND mi.locale = $1"
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT
			m.id::text,
			%s,
			m.type,
			%s,
			m.orderable,
			m.field_key
		FROM metrics_definitions m
		%s
		ORDER BY %s ASC
	`, nameExpr, descExpr, join, nameExpr), args...)
	if err != nil {
		return nil, fmt.Errorf("list metrics: %w", err)
	}
	defer rows.Close()

	var out []domain.MetricDefinition
	for rows.Next() {
		var it domain.MetricDefinition
		if err := rows.Scan(&it.ID, &it.Name, &it.Type, &it.Description, &it.Orderable, &it.FieldKey); err != nil {
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
			o.description,
			o.website,
			o.headquarters,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		WHERE o.slug = $1
		GROUP BY o.id, o.slug, o.name, o.logo_url, o.description, o.website, o.headquarters
	`, slug)

	var it domain.Organization
	if err := row.Scan(
		&it.ID,
		&it.Slug,
		&it.Name,
		&it.LogoURL,
		&it.Description,
		&it.Website,
		&it.Headquarters,
		&it.TechnologiesCount,
	); err != nil {
		if err.Error() == "no rows in result set" {
			return domain.Organization{}, false, nil
		}
		return domain.Organization{}, false, fmt.Errorf("get organization: %w", err)
	}
	return it, true, nil
}
func (r *CatalogRepo) GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error) {
	const q = `
SELECT
	m.id,
	m.name,
	m.type,
	m.field_key,
	t.id,
	CASE m.field_key
		WHEN 'readiness_level' THEN t.readiness_level::double precision
		WHEN 'list_index' THEN t.list_index::double precision
		WHEN 'custom_metric_1' THEN t.custom_metric_1
		WHEN 'custom_metric_2' THEN t.custom_metric_2
		WHEN 'custom_metric_3' THEN t.custom_metric_3
		WHEN 'custom_metric_4' THEN t.custom_metric_4
		ELSE NULL
	END AS value
FROM metrics_definitions m
JOIN technologies t ON t.id = $2
WHERE m.id = $1
LIMIT 1;
`
	var mid, mname, mtype string
	var fieldKey *string
	var tid string
	var value sql.NullFloat64

	err := r.db.QueryRow(ctx, q, metricID, technologyID).
		Scan(&mid, &mname, &mtype, &fieldKey, &tid, &value)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var v any = nil
	if value.Valid {
		v = value.Float64
	}

	return map[string]any{
		"metric_id":     mid,
		"metric_name":   mname,
		"type":          mtype,
		"field_key":     fieldKey,
		"technology_id": tid,
		"value":         v,
	}, true, nil
}
