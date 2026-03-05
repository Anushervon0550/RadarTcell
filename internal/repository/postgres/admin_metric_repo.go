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
		WHERE id = $1::uuid;
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
	ct, err := r.db.Exec(ctx, `DELETE FROM metrics_definitions WHERE id = $1::uuid`, id)
	if err != nil {
		return false, mapMetricPGErr(err, "metric is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminMetricRepo) List(ctx context.Context) ([]domain.MetricDefinition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, name, type, description, orderable, field_key
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
		WHERE id = $1::uuid
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

func mapMetricPGErr(err error, msg string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "23503": // foreign_key_violation
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "22P02": // invalid_text_representation (например, плохой UUID)
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
