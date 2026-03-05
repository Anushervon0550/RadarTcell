package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TechnologyRepo struct {
	db *pgxpool.Pool
}

func NewTechnologyRepo(db *pgxpool.Pool) *TechnologyRepo {
	return &TechnologyRepo{db: db}
}

func (r *TechnologyRepo) ListTrendIDsOrdered(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text
		FROM trends
		ORDER BY order_index ASC, name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list trend ids: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan trend id: %w", err)
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

func (r *TechnologyRepo) ListTechnologies(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 200 {
		p.Limit = 200
	}
	offset := (p.Page - 1) * p.Limit

	where, args := buildTechWhere(p, 0)

	sortExpr, err := r.normalizeSortBy(ctx, p.SortBy)
	if err != nil {
		return nil, 0, err
	}
	orderDir := normalizeOrder(p.Order)

	// Count
	var total int
	countSQL := `
		SELECT COUNT(*)
		FROM technologies tech
		JOIN trends tr ON tr.id = tech.trend_id
	` + where
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count technologies: %w", err)
	}

	nameExpr := "tech.name"
	descShortExpr := "tech.description_short"
	descFullExpr := "tech.description_full"
	joinI18n := ""
	listArgs := args
	if p.Locale != "" {
		nameExpr = "COALESCE(ti.name, tech.name)"
		descShortExpr = "COALESCE(ti.description_short, tech.description_short)"
		descFullExpr = "COALESCE(ti.description_full, tech.description_full)"
		joinI18n = fmt.Sprintf("LEFT JOIN technology_i18n ti ON ti.technology_id = tech.id AND ti.locale = $%d", len(args)+1)
		listArgs = append(listArgs, p.Locale)
	}

	// List
	listSQL := fmt.Sprintf(`
		SELECT
			tech.id::text,
			tech.slug,
			tech.list_index,
			%s,
			%s,
			%s,
			tech.readiness_level,
			tech.custom_metric_1,
			tech.custom_metric_2,
			tech.custom_metric_3,
			tech.custom_metric_4,
			tech.image_url,
			tech.source_link,
			tech.trend_id::text,
			tr.slug,
			tr.name
		FROM technologies tech
		JOIN trends tr ON tr.id = tech.trend_id
		%s
		%s
		ORDER BY %s %s NULLS LAST, tech.list_index ASC
		LIMIT %d OFFSET %d
	`, nameExpr, descShortExpr, descFullExpr, joinI18n, where, sortExpr, orderDir, p.Limit, offset)

	rows, err := r.db.Query(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list technologies: %w", err)
	}
	defer rows.Close()

	var out []domain.Technology
	for rows.Next() {
		var t domain.Technology
		if err := rows.Scan(
			&t.ID,
			&t.Slug,
			&t.Index,
			&t.Name,
			&t.DescriptionShort,
			&t.DescriptionFull,
			&t.TRL,
			&t.CustomMetric1,
			&t.CustomMetric2,
			&t.CustomMetric3,
			&t.CustomMetric4,
			&t.ImageURL,
			&t.SourceLink,
			&t.TrendID,
			&t.TrendSlug,
			&t.TrendName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan technology: %w", err)
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows technologies: %w", err)
	}

	return out, total, nil
}

func buildTechWhere(p domain.TechnologyListParams, startIndex int) (string, []any) {
	var b strings.Builder
	args := make([]any, 0, 10)

	b.WriteString(" WHERE 1=1 ")

	if s := strings.TrimSpace(p.Search); s != "" {
		args = append(args, s)
		n := startIndex + len(args)
		b.WriteString(fmt.Sprintf(
			" AND (tech.name ILIKE '%%' || $%d || '%%' OR tech.slug ILIKE '%%' || $%d || '%%') ",
			n, n,
		))
	}

	if p.TrendID != "" {
		args = append(args, p.TrendID)
		b.WriteString(fmt.Sprintf(" AND tech.trend_id = $%d::uuid ", startIndex+len(args)))
	}
	if p.SDGID != "" {
		args = append(args, p.SDGID)
		b.WriteString(fmt.Sprintf(
			" AND EXISTS (SELECT 1 FROM technology_sdgs ts WHERE ts.technology_id = tech.id AND ts.sdg_id = $%d::uuid) ",
			startIndex+len(args),
		))
	}
	if p.TagID != "" {
		args = append(args, p.TagID)
		b.WriteString(fmt.Sprintf(
			" AND EXISTS (SELECT 1 FROM technology_tags tt WHERE tt.technology_id = tech.id AND tt.tag_id = $%d::uuid) ",
			startIndex+len(args),
		))
	}
	if p.OrganizationID != "" {
		args = append(args, p.OrganizationID)
		b.WriteString(fmt.Sprintf(
			" AND EXISTS (SELECT 1 FROM technology_organizations to2 WHERE to2.technology_id = tech.id AND to2.organization_id = $%d::uuid) ",
			startIndex+len(args),
		))
	}

	if p.HasTRLMin {
		args = append(args, p.TRLMin)
		b.WriteString(fmt.Sprintf(" AND tech.readiness_level >= $%d ", startIndex+len(args)))
	}
	if p.HasTRLMax {
		args = append(args, p.TRLMax)
		b.WriteString(fmt.Sprintf(" AND tech.readiness_level <= $%d ", startIndex+len(args)))
	}

	if len(p.OnlyIDs) > 0 {
		args = append(args, p.OnlyIDs) // []string
		b.WriteString(fmt.Sprintf(" AND tech.id::text = ANY($%d::text[]) ", startIndex+len(args)))
	}

	return b.String(), args
}

func normalizeOrder(order string) string {
	switch strings.ToLower(strings.TrimSpace(order)) {
	case "desc":
		return "DESC"
	default:
		return "ASC"
	}
}

// ------ lookups by slug/code ------

func (r *TechnologyRepo) GetTrendIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM trends WHERE slug=$1`, slug).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	return "", false, fmt.Errorf("get trend id: %w", err)
}

func (r *TechnologyRepo) GetSDGIDByCode(ctx context.Context, code string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM sdgs WHERE code=$1`, code).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	return "", false, fmt.Errorf("get sdg id: %w", err)
}

