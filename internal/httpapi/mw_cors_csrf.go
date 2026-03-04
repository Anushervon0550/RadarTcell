package httpapi

import (
	"net/http"
	"strings"
)

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
			}

			next.ServeHTTP(w, r)
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
	// Expect "scheme://host/..."
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
