package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

const (
	cacheVersionCatalog      = "v:catalog"
	cacheVersionTechnologies = "v:technologies"
)

func cacheVersion(ctx context.Context, cache ports.Cache, key string) string {
	if cache == nil {
		return "0"
	}
	b, ok, err := cache.Get(ctx, key)
	if err != nil || !ok {
		return "0"
	}
	return string(b)
}

func bumpCacheVersion(ctx context.Context, cache ports.Cache, key string) {
	if cache == nil {
		return
	}
	_ = cache.Set(ctx, key, []byte(strconv.FormatInt(time.Now().UnixNano(), 10)), 0)
}
