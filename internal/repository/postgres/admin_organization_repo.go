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
		INSERT INTO organizations (slug, name, logo_url)
		VALUES ($1, $2, $3)
		RETURNING id::text
	`, strings.TrimSpace(cmd.Slug), strings.TrimSpace(cmd.Name), cmd.LogoURL).Scan(&id)

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
		SET name=$2, logo_url=$3, updated_at=now()
		WHERE id=$1::uuid
	`, id, strings.TrimSpace(cmd.Name), cmd.LogoURL)
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
