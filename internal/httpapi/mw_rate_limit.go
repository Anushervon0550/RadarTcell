package httpapi

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type RateLimitConfig struct {
	// Name — уникальное пространство имён лимитера (для ключей в Redis).
	Name    string
	Limit   int
	Window  time.Duration
	Message string
	KeyFunc func(r *http.Request) string
	// Store — распределённый счётчик (Redis). Если nil — используется in-memory.
	Store ports.RateLimiter
}

type rateLimitState struct {
	count   int
	resetAt time.Time
}

func RateLimit(cfg RateLimitConfig) func(http.Handler) http.Handler {
	if cfg.Limit <= 0 {
		cfg.Limit = 60
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if strings.TrimSpace(cfg.Message) == "" {
		cfg.Message = "too many requests"
	}
	if strings.TrimSpace(cfg.Name) == "" {
		cfg.Name = "default"
	}

	keyFn := cfg.KeyFunc
	if keyFn == nil {
		keyFn = clientIPKey
	}

	var mu sync.Mutex
	state := map[string]rateLimitState{}
	nextCleanupAt := time.Now().Add(cfg.Window)

	// inMemory — фоллбэк, используется при отсутствии store или при ошибке Redis.
	inMemory := func(key string, now time.Time) (count int, retry time.Duration) {
		mu.Lock()
		defer mu.Unlock()
		if now.After(nextCleanupAt) {
			for k, v := range state {
				if now.After(v.resetAt.Add(cfg.Window)) {
					delete(state, k)
				}
			}
			nextCleanupAt = now.Add(cfg.Window)
		}
		s := state[key]
		if now.After(s.resetAt) || s.resetAt.IsZero() {
			s = rateLimitState{count: 0, resetAt: now.Add(cfg.Window)}
		}
		s.count++
		state[key] = s
		return s.count, s.resetAt.Sub(now)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := strings.TrimSpace(keyFn(r))
			if key == "" {
				key = "unknown"
			}

			now := time.Now()
			var count int
			var retry time.Duration

			if cfg.Store != nil {
				c, ttl, err := cfg.Store.Incr(r.Context(), "rl:"+cfg.Name+":"+key, cfg.Window)
				if err == nil {
					count = int(c)
					retry = ttl
				} else {
					// Redis недоступен — деградируем на локальный счётчик.
					count, retry = inMemory(key, now)
				}
			} else {
				count, retry = inMemory(key, now)
			}

			if count > cfg.Limit {
				secs := int(retry.Seconds())
				if secs < 1 {
					secs = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(secs))
				writeError(w, http.StatusTooManyRequests, cfg.Message)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIPKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}