func (r *TechnologyRepo) GetTagIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM tags WHERE slug=$1`, slug).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	return "", false, fmt.Errorf("get tag id: %w", err)
}

func (r *TechnologyRepo) GetOrganizationIDBySlug(ctx context.Context, slug string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id::text FROM organizations WHERE slug=$1`, slug).Scan(&id)
	if err == nil {
		return id, true, nil
	}
	if err == pgx.ErrNoRows {
		return "", false, nil
	}
	return "", false, fmt.Errorf("get org id: %w", err)
}

// ------ ids by relation (for highlight + future group endpoints) ------

func (r *TechnologyRepo) ListTechnologyIDsByTrendID(ctx context.Context, trendID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `SELECT id::text FROM technologies WHERE trend_id=$1::uuid`, trendID)
	if err != nil {
		return nil, fmt.Errorf("list tech ids by trend: %w", err)
	}
	defer rows.Close()
	return scanIDs(rows)
}

func (r *TechnologyRepo) ListTechnologyIDsBySDGID(ctx context.Context, sdgID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT technology_id::text
		FROM technology_sdgs
		WHERE sdg_id=$1::uuid
	`, sdgID)
	if err != nil {
		return nil, fmt.Errorf("list tech ids by sdg: %w", err)
	}
	defer rows.Close()
	return scanIDs(rows)
}

func (r *TechnologyRepo) ListTechnologyIDsByTagID(ctx context.Context, tagID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT technology_id::text
		FROM technology_tags
		WHERE tag_id=$1::uuid
	`, tagID)
	if err != nil {
		return nil, fmt.Errorf("list tech ids by tag: %w", err)
	}
	defer rows.Close()
	return scanIDs(rows)
}

