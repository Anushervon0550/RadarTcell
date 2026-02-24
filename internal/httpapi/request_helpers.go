package httpapi

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func pathParamRequired(r *http.Request, name string) (string, bool) {
	v := strings.TrimSpace(chiURLParam(r, name))
	if v == "" {
		return "", false
	}
	return v, true
}

func queryParamRequired(r *http.Request, name string) (string, bool) {
	v := strings.TrimSpace(r.URL.Query().Get(name))
	if v == "" {
		return "", false
	}
	return v, true
}
func chiURLParam(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}
