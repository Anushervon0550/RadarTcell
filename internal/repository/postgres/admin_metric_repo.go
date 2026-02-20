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
	var id string
	err := r.db.QueryRow(ctx, `
		INSERT INTO metrics_definitions (name, type, description, orderable)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text
	`, strings.TrimSpace(cmd.Name), strings.TrimSpace(cmd.Type), cmd.Description, cmd.Orderable).Scan(&id)

	if err != nil {
		return "", mapMetricPGErr(err, "metric already exists")
	}
	return id, nil
}

func (r *AdminMetricRepo) Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		UPDATE metrics_definitions
		SET name=$2, type=$3, description=$4, orderable=$5, updated_at=now()
		WHERE id=$1::uuid
	`, id, strings.TrimSpace(cmd.Name), strings.TrimSpace(cmd.Type), cmd.Description, cmd.Orderable)

	if err != nil {
		return false, mapMetricPGErr(err, "metric conflict")
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminMetricRepo) Delete(ctx context.Context, id string) (bool, error) {
	ct, err := r.db.Exec(ctx, `DELETE FROM metrics_definitions WHERE id=$1::uuid`, id)
	if err != nil {
		return false, mapMetricPGErr(err, "metric is referenced")
	}
	return ct.RowsAffected() > 0, nil
}

func mapMetricPGErr(err error, msg string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		case "23503":
			return fmt.Errorf("%w: %s", domain.ErrConflict, msg)
		}
	}
	if err == pgx.ErrNoRows {
		return fmt.Errorf("%w: not found", domain.ErrNotFound)
	}
	return fmt.Errorf("db error: %w", err)
}
