package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// =============================== PUBLIC =====================================

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
		WHERE m.deleted_at IS NULL
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

func (r *CatalogRepo) GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error) {
	const q = `
SELECT
	m.id::text,
	m.name,
	m.type,
	m.field_key,
	t.id::text,
	COALESCE(tmv.value,
	CASE m.field_key
		WHEN 'readiness_level' THEN t.readiness_level::double precision
		WHEN 'list_index' THEN t.list_index::double precision
		WHEN 'custom_metric_1' THEN t.custom_metric_1
		WHEN 'custom_metric_2' THEN t.custom_metric_2
		WHEN 'custom_metric_3' THEN t.custom_metric_3
		WHEN 'custom_metric_4' THEN t.custom_metric_4
		ELSE NULL
	END) AS value
FROM metrics_definitions m
JOIN technologies t ON t.id = $2 AND t.deleted_at IS NULL
LEFT JOIN technology_metric_values tmv ON tmv.metric_id = m.id AND tmv.technology_id = t.id
WHERE m.id = $1 AND m.deleted_at IS NULL
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
		return nil, false, fmt.Errorf("get metric value: %w", err)
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

// =============================== ADMIN ======================================

type AdminMetricRepo struct {
	db *pgxpool.Pool
}

func NewAdminMetricRepo(db *pgxpool.Pool) *AdminMetricRepo {
	return &AdminMetricRepo{db: db}
}

var _ ports.AdminMetricRepository = (*AdminMetricRepo)(nil)

func (r *AdminMetricRepo) Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (string, error) {
	const q = `
		INSERT INTO metrics_definitions (name, type, description, orderable, field_key)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	fieldKeyArg := nullableTrimmedString(cmd.FieldKey)

	var id string
	err := r.db.QueryRow(ctx, q,
		strings.TrimSpace(cmd.Name),
		strings.TrimSpace(cmd.Type),
		cmd.Description,
		cmd.Orderable,
		fieldKeyArg,
	).Scan(&id)
	if err != nil {
		return "", mapMetricPGErr(err, "metric already exists")
	}

	return id, nil
}

func (r *AdminMetricRepo) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	const q = `
		UPDATE metrics_definitions
		SET
			name = $2,
			type = $3,
			description = $4,
			orderable = $5,
			field_key = $6,
			updated_at = now()
		WHERE id = $1::uuid AND deleted_at IS NULL;
	`

	fieldKeyArg := nullableTrimmedString(cmd.FieldKey)

	ct, err := r.db.Exec(ctx, q,
		id,
		strings.TrimSpace(cmd.Name),
		strings.TrimSpace(cmd.Type),
		cmd.Description,
		cmd.Orderable,
		fieldKeyArg,
	)
	if err != nil {
		return false, mapMetricPGErr(err, "metric conflict")
	}

	return ct.RowsAffected() > 0, nil
}

func (r *AdminMetricRepo) Delete(ctx context.Context, id string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		UPDATE metrics_definitions
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1::uuid AND deleted_at IS NULL
	`, id)
	if err != nil {
		return false, mapMetricPGErr(err, "metric is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminMetricRepo) List(ctx context.Context) ([]domain.MetricDefinition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, name, type, description, orderable, field_key
		FROM metrics_definitions
		WHERE deleted_at IS NULL
		ORDER BY name ASC
	`)
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

func (r *AdminMetricRepo) Get(ctx context.Context, id string) (domain.MetricDefinition, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, name, type, description, orderable, field_key
		FROM metrics_definitions
		WHERE id = $1::uuid AND deleted_at IS NULL
	`, strings.TrimSpace(id))

	var it domain.MetricDefinition
	if err := row.Scan(&it.ID, &it.Name, &it.Type, &it.Description, &it.Orderable, &it.FieldKey); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MetricDefinition{}, false, nil
		}
		return domain.MetricDefinition{}, false, fmt.Errorf("get metric: %w", err)
	}
	return it, true, nil
}

// metric-specific PG error mapping: handles invalid UUID and "not found"
// in addition to common conflict cases.
func mapMetricPGErr(err error, msg string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "22P02": // invalid_text_representation (bad UUID)
			return fmt.Errorf("%w: invalid id", domain.ErrInvalid)
		}
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: not found", domain.ErrNotFound)
	}
	return fmt.Errorf("db error: %w", err)
}

func nullableTrimmedString(v *string) any {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return s
}
