package ports

import (
	"context"
	"encoding/json"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type CatalogRepository interface {
	ListTrends(ctx context.Context, locale string) ([]domain.Trend, error)
	ListSDGs(ctx context.Context, locale string) ([]domain.SDG, error)
	ListTags(ctx context.Context, locale string) ([]domain.Tag, error)
	ListOrganizations(ctx context.Context, locale string) ([]domain.Organization, error)
	ListMetrics(ctx context.Context, locale string) ([]domain.MetricDefinition, error)

	GetOrganizationBySlug(ctx context.Context, slug string) (domain.Organization, bool, error)

	GetMetricValue(ctx context.Context, metricID, technologyID string) (map[string]any, bool, error)
}

type TechnologyRepository interface {
	ListTrendIDsOrdered(ctx context.Context) ([]string, error)
	ListTechnologies(ctx context.Context, p domain.TechnologyListParams) ([]domain.Technology, int, error)
	GetMetricRanges(ctx context.Context) (map[string]domain.MetricRange, error)

	GetTechnologyBySlug(ctx context.Context, slug, locale string) (*domain.Technology, bool, error)
	GetTechnologyCardDataBySlug(ctx context.Context, slug, locale string) (domain.TechnologyCardData, bool, error)

	GetTrendIDBySlug(ctx context.Context, slug string) (string, bool, error)
	GetSDGIDByCode(ctx context.Context, code string) (string, bool, error)
	GetTagIDBySlug(ctx context.Context, slug string) (string, bool, error)
	GetOrganizationIDBySlug(ctx context.Context, slug string) (string, bool, error)

	ListTechnologyIDsByTrendID(ctx context.Context, trendID string) ([]string, error)
	ListTechnologyIDsBySDGID(ctx context.Context, sdgID string) ([]string, error)
	ListTechnologyIDsByTagID(ctx context.Context, tagID string) ([]string, error)
	ListTechnologyIDsByOrganizationID(ctx context.Context, orgID string) ([]string, error)

	ListTagsByTechnologyID(ctx context.Context, techID string) ([]domain.Tag, error)
	ListSDGsByTechnologyID(ctx context.Context, techID string) ([]domain.SDG, error)
	ListOrganizationsByTechnologyID(ctx context.Context, techID string) ([]domain.Organization, error)
	ListDynamicMetricValuesByTechnologyIDs(ctx context.Context, techIDs []string) (map[string][]domain.TechnologyMetricValue, error)
	ListDynamicMetricValuesByTechnologyID(ctx context.Context, techID string) ([]domain.TechnologyMetricValue, error)
}

type PreferencesRepository interface {
	UpsertPreferences(ctx context.Context, userID string, settings json.RawMessage) error
	GetPreferences(ctx context.Context, userID string) (json.RawMessage, bool, error)
}
type AdminSDGRepository interface {
	List(ctx context.Context) ([]domain.AdminSDG, error)
	Get(ctx context.Context, code string) (domain.AdminSDG, bool, error)
	Create(ctx context.Context, cmd domain.SDGUpsert) (string, error)
	Update(ctx context.Context, code string, cmd domain.SDGUpsert) (bool, error)
	Delete(ctx context.Context, code string) (bool, error)
}
