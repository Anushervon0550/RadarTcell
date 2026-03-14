package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminI18nRepo struct {
	db *pgxpool.Pool
}

func NewAdminI18nRepo(db *pgxpool.Pool) *AdminI18nRepo {
	return &AdminI18nRepo{db: db}
}

var _ ports.AdminI18nRepository = (*AdminI18nRepo)(nil)

func (r *AdminI18nRepo) UpsertTrend(ctx context.Context, trendSlug string, cmd domain.TrendI18nUpsert) error {
	const q = `
		INSERT INTO trend_i18n (trend_id, locale, name, description)
		SELECT id, $2, $3, $4
		FROM trends
		WHERE slug = $1 AND deleted_at IS NULL
		ON CONFLICT (trend_id, locale)
		DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description;
	`
	ct, err := r.db.Exec(ctx, q, strings.TrimSpace(trendSlug), strings.TrimSpace(cmd.Locale), strings.TrimSpace(cmd.Name), cmd.Description)
	if err != nil {
		return fmt.Errorf("upsert trend i18n: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%w: trend not found", domain.ErrNotFound)
	}
	return nil
}

func (r *AdminI18nRepo) GetTrend(ctx context.Context, trendSlug, locale string) (domain.TrendI18n, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT ti.locale, ti.name, ti.description
		FROM trend_i18n ti
		JOIN trends t ON t.id = ti.trend_id
		WHERE t.slug = $1 AND ti.locale = $2
	`, strings.TrimSpace(trendSlug), strings.TrimSpace(locale))

	var it domain.TrendI18n
	if err := row.Scan(&it.Locale, &it.Name, &it.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TrendI18n{}, false, nil
		}
		return domain.TrendI18n{}, false, fmt.Errorf("get trend i18n: %w", err)
	}
	return it, true, nil
}

func (r *AdminI18nRepo) DeleteTrend(ctx context.Context, trendSlug, locale string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		DELETE FROM trend_i18n ti
		USING trends t
		WHERE t.id = ti.trend_id AND t.slug = $1 AND ti.locale = $2
	`, strings.TrimSpace(trendSlug), strings.TrimSpace(locale))
	if err != nil {
		return false, fmt.Errorf("delete trend i18n: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminI18nRepo) UpsertTechnology(ctx context.Context, techSlug string, cmd domain.TechnologyI18nUpsert) error {
	const q = `
		INSERT INTO technology_i18n (technology_id, locale, name, description_short, description_full)
		SELECT id, $2, $3, $4, $5
		FROM technologies
		WHERE slug = $1 AND deleted_at IS NULL
		ON CONFLICT (technology_id, locale)
		DO UPDATE SET name = EXCLUDED.name,
			description_short = EXCLUDED.description_short,
			description_full = EXCLUDED.description_full;
	`
	ct, err := r.db.Exec(ctx, q, strings.TrimSpace(techSlug), strings.TrimSpace(cmd.Locale), strings.TrimSpace(cmd.Name), cmd.DescriptionShort, cmd.DescriptionFull)
	if err != nil {
		return fmt.Errorf("upsert tech i18n: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%w: technology not found", domain.ErrNotFound)
	}
	return nil
}

func (r *AdminI18nRepo) GetTechnology(ctx context.Context, techSlug, locale string) (domain.TechnologyI18n, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT ti.locale, ti.name, ti.description_short, ti.description_full
		FROM technology_i18n ti
		JOIN technologies t ON t.id = ti.technology_id
		WHERE t.slug = $1 AND ti.locale = $2
	`, strings.TrimSpace(techSlug), strings.TrimSpace(locale))

	var it domain.TechnologyI18n
	if err := row.Scan(&it.Locale, &it.Name, &it.DescriptionShort, &it.DescriptionFull); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TechnologyI18n{}, false, nil
		}
		return domain.TechnologyI18n{}, false, fmt.Errorf("get tech i18n: %w", err)
	}
	return it, true, nil
}

func (r *AdminI18nRepo) DeleteTechnology(ctx context.Context, techSlug, locale string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		DELETE FROM technology_i18n ti
		USING technologies t
		WHERE t.id = ti.technology_id AND t.slug = $1 AND ti.locale = $2
	`, strings.TrimSpace(techSlug), strings.TrimSpace(locale))
	if err != nil {
		return false, fmt.Errorf("delete tech i18n: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}

func (r *AdminI18nRepo) UpsertMetric(ctx context.Context, metricID string, cmd domain.MetricI18nUpsert) error {
	const q = `
		INSERT INTO metric_definition_i18n (metric_id, locale, name, description)
		VALUES ($1::uuid, $2, $3, $4)
		ON CONFLICT (metric_id, locale)
		DO UPDATE SET name = EXCLUDED.name, description = EXCLUDED.description;
	`
	ct, err := r.db.Exec(ctx, q, strings.TrimSpace(metricID), strings.TrimSpace(cmd.Locale), strings.TrimSpace(cmd.Name), cmd.Description)
	if err != nil {
		return fmt.Errorf("upsert metric i18n: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%w: metric not found", domain.ErrNotFound)
	}
	return nil
}

func (r *AdminI18nRepo) GetMetric(ctx context.Context, metricID, locale string) (domain.MetricI18n, bool, error) {
	row := r.db.QueryRow(ctx, `
		SELECT locale, name, description
		FROM metric_definition_i18n
		WHERE metric_id = $1::uuid AND locale = $2
	`, strings.TrimSpace(metricID), strings.TrimSpace(locale))

	var it domain.MetricI18n
	if err := row.Scan(&it.Locale, &it.Name, &it.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.MetricI18n{}, false, nil
		}
		return domain.MetricI18n{}, false, fmt.Errorf("get metric i18n: %w", err)
	}
	return it, true, nil
}

func (r *AdminI18nRepo) DeleteMetric(ctx context.Context, metricID, locale string) (bool, error) {
	ct, err := r.db.Exec(ctx, `
		DELETE FROM metric_definition_i18n
		WHERE metric_id = $1::uuid AND locale = $2
	`, strings.TrimSpace(metricID), strings.TrimSpace(locale))
	if err != nil {
		return false, fmt.Errorf("delete metric i18n: %w", err)
	}
	return ct.RowsAffected() > 0, nil
}
