package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminOrganizationService interface {
	List(ctx context.Context) ([]domain.Organization, error)
	Get(ctx context.Context, slug string) (domain.Organization, bool, error)
	Create(ctx context.Context, cmd domain.OrganizationUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminOrganizationRepository interface {
	List(ctx context.Context) ([]domain.Organization, error)
	Get(ctx context.Context, slug string) (domain.Organization, bool, error)
	Create(ctx context.Context, cmd domain.OrganizationUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminMetricService interface {
	List(ctx context.Context) ([]domain.MetricDefinition, error)
	Get(ctx context.Context, id string) (domain.MetricDefinition, bool, error)
	Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (id string, err error)
	Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (ok bool, err error)
	Delete(ctx context.Context, id string) (ok bool, err error)
}

type AdminMetricRepository interface {
	List(ctx context.Context) ([]domain.MetricDefinition, error)
	Get(ctx context.Context, id string) (domain.MetricDefinition, bool, error)
	Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (id string, err error)
	Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (ok bool, err error)
	Delete(ctx context.Context, id string) (ok bool, err error)
}
