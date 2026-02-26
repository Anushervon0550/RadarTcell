package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type CatalogService interface {
	ListTrends(ctx context.Context) ([]domain.Trend, error)
	ListSDGs(ctx context.Context) ([]domain.SDG, error)
	ListTags(ctx context.Context) ([]domain.Tag, error)
	ListOrganizations(ctx context.Context) ([]domain.Organization, error)
	ListMetrics(ctx context.Context) ([]domain.MetricDefinition, error)
	GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error)
	GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error)
}

type TechnologyService interface {
	List(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Technology, bool, error)

	GetCard(ctx context.Context, slug string) (domain.TechnologyCard, bool, error)

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
	Create(ctx context.Context, cmd domain.SDGUpsert) (string, error)
	Update(ctx context.Context, code string, cmd domain.SDGUpsert) (bool, error)
	Delete(ctx context.Context, code string) (bool, error)
}
