package ports

import (
	"context"
	"encoding/json"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type CatalogRepository interface {
	ListTrends(ctx context.Context) ([]domain.Trend, error)
	ListSDGs(ctx context.Context) ([]domain.SDG, error)
	ListTags(ctx context.Context) ([]domain.Tag, error)
	ListOrganizations(ctx context.Context) ([]domain.Organization, error)
	ListMetrics(ctx context.Context) ([]domain.MetricDefinition, error)
}

type TechnologyRepository interface {
	ListTrendIDsOrdered(ctx context.Context) ([]string, error)
	ListTechnologies(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error)

	GetTechnologyBySlug(ctx context.Context, slug string) (*domain.Technology, bool, error)

	GetTrendIDBySlug(ctx context.Context, slug string) (string, bool, error)
	GetSDGIDByCode(ctx context.Context, code string) (string, bool, error)
	GetTagIDBySlug(ctx context.Context, slug string) (string, bool, error)
	GetOrganizationIDBySlug(ctx context.Context, slug string) (string, bool, error)

	ListTechnologyIDsByTrendID(ctx context.Context, trendID string) ([]string, error)
	ListTechnologyIDsBySDGID(ctx context.Context, sdgID string) ([]string, error)
	ListTechnologyIDsByTagID(ctx context.Context, tagID string) ([]string, error)
	ListTechnologyIDsByOrganizationID(ctx context.Context, orgID string) ([]string, error)
}

type PreferencesRepository interface {
	UpsertPreferences(ctx context.Context, userID string, settings json.RawMessage) error
	GetPreferences(ctx context.Context, userID string) (json.RawMessage, bool, error)
}
