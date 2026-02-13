package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

type TrendListItem struct {
	ID                string `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	TechnologiesCount int    `json:"technologies_count"`
}

func (r *Repo) ListTrends(ctx context.Context) ([]TrendListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			t.id::text,
			t.slug,
			t.name,
			COUNT(tech.id)::int AS technologies_count
		FROM trends t
		LEFT JOIN technologies tech ON tech.trend_id = t.id
		GROUP BY t.id
		ORDER BY t.order_index ASC, t.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query trends: %w", err)
	}
	defer rows.Close()

	var out []TrendListItem
	for rows.Next() {
		var it TrendListItem
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan trends: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

type SDGListItem struct {
	ID                string `json:"id"`
	Code              string `json:"code"`
	Title             string `json:"title"`
	TechnologiesCount int    `json:"technologies_count"`
}

func (r *Repo) ListSDGs(ctx context.Context) ([]SDGListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			s.id::text,
			s.code,
			s.title,
			COUNT(ts.technology_id)::int AS technologies_count
		FROM sdgs s
		LEFT JOIN technology_sdgs ts ON ts.sdg_id = s.id
		GROUP BY s.id
		ORDER BY s.code ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query sdgs: %w", err)
	}
	defer rows.Close()

	var out []SDGListItem
	for rows.Next() {
		var it SDGListItem
		if err := rows.Scan(&it.ID, &it.Code, &it.Title, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan sdgs: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

type TagItem struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
	Category    *string `json:"category,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (r *Repo) ListTags(ctx context.Context) ([]TagItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, slug, title, category, description
		FROM tags
		ORDER BY COALESCE(category,''), title
	`)
	if err != nil {
		return nil, fmt.Errorf("query tags: %w", err)
	}
	defer rows.Close()

	var out []TagItem
	for rows.Next() {
		var it TagItem
		if err := rows.Scan(&it.ID, &it.Slug, &it.Title, &it.Category, &it.Description); err != nil {
			return nil, fmt.Errorf("scan tags: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

type OrganizationListItem struct {
	ID                string  `json:"id"`
	Slug              string  `json:"slug"`
	Name              string  `json:"name"`
	LogoURL           *string `json:"logo_url,omitempty"`
	TechnologiesCount int     `json:"technologies_count"`
}

func (r *Repo) ListOrganizations(ctx context.Context) ([]OrganizationListItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			o.id::text,
			o.slug,
			o.name,
			o.logo_url,
			COUNT(to2.technology_id)::int AS technologies_count
		FROM organizations o
		LEFT JOIN technology_organizations to2 ON to2.organization_id = o.id
		GROUP BY o.id
		ORDER BY o.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query organizations: %w", err)
	}
	defer rows.Close()

	var out []OrganizationListItem
	for rows.Next() {
		var it OrganizationListItem
		if err := rows.Scan(&it.ID, &it.Slug, &it.Name, &it.LogoURL, &it.TechnologiesCount); err != nil {
			return nil, fmt.Errorf("scan organizations: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

type MetricDefinition struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // distance|bubble|bar
	Description *string `json:"description,omitempty"`
	Orderable   bool    `json:"orderable"`
}

func (r *Repo) ListMetrics(ctx context.Context) ([]MetricDefinition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, name, type, description, orderable
		FROM metrics_definitions
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query metrics: %w", err)
	}
	defer rows.Close()

	var out []MetricDefinition
	for rows.Next() {
		var it MetricDefinition
		if err := rows.Scan(&it.ID, &it.Name, &it.Type, &it.Description, &it.Orderable); err != nil {
			return nil, fmt.Errorf("scan metrics: %w", err)
		}
		out = append(out, it)
	}
	return out, rows.Err()
}
