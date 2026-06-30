package cache

import (
	"context"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	c := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  3 * time.Second,
	})
	return &RedisCache{client: c}
}

// Ping проверяет доступность Redis (используется при старте для ранней диагностики).
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

var _ ports.Cache = (*RedisCache)(nil)

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return val, true, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Incr реализует ports.RateLimiter поверх Redis (фиксированное окно):
// INCR + ExpireNX (TTL ставится только при первом инкременте окна).
func (r *RedisCache) Incr(ctx context.Context, key string, window time.Duration) (int64, time.Duration, error) {
	pipe := r.client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, window)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, 0, err
	}
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil || ttl < 0 {
		ttl = window
	}
	return incr.Val(), ttl, nil
}
