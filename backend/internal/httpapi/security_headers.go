package httpapi

import (
	"net/http"
	"strings"
)

// SecurityHeadersConfig configures the SecurityHeaders middleware.
// All fields are optional and use safe defaults when empty.
type SecurityHeadersConfig struct {
	// HSTS enables Strict-Transport-Security. Only send in production behind TLS.
	HSTS bool
	// HSTSMaxAge in seconds. Defaults to 31536000 (1 year).
	HSTSMaxAge int
	// ContentSecurityPolicy overrides the default CSP. Leave empty for a strict
	// API-only default appropriate for JSON responses and Swagger UI.
	ContentSecurityPolicy string
	// FrameAncestors controls Content-Security-Policy `frame-ancestors`.
	// Defaults to "'none'" to prevent clickjacking.
	FrameAncestors string
}

// SecurityHeaders sets industry-standard HTTP security headers on every
// response. Values follow OWASP secure headers recommendations.
func SecurityHeaders(cfg SecurityHeadersConfig) func(http.Handler) http.Handler {
	if cfg.HSTSMaxAge <= 0 {
		cfg.HSTSMaxAge = 31536000
	}
	if strings.TrimSpace(cfg.FrameAncestors) == "" {
		cfg.FrameAncestors = "'none'"
	}
	if strings.TrimSpace(cfg.ContentSecurityPolicy) == "" {
		cfg.ContentSecurityPolicy = "default-src 'none'; frame-ancestors " + cfg.FrameAncestors +
			"; base-uri 'none'; form-action 'self'"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")
			h.Set("Cross-Origin-Opener-Policy", "same-origin")
			h.Set("Cross-Origin-Resource-Policy", "same-site")
			if cfg.ContentSecurityPolicy != "" {
				h.Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
			}
			if cfg.HSTS {
				h.Set("Strict-Transport-Security", "max-age="+itoaInt(cfg.HSTSMaxAge)+"; includeSubDomains")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func itoaInt(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
