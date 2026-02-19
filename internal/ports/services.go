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
}

type TechnologyService interface {
	List(ctx context.Context, p domain.TechnologyListParams) (domain.TechnologyListResult, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Technology, bool, error)
}

type PreferencesService interface {
	Save(ctx context.Context, p domain.Preferences) error
	Get(ctx context.Context, userID string) (domain.Preferences, bool, error)
}
