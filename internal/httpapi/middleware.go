package httpapi

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// =============================== AUTH =======================================

type ctxKey string

const ctxAdminSubject ctxKey = "admin_subject"

func AdminSubject(r *http.Request) string {
	v, _ := r.Context().Value(ctxAdminSubject).(string)
	return v
}

func AuthRequired(auth ports.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			const pref = "Bearer "
			if !strings.HasPrefix(h, pref) {
				writeError(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(h, pref))
			sub, ok, err := auth.Verify(r.Context(), token)
			if err != nil {
				writeInternalError(w)
				return
			}
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), ctxAdminSubject, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// =============================== CORS / CSRF ================================

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	AllowedMethods   []string
	AllowCredentials bool
}

type CSRFConfig struct {
	TrustedOrigins []string
}

func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedOrigins := normalizeCSV(cfg.AllowedOrigins)
	allowedHeaders := cfg.AllowedHeaders
	if len(allowedHeaders) == 0 {
		allowedHeaders = []string{"Authorization", "Content-Type", "Accept"}
	}
	allowedMethods := cfg.AllowedMethods
	if len(allowedMethods) == 0 {
		allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && isOriginAllowed(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
				if cfg.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func CSRF(cfg CSRFConfig) func(http.Handler) http.Handler {
	trusted := normalizeCSV(cfg.TrustedOrigins)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isStateChanging(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")
			referer := r.Header.Get("Referer")

			if origin != "" {
				if !isOriginAllowed(origin, trusted) {
					writeError(w, http.StatusForbidden, "csrf: origin not allowed")
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if referer != "" {
				if !isOriginAllowed(extractOrigin(referer), trusted) {
					writeError(w, http.StatusForbidden, "csrf: referer not allowed")
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			writeError(w, http.StatusForbidden, "csrf: missing origin and referer")
		})
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	if len(allowed) == 0 {
		return false
	}
	if allowed[0] == "*" {
		return true
	}
	origin = strings.ToLower(strings.TrimSpace(origin))
	for _, a := range allowed {
		if origin == a {
			return true
		}
	}
	return false
}

func extractOrigin(referer string) string {
	if referer == "" {
		return ""
	}
	parts := strings.SplitN(referer, "/", 4)
	if len(parts) >= 3 {
		return strings.ToLower(strings.TrimSpace(parts[0] + "//" + parts[2]))
	}
	return ""
}

func normalizeCSV(in []string) []string {
	var out []string
	for _, raw := range in {
		parts := strings.Split(raw, ",")
		for _, p := range parts {
			p = strings.ToLower(strings.TrimSpace(p))
			if p != "" {
				out = append(out, p)
			}
		}
	}
	return out
}

func isStateChanging(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// =============================== IP ALLOW ===================================

var privateNetworks = mustParseCIDRs([]string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"::1/128",
	"fc00::/7",
})

func AllowPrivateNetworks() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIPKey(r)
			parsed := net.ParseIP(strings.TrimSpace(ip))
			if parsed == nil || !ipInAnyCIDR(parsed, privateNetworks) {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func mustParseCIDRs(values []string) []*net.IPNet {
	out := make([]*net.IPNet, 0, len(values))
	for _, v := range values {
		_, n, err := net.ParseCIDR(v)
		if err == nil {
			out = append(out, n)
		}
	}
	return out
}

func ipInAnyCIDR(ip net.IP, nets []*net.IPNet) bool {
	for _, n := range nets {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// =============================== LOGGING ====================================

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

func StructuredLogger(log *zap.Logger) func(http.Handler) http.Handler {
	if log == nil {
		log = zap.NewNop()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			requestID := middleware.GetReqID(r.Context())
			subject := AdminSubject(r)

			fields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", rec.status),
				zap.Int("bytes", rec.bytes),
				zap.Duration("duration", duration),
				zap.String("remote_ip", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			}
			if subject != "" {
				fields = append(fields, zap.String("subject", subject))
			}

			if rec.status >= 500 {
				log.Error("http_request", fields...)
				return
			}
			if rec.status >= 400 {
				log.Warn("http_request", fields...)
				return
			}
			log.Info("http_request", fields...)
		})
	}
}

// =============================== RATE LIMIT =================================

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
