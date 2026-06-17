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

type AdminTechnologyRepo struct {
	db *pgxpool.Pool
}

func NewAdminTechnologyRepo(db *pgxpool.Pool) *AdminTechnologyRepo {
	return &AdminTechnologyRepo{db: db}
}

var _ ports.AdminTechnologyRepository = (*AdminTechnologyRepo)(nil)

func (r *AdminTechnologyRepo) Create(ctx context.Context, cmd domain.TechnologyUpsert) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	trendID, err := resolveTrendID(ctx, tx, cmd.TrendSlug)
	if err != nil {
		return "", err
	}

	var techID string
	err = tx.QueryRow(ctx, `
		INSERT INTO technologies (
			slug, list_index, name,
			description_short, description_full,
			readiness_level,
			custom_metric_1, custom_metric_2, custom_metric_3, custom_metric_4,
			image_url, source_link,
			trend_id
		)
		VALUES (
			$1, $2, $3,
			$4, $5,
			$6,
			$7, $8, $9, $10,
			$11, $12,
			$13::uuid
		)
		RETURNING id::text
	`,
		strings.TrimSpace(cmd.Slug),
		cmd.Index,
		strings.TrimSpace(cmd.Name),
		cmd.DescriptionShort,
		cmd.DescriptionFull,
		cmd.TRL,
		cmd.CustomMetric1, cmd.CustomMetric2, cmd.CustomMetric3, cmd.CustomMetric4,
		cmd.ImageURL,
		cmd.SourceLink,
		trendID,
	).Scan(&techID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%w: technology slug already exists", domain.ErrConflict)
		}
		return "", fmt.Errorf("insert technology: %w", err)
	}

	if err := replaceLinks(ctx, tx, techID, cmd); err != nil {
		return "", err
	}
	if err := replaceMetricValues(ctx, tx, techID, cmd.CustomMetrics); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return techID, nil
}

