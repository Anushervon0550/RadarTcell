package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
)

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
				writeError(w, http.StatusInternalServerError, err.Error())
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
