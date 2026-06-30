package httpapi

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// sensitiveQueryKeys — параметры, значения которых не должны попадать в логи.
var sensitiveQueryKeys = map[string]struct{}{
	"token": {}, "access_token": {}, "password": {}, "secret": {},
	"api_key": {}, "apikey": {}, "key": {}, "authorization": {},
}

// sanitizeQuery редактирует значения чувствительных query-параметров,
// сохраняя остальные для диагностики.
func sanitizeQuery(values url.Values) string {
	if len(values) == 0 {
		return ""
	}
	out := make(url.Values, len(values))
	for k, vs := range values {
		if _, sensitive := sensitiveQueryKeys[strings.ToLower(k)]; sensitive {
			out.Set(k, "[REDACTED]")
			continue
		}
		out[k] = vs
	}
	return out.Encode()
}

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
				zap.String("query", sanitizeQuery(r.URL.Query())),
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