func (r *AdminTechnologyRepo) Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (string, bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", false, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var techID string
	err = tx.QueryRow(ctx, `SELECT id::text FROM technologies WHERE slug=$1 AND deleted_at IS NULL`, slug).Scan(&techID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("find technology: %w", err)
	}

	trendID, err := resolveTrendID(ctx, tx, cmd.TrendSlug)
	if err != nil {
		return "", false, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE technologies
		SET
			list_index=$2,
			name=$3,
			description_short=$4,
			description_full=$5,
			readiness_level=$6,
			custom_metric_1=$7,
			custom_metric_2=$8,
			custom_metric_3=$9,
			custom_metric_4=$10,
			image_url=$11,
			source_link=$12,
			trend_id=$13::uuid,
			updated_at=now()
		WHERE id=$1::uuid
	`,
		techID,
		cmd.Index,
		strings.TrimSpace(cmd.Name),
		cmd.DescriptionShort,
		cmd.DescriptionFull,
		cmd.TRL,
		cmd.CustomMetric1, cmd.CustomMetric2, cmd.CustomMetric3, cmd.CustomMetric4,
		cmd.ImageURL,
		cmd.SourceLink,
		trendID,
	)
	if err != nil {
		return "", false, fmt.Errorf("update technology: %w", err)
	}

	if err := replaceLinks(ctx, tx, techID, cmd); err != nil {
		return "", false, err
	}
	if err := replaceMetricValues(ctx, tx, techID, cmd.CustomMetrics); err != nil {
		return "", false, err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", false, fmt.Errorf("commit: %w", err)
	}
	return techID, true, nil
}

func (r *AdminTechnologyRepo) Delete(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		UPDATE technologies
		SET deleted_at = now(), updated_at = now()
		WHERE slug=$1 AND deleted_at IS NULL
	`, slug)
	if err != nil {
		return false, fmt.Errorf("delete technology: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminTechnologyRepo) Restore(ctx context.Context, slug string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		UPDATE technologies
		SET deleted_at = NULL, updated_at = now()
		WHERE slug=$1 AND deleted_at IS NOT NULL
	`, slug)
	if err != nil {
		return false, fmt.Errorf("restore technology: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminTechnologyRepo) List(ctx context.Context, p domain.AdminTechnologyListParams) ([]domain.TechnologyAdmin, int, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 50
	}
	offset := (p.Page - 1) * p.Limit
	whereClause := "WHERE tech.deleted_at IS NULL"
	if p.IncludeDeleted {
		whereClause = ""
	}

	rows, err := r.db.Query(ctx, fmt.Sprintf(`
		SELECT
			COUNT(*) OVER()::int,
			tech.id::text,
			tech.slug,
			tech.list_index,
			tech.name,
			tech.readiness_level,
			tech.description_short,
			tech.description_full,
			tech.custom_metric_1,
			tech.custom_metric_2,
			tech.custom_metric_3,
			tech.custom_metric_4,
			tech.image_url,
			tech.source_link,
			tech.deleted_at,
			tr.slug,
			tr.name
		FROM technologies tech
		JOIN trends tr ON tr.id = tech.trend_id
		%s
		ORDER BY tech.list_index ASC, tech.name ASC
		LIMIT $1 OFFSET $2
	`, whereClause), p.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list technologies: %w", err)
	}
	defer rows.Close()

	total := 0
	var out []domain.TechnologyAdmin
	for rows.Next() {
		var it domain.TechnologyAdmin
		var rowTotal int
		if err := rows.Scan(
			&rowTotal,
			&it.ID,
			&it.Slug,
			&it.Index,
			&it.Name,
			&it.TRL,
			&it.DescriptionShort,
			&it.DescriptionFull,
			&it.CustomMetric1,
			&it.CustomMetric2,
			&it.CustomMetric3,
			&it.CustomMetric4,
			&it.ImageURL,
			&it.SourceLink,
			&it.DeletedAt,
			&it.TrendSlug,
			&it.TrendName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan technology: %w", err)
		}
		total = rowTotal
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	if len(out) > 0 {
		ids := make([]string, 0, len(out))
		for _, it := range out {
			ids = append(ids, it.ID)
		}
		byTech, err := listDynamicMetricValuesByTechnologyIDs(ctx, r.db, ids)
		if err != nil {
			return nil, 0, err
		}
		for i := range out {
			out[i].CustomMetrics = byTech[out[i].ID]
		}
	}
	return out, total, nil
}

func (r *AdminTechnologyRepo) Get(ctx context.Context, slug string) (domain.TechnologyAdmin, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			tech.id::text,
			tech.slug,
			tech.list_index,
			tech.name,
			tech.readiness_level,
			tech.description_short,
			tech.description_full,
			tech.custom_metric_1,
			tech.custom_metric_2,
			tech.custom_metric_3,
			tech.custom_metric_4,
			tech.image_url,
			tech.source_link,
			tech.deleted_at,
			tr.slug,
			tr.name
		FROM technologies tech
		JOIN trends tr ON tr.id = tech.trend_id
		WHERE tech.slug = $1 AND tech.deleted_at IS NULL
	`, strings.TrimSpace(slug))

	var it domain.TechnologyAdmin
	if err := row.Scan(
		&it.ID,
		&it.Slug,
		&it.Index,
		&it.Name,
		&it.TRL,
		&it.DescriptionShort,
		&it.DescriptionFull,
		&it.CustomMetric1,
		&it.CustomMetric2,
		&it.CustomMetric3,
		&it.CustomMetric4,
		&it.ImageURL,
		&it.SourceLink,
		&it.DeletedAt,
		&it.TrendSlug,
		&it.TrendName,
	); err != nil {
		if err == pgx.ErrNoRows {
			return domain.TechnologyAdmin{}, false, nil
		}
		return domain.TechnologyAdmin{}, false, fmt.Errorf("get technology: %w", err)
	}

	var err error
	it.TagSlugs, err = listTextByTech(ctx, r.db, `
		SELECT t.slug
		FROM technology_tags tt
		JOIN tags t ON t.id = tt.tag_id
		WHERE tt.technology_id = $1::uuid
		ORDER BY t.slug ASC
	`, it.ID)
	if err != nil {
		return domain.TechnologyAdmin{}, false, err
	}

	it.SDGCodes, err = listTextByTech(ctx, r.db, `
		SELECT s.code
		FROM technology_sdgs ts
		JOIN sdgs s ON s.id = ts.sdg_id
		WHERE ts.technology_id = $1::uuid
		ORDER BY s.code ASC
	`, it.ID)
	if err != nil {
		return domain.TechnologyAdmin{}, false, err
	}

	it.OrganizationSlugs, err = listTextByTech(ctx, r.db, `
		SELECT o.slug
		FROM technology_organizations to2
		JOIN organizations o ON o.id = to2.organization_id
		WHERE to2.technology_id = $1::uuid
		ORDER BY o.slug ASC
	`, it.ID)
	if err != nil {
		return domain.TechnologyAdmin{}, false, err
	}
	it.CustomMetrics, err = listDynamicMetricValuesByTechnologyID(ctx, r.db, it.ID)
	if err != nil {
		return domain.TechnologyAdmin{}, false, err
	}

	return it, true, nil
}

func listDynamicMetricValuesByTechnologyID(ctx context.Context, db *pgxpool.Pool, techID string) ([]domain.TechnologyMetricValue, error) {
	byTech, err := listDynamicMetricValuesByTechnologyIDs(ctx, db, []string{techID})
	if err != nil {
		return nil, err
	}
	return byTech[techID], nil
}

func listDynamicMetricValuesByTechnologyIDs(ctx context.Context, db *pgxpool.Pool, techIDs []string) (map[string][]domain.TechnologyMetricValue, error) {
	out := make(map[string][]domain.TechnologyMetricValue, len(techIDs))
	if len(techIDs) == 0 {
		return out, nil
	}

	rows, err := db.Query(ctx, `
		SELECT
			tmv.technology_id::text,
			m.id::text,
			m.field_key,
			tmv.value
		FROM technology_metric_values tmv
		JOIN metrics_definitions m ON m.id = tmv.metric_id
		WHERE tmv.technology_id::text = ANY($1::text[])
		ORDER BY m.name ASC
	`, techIDs)
	if err != nil {
		return nil, fmt.Errorf("list admin tech dynamic metric values: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var techID, metricID string
		var fieldKey *string
		var value *float64
		if err := rows.Scan(&techID, &metricID, &fieldKey, &value); err != nil {
			return nil, fmt.Errorf("scan admin tech dynamic metric value: %w", err)
		}
		out[techID] = append(out[techID], domain.TechnologyMetricValue{
			MetricID: metricID,
			FieldKey: fieldKey,
			Value:    value,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows admin tech dynamic metric values: %w", err)
	}

	return out, nil
}

func listTextByTech(ctx context.Context, db *pgxpool.Pool, q string, techID string) ([]string, error) {
	rows, err := db.Query(ctx, q, techID)
	if err != nil {
		return nil, fmt.Errorf("list tech refs: %w", err)
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, fmt.Errorf("scan tech ref: %w", err)
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

// ---- helpers ----

func resolveTrendID(ctx context.Context, q pgx.Tx, trendSlug string) (string, error) {
	trendSlug = strings.TrimSpace(trendSlug)
	var id string
	err := q.QueryRow(ctx, `SELECT id::text FROM trends WHERE slug=$1 AND deleted_at IS NULL`, trendSlug).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err == pgx.ErrNoRows {
		return "", fmt.Errorf("%w: trend not found", domain.ErrInvalid)
	}
	return "", fmt.Errorf("resolve trend: %w", err)
}

func replaceLinks(ctx context.Context, tx pgx.Tx, techID string, cmd domain.TechnologyUpsert) error {
	// очистить старые связи
	if _, err := tx.Exec(ctx, `DELETE FROM technology_tags WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech tags: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM technology_sdgs WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech sdgs: %w", err)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM technology_organizations WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech orgs: %w", err)
	}

	// tags by slug
	if len(cmd.TagSlugs) > 0 {
		tagIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, slug FROM tags WHERE slug = ANY($1::text[]) AND deleted_at IS NULL`,
			cmd.TagSlugs,
		)
		if err != nil {
			return err
		}
		if err := insertLinkBatch(ctx, tx, `
			INSERT INTO technology_tags (technology_id, tag_id)
			SELECT $1::uuid, x::uuid
			FROM unnest($2::text[]) AS x
			ON CONFLICT DO NOTHING
		`, techID, tagIDs, "insert tech tag"); err != nil {
			return err
		}
	}

	// sdgs by code
	if len(cmd.SDGCodes) > 0 {
		sdgIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, code FROM sdgs WHERE code = ANY($1::text[]) AND deleted_at IS NULL`,
			cmd.SDGCodes,
		)
		if err != nil {
			return err
		}
		if err := insertLinkBatch(ctx, tx, `
			INSERT INTO technology_sdgs (technology_id, sdg_id)
			SELECT $1::uuid, x::uuid
			FROM unnest($2::text[]) AS x
			ON CONFLICT DO NOTHING
		`, techID, sdgIDs, "insert tech sdg"); err != nil {
			return err
		}
	}

	// orgs by slug
	if len(cmd.OrganizationSlugs) > 0 {
		orgIDs, err := resolveIDsByText(ctx, tx,
			`SELECT id::text, slug FROM organizations WHERE slug = ANY($1::text[]) AND deleted_at IS NULL`,
			cmd.OrganizationSlugs,
		)
		if err != nil {
			return err
		}
		if err := insertLinkBatch(ctx, tx, `
			INSERT INTO technology_organizations (technology_id, organization_id)
			SELECT $1::uuid, x::uuid
			FROM unnest($2::text[]) AS x
			ON CONFLICT DO NOTHING
		`, techID, orgIDs, "insert tech org"); err != nil {
			return err
		}
	}

	return nil
}

func replaceMetricValues(ctx context.Context, tx pgx.Tx, techID string, items []domain.TechnologyMetricValueUpsert) error {
	if _, err := tx.Exec(ctx, `DELETE FROM technology_metric_values WHERE technology_id=$1::uuid`, techID); err != nil {
		return fmt.Errorf("clear tech metric values: %w", err)
	}

	for _, it := range items {
		res, err := tx.Exec(ctx, `
			INSERT INTO technology_metric_values (technology_id, metric_id, value)
			SELECT $1::uuid, m.id, $3
			FROM metrics_definitions m
			WHERE m.id = $2::uuid AND m.deleted_at IS NULL
			ON CONFLICT (technology_id, metric_id)
			DO UPDATE SET value = EXCLUDED.value, updated_at = now()
		`, techID, strings.TrimSpace(it.MetricID), it.Value)
		if err != nil {
			return fmt.Errorf("upsert tech metric value: %w", err)
		}
		if res.RowsAffected() == 0 {
			return fmt.Errorf("%w: custom_metrics.metric_id is not found or deleted: %s", domain.ErrInvalid, strings.TrimSpace(it.MetricID))
		}
	}

	return nil
}

func insertLinkBatch(ctx context.Context, tx pgx.Tx, query, techID string, ids []string, op string) error {
	if len(ids) == 0 {
		return nil
	}
	if _, err := tx.Exec(ctx, query, techID, ids); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Возвращает ID-шники и валидирует, что все requested значения существуют.
func resolveIDsByText(ctx context.Context, tx pgx.Tx, sql string, wanted []string) ([]string, error) {
	wantedSet := map[string]struct{}{}
	for _, v := range wanted {
		v = strings.TrimSpace(v)
		if v != "" {
			wantedSet[v] = struct{}{}
		}
	}
	if len(wantedSet) == 0 {
		return nil, nil
	}

	uniq := make([]string, 0, len(wantedSet))
	for k := range wantedSet {
		uniq = append(uniq, k)
	}

	rows, err := tx.Query(ctx, sql, uniq)
	if err != nil {
		return nil, fmt.Errorf("resolve ids query: %w", err)
	}
	defer rows.Close()

	found := map[string]string{} // key -> id
	for rows.Next() {
		var id, key string
		if err := rows.Scan(&id, &key); err != nil {
			return nil, fmt.Errorf("resolve ids scan: %w", err)
		}
		found[key] = id
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("resolve ids rows: %w", err)
	}

	// validate missing
	for k := range wantedSet {
		if _, ok := found[k]; !ok {
			return nil, fmt.Errorf("%w: reference not found: %s", domain.ErrInvalid, k)
		}
	}

	out := make([]string, 0, len(found))
	for _, id := range found {
		out = append(out, id)
	}
	return out, nil
}
