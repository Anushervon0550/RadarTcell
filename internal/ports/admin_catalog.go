package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminTrendService interface {
	Create(ctx context.Context, cmd domain.TrendUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTrendRepository interface {
	Create(ctx context.Context, cmd domain.TrendUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTagService interface {
	Create(ctx context.Context, cmd domain.TagUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TagUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTagRepository interface {
	Create(ctx context.Context, cmd domain.TagUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TagUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}
