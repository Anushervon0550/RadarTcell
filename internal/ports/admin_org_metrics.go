package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminOrganizationService interface {
	Create(ctx context.Context, cmd domain.OrganizationUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminOrganizationRepository interface {
	Create(ctx context.Context, cmd domain.OrganizationUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.OrganizationUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminMetricService interface {
	Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (id string, err error)
	Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (ok bool, err error)
	Delete(ctx context.Context, id string) (ok bool, err error)
}

type AdminMetricRepository interface {
	Create(ctx context.Context, cmd domain.MetricDefinitionUpsert) (id string, err error)
	Update(ctx context.Context, id string, cmd domain.MetricDefinitionUpsert) (ok bool, err error)
	Delete(ctx context.Context, id string) (ok bool, err error)
}
