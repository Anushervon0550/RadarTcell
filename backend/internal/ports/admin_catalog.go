package ports

import (
	"context"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
)

type AdminTrendService interface {
	List(ctx context.Context) ([]domain.AdminTrend, error)
	Get(ctx context.Context, slug string) (domain.AdminTrend, bool, error)
	Create(ctx context.Context, cmd domain.TrendUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTrendRepository interface {
	List(ctx context.Context) ([]domain.AdminTrend, error)
	Get(ctx context.Context, slug string) (domain.AdminTrend, bool, error)
	Create(ctx context.Context, cmd domain.TrendUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TrendUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTagService interface {
	List(ctx context.Context) ([]domain.Tag, error)
	Get(ctx context.Context, slug string) (domain.Tag, bool, error)
	Create(ctx context.Context, cmd domain.TagUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TagUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}

type AdminTagRepository interface {
	List(ctx context.Context) ([]domain.Tag, error)
	Get(ctx context.Context, slug string) (domain.Tag, bool, error)
	Create(ctx context.Context, cmd domain.TagUpsert) (id string, err error)
	Update(ctx context.Context, slug string, cmd domain.TagUpsert) (id string, ok bool, err error)
	Delete(ctx context.Context, slug string) (ok bool, err error)
}