func (r *TechnologyRepo) ListTechnologyIDsByOrganizationID(ctx context.Context, orgID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT technology_id::text
		FROM technology_organizations
		WHERE organization_id=$1::uuid
	`, orgID)
	if err != nil {
		return nil, fmt.Errorf("list tech ids by org: %w", err)
	}
	defer rows.Close()
	return scanIDs(rows)
}

func scanIDs(rows pgx.Rows) ([]string, error) {
	var out []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan id: %w", err)
		}
		out = append(out, id)
	}
	return out, rows.Err()
}

// ------ card (пока просто база; расширим на следующем шаге) ------

func (r *TechnologyRepo) GetTechnologyBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error) {
	locale = strings.TrimSpace(locale)
	nameExpr := "tech.name"
	descShortExpr := "tech.description_short"
	descFullExpr := "tech.description_full"
	joinI18n := ""
	args := []any{slug}
	if locale != "" {
		nameExpr = "COALESCE(ti.name, tech.name)"
		descShortExpr = "COALESCE(ti.description_short, tech.description_short)"
		descFullExpr = "COALESCE(ti.description_full, tech.description_full)"
		joinI18n = "LEFT JOIN technology_i18n ti ON ti.technology_id = tech.id AND ti.locale = $2"
		args = append(args, locale)
	}

	row := r.db.QueryRow(ctx, fmt.Sprintf(`
		SELECT
			tech.id::text,
			tech.slug,
			tech.list_index,
			%s,
			%s,
			%s,
			tech.readiness_level,
			tech.custom_metric_1,
			tech.custom_metric_2,
			tech.custom_metric_3,
			tech.custom_metric_4,
			tech.image_url,
			tech.source_link,
			tech.trend_id::text,
			tr.slug,
			tr.name
		FROM technologies tech
		JOIN trends tr ON tr.id = tech.trend_id
		%s
		WHERE tech.slug = $1
	`, nameExpr, descShortExpr, descFullExpr, joinI18n), args...)

	var t domain.Technology
	if err := row.Scan(
		&t.ID,
		&t.Slug,
		&t.Index,
		&t.Name,
		&t.DescriptionShort,
		&t.DescriptionFull,
		&t.TRL,
		&t.CustomMetric1,
		&t.CustomMetric2,
		&t.CustomMetric3,
		&t.CustomMetric4,
		&t.ImageURL,
		&t.SourceLink,
		&t.TrendID,
		&t.TrendSlug,
		&t.TrendName,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get technology: %w", err)
	}

	return &t, true, nil
}

func (r *TechnologyRepo) ListTagsByTechnologyID(ctx context.Context, techID string) ([]domain.Tag, error) {
	rows, err := r.db.Query(ctx, `
		SELECT t.id::text, t.slug, t.title, t.category, t.description
		FROM tags t
		JOIN technology_tags tt ON tt.tag_id = t.id
		WHERE tt.technology_id = $1::uuid
		ORDER BY COALESCE(t.category,''), t.title ASC
	`, techID)
	if err != nil {
		return nil, fmt.Errorf("list tags by tech: %w", err)
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

func (r *TechnologyRepo) ListSDGsByTechnologyID(ctx context.Context, techID string) ([]domain.SDG, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s.id::text, s.code, s.title, 0::int
		FROM sdgs s
		JOIN technology_sdgs ts ON ts.sdg_id = s.id
		WHERE ts.technology_id = $1::uuid
		ORDER BY s.code ASC
	`, techID)
	if err != nil {
		return nil, fmt.Errorf("list sdgs by tech: %w", err)
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

func (r *TechnologyRepo) ListOrganizationsByTechnologyID(ctx context.Context, techID string) ([]domain.Organization, error) {
	rows, err := r.db.Query(ctx, `
		SELECT o.id::text, o.slug, o.name, o.logo_url, 0::int
		FROM organizations o
		JOIN technology_organizations to2 ON to2.organization_id = o.id
		WHERE to2.technology_id = $1::uuid
		ORDER BY o.name ASC
	`, techID)
	if err != nil {
		return nil, fmt.Errorf("list orgs by tech: %w", err)
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

func (r *TechnologyRepo) normalizeSortBy(ctx context.Context, sortBy string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(sortBy))
	switch s {
	case "", "index", "list_index":
		return "tech.list_index", nil
	case "name":
		return "tech.name", nil
	case "trl", "readiness_level":
		return "tech.readiness_level", nil
	case "custom_metric_1":
		return "tech.custom_metric_1", nil
	case "custom_metric_2":
		return "tech.custom_metric_2", nil
	case "custom_metric_3":
		return "tech.custom_metric_3", nil
	case "custom_metric_4":
		return "tech.custom_metric_4", nil
	case "trend":
		return "tr.name", nil
	}

	allowed, err := r.listOrderableFieldKeys(ctx)
	if err != nil {
		return "", err
	}
	if _, ok := allowed[s]; ok {
		return "tech." + s, nil
	}

	return "", fmt.Errorf("%w: sort_by must be name or any orderable metric field_key", domain.ErrInvalid)
}

func (r *TechnologyRepo) listOrderableFieldKeys(ctx context.Context) (map[string]struct{}, error) {
	rows, err := r.db.Query(ctx, `
		SELECT field_key
		FROM metrics_definitions
		WHERE orderable = TRUE AND field_key IS NOT NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("list orderable field keys: %w", err)
	}
	defer rows.Close()

	out := map[string]struct{}{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, fmt.Errorf("scan orderable field key: %w", err)
		}
		key = strings.ToLower(strings.TrimSpace(key))
		if key != "" {
			out[key] = struct{}{}
		}
	}
	return out, rows.Err()
}
