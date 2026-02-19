package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	DB          ports.DBPinger
	Catalog     ports.CatalogService
	Technology  ports.TechnologyService
	Preferences ports.PreferencesService
	Auth        ports.AuthService
}

func NewRouter(d RouterDeps) http.Handler {
	prefs := NewPreferencesHandler(d.Preferences)
	admin := NewAdminHandler(d.Auth)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	catalog := NewCatalogHandler(d.Catalog)
	tech := NewTechnologyHandler(d.Technology)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
		defer cancel()

		if err := d.DB.Ping(ctx); err != nil {
			writeError(w, http.StatusServiceUnavailable, "db not ready")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "ready"})
	})

	r.Route("/api", func(api chi.Router) {
		api.Get("/trends", catalog.ListTrends)
		api.Get("/sdgs", catalog.ListSDGs)
		api.Get("/tags", catalog.ListTags)
		api.Get("/organizations", catalog.ListOrganizations)
		api.Get("/metrics", catalog.ListMetrics)

		api.Get("/technologies", tech.List)

		api.Get("/technologies/{slug}", tech.Get)

		api.Get("/trends/{slug}/technologies", tech.ListByTrend)
		api.Get("/sdgs/{code}/technologies", tech.ListBySDG)
		api.Get("/tags/{slug}/technologies", tech.ListByTag)
		api.Get("/organizations/{slug}/technologies", tech.ListByOrganization)

		api.Get("/organizations/{slug}", catalog.GetOrganization)

		api.Post("/preferences", prefs.Save)
		api.Get("/preferences/{user_id}", prefs.Get)

	})

	r.Route("/api/admin", func(a chi.Router) {
		a.Post("/login", admin.Login)

		a.Group(func(pr chi.Router) {
			pr.Use(AuthRequired(d.Auth))
			pr.Get("/me", admin.Me)
		})
	})

	return r
}
