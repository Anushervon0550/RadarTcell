package httpapi

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RateLimitConfig struct {
	Limit   int
	Window  time.Duration
	Message string
	KeyFunc func(r *http.Request) string
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

	keyFn := cfg.KeyFunc
	if keyFn == nil {
		keyFn = clientIPKey
	}

	var mu sync.Mutex
	state := map[string]rateLimitState{}
	nextCleanupAt := time.Now().Add(cfg.Window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := strings.TrimSpace(keyFn(r))
			if key == "" {
				key = "unknown"
			}

			now := time.Now()
			mu.Lock()
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
			remaining := s.resetAt.Sub(now)
			allowed := s.count <= cfg.Limit
			mu.Unlock()

			if !allowed {
				secs := int(remaining.Seconds())
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
