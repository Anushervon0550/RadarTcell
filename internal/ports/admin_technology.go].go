package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminTechnologyService interface {
	List(ctx context.Context) ([]domain.TechnologyAdmin, error)
	Get(ctx context.Context, slug string) (domain.TechnologyAdmin, bool, error)
	Create(ctx context.Context, cmd domain.TechnologyUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTechnologyRepository interface {
	List(ctx context.Context) ([]domain.TechnologyAdmin, error)
	Get(ctx context.Context, slug string) (domain.TechnologyAdmin, bool, error)
	Create(ctx context.Context, cmd domain.TechnologyUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TechnologyUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}
