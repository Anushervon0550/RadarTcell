package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type CatalogService interface {
	ListTrends(ctx context.Context, locale string) ([]domain.Trend, error)
	ListSDGs(ctx context.Context, locale string) ([]domain.SDG, error)
	ListTags(ctx context.Context, locale string) ([]domain.Tag, error)
	ListOrganizations(ctx context.Context, locale string) ([]domain.Organization, error)
	ListMetrics(ctx context.Context, locale string) ([]domain.MetricDefinition, error)
	GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error)
	GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error)
}

type TechnologyService interface {
	List(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, error)
	GetBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error)

	GetCard(ctx context.Context, slug, locale string) (domain.TechnologyCard, bool, error)

	ListByTrendSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error)
	ListBySDGCode(ctx context.Context, code string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error)
	ListByTagSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error)
	ListByOrganizationSlug(ctx context.Context, slug string, p domain.TechnologyListParams) (domain.TechnologyListResult, bool, error)
}

type PreferencesService interface {
	Save(ctx context.Context, p domain.Preferences) error
	Get(ctx context.Context, userID string) (domain.Preferences, bool, error)
}
type AdminSDGService interface {
	List(ctx context.Context) ([]domain.AdminSDG, error)
	Get(ctx context.Context, code string) (domain.AdminSDG, bool, error)
	Create(ctx context.Context, cmd domain.SDGUpsert) (string, error)
	Update(ctx context.Context, code string, cmd domain.SDGUpsert) (bool, error)
	Delete(ctx context.Context, code string) (bool, error)
}
