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
	DB                ports.DBPinger
	Catalog           ports.CatalogService
	Technology        ports.TechnologyService
	Preferences       ports.PreferencesService
	Auth              ports.AuthService
	AdminTechnology   ports.AdminTechnologyService
	AdminTrend        ports.AdminTrendService
	AdminTag          ports.AdminTagService
	AdminOrganization ports.AdminOrganizationService
	AdminMetric       ports.AdminMetricService
}

func NewRouter(d RouterDeps) http.Handler {
	prefs := NewPreferencesHandler(d.Preferences)
	admin := NewAdminHandler(d.Auth)
	adminTech := NewAdminTechnologyHandler(d.AdminTechnology)
	adminCatalog := NewAdminCatalogHandler(d.AdminTrend, d.AdminTag)
	adminOrg := NewAdminOrganizationHandler(d.AdminOrganization)
	adminMetrics := NewAdminMetricsHandler(d.AdminMetric)

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

		r.Post("/api/preferences", prefs.Save)
		r.Get("/api/preferences/{user_id}", prefs.Get)

	})
	r.Get("/api/metrics/{id}/values", catalog.GetMetricValue)

	r.Route("/api/admin", func(a chi.Router) {
		a.Post("/login", admin.Login)

		a.Group(func(pr chi.Router) {
			pr.Use(AuthRequired(d.Auth))

			pr.Get("/me", admin.Me)

			pr.Post("/technologies", adminTech.Create)
			pr.Put("/technologies/{slug}", adminTech.Update)
			pr.Delete("/technologies/{slug}", adminTech.Delete)

			pr.Post("/trends", adminCatalog.CreateTrend)
			pr.Put("/trends/{slug}", adminCatalog.UpdateTrend)
			pr.Delete("/trends/{slug}", adminCatalog.DeleteTrend)

			pr.Post("/tags", adminCatalog.CreateTag)
			pr.Put("/tags/{slug}", adminCatalog.UpdateTag)
			pr.Delete("/tags/{slug}", adminCatalog.DeleteTag)

			pr.Post("/organizations", adminOrg.Create)
			pr.Put("/organizations/{slug}", adminOrg.Update)
			pr.Delete("/organizations/{slug}", adminOrg.Delete)

			pr.Post("/metrics", adminMetrics.Create)
			pr.Put("/metrics/{id}", adminMetrics.Update)
			pr.Delete("/metrics/{id}", adminMetrics.Delete)
		})

	})
	return r
}
