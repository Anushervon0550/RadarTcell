package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminI18nService interface {
	UpsertTrend(ctx context.Context, trendSlug string, cmd domain.TrendI18nUpsert) error
	GetTrend(ctx context.Context, trendSlug, locale string) (domain.TrendI18n, bool, error)
	DeleteTrend(ctx context.Context, trendSlug, locale string) (bool, error)

	UpsertTechnology(ctx context.Context, techSlug string, cmd domain.TechnologyI18nUpsert) error
	GetTechnology(ctx context.Context, techSlug, locale string) (domain.TechnologyI18n, bool, error)
	DeleteTechnology(ctx context.Context, techSlug, locale string) (bool, error)

	UpsertMetric(ctx context.Context, metricID string, cmd domain.MetricI18nUpsert) error
	GetMetric(ctx context.Context, metricID, locale string) (domain.MetricI18n, bool, error)
	DeleteMetric(ctx context.Context, metricID, locale string) (bool, error)
}

type AdminI18nRepository interface {
	UpsertTrend(ctx context.Context, trendSlug string, cmd domain.TrendI18nUpsert) error
	GetTrend(ctx context.Context, trendSlug, locale string) (domain.TrendI18n, bool, error)
	DeleteTrend(ctx context.Context, trendSlug, locale string) (bool, error)

	UpsertTechnology(ctx context.Context, techSlug string, cmd domain.TechnologyI18nUpsert) error
	GetTechnology(ctx context.Context, techSlug, locale string) (domain.TechnologyI18n, bool, error)
	DeleteTechnology(ctx context.Context, techSlug, locale string) (bool, error)

	UpsertMetric(ctx context.Context, metricID string, cmd domain.MetricI18nUpsert) error
	GetMetric(ctx context.Context, metricID, locale string) (domain.MetricI18n, bool, error)
	DeleteMetric(ctx context.Context, metricID, locale string) (bool, error)
}
