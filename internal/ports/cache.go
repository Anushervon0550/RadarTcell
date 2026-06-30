package ports

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

// RateLimiter — распределённый счётчик для rate-limiting (например, Redis).
type RateLimiter interface {
	// Incr увеличивает счётчик ключа в рамках фиксированного окна и
	// возвращает текущее значение счётчика и оставшийся TTL окна.
	Incr(ctx context.Context, key string, window time.Duration) (count int64, ttl time.Duration, err error)
}
