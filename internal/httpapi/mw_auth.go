package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/domain"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

type ctxKey string

const ctxPrincipal ctxKey = "principal"

// AdminPrincipal возвращает аутентифицированного субъекта запроса.
func AdminPrincipal(r *http.Request) domain.Principal {
	v, _ := r.Context().Value(ctxPrincipal).(domain.Principal)
	return v
}

// AdminSubject возвращает имя субъекта (для логов/ответов).
func AdminSubject(r *http.Request) string {
	return AdminPrincipal(r).Subject
}

// AdminRole возвращает роль субъекта.
func AdminRole(r *http.Request) string {
	return AdminPrincipal(r).Role
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
			principal, ok, err := auth.Verify(r.Context(), token)
			if err != nil {
				writeInternalError(w)
				return
			}
			if !ok {
				writeError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), ctxPrincipal, principal)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole пропускает запрос только при совпадении роли субъекта.
// Должен применяться после AuthRequired.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if AdminRole(r) != role {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
