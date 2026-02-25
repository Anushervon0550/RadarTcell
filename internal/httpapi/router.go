package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

	catalog := NewCatalogHandler(d.Catalog)
	tech := NewTechnologyHandler(d.Technology)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Health
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	// Readiness
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
		defer cancel()

		if err := d.DB.Ping(ctx); err != nil {
			writeError(w, http.StatusServiceUnavailable, "db not ready")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{"status": "ready"})
	})

	// OpenAPI file + Swagger UI
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/openapi.yaml")
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/openapi.yaml"),
	))

	// Public API
	r.Route("/api", func(api chi.Router) {
		// Catalog
		api.Get("/trends", catalog.ListTrends)
		api.Get("/sdgs", catalog.ListSDGs)
		api.Get("/tags", catalog.ListTags)
		api.Get("/organizations", catalog.ListOrganizations)
		api.Get("/metrics", catalog.ListMetrics)
		api.Get("/metrics/{id}/values", catalog.GetMetricValue)
		api.Get("/organizations/{slug}", catalog.GetOrganization)

		// Technologies
		api.Get("/technologies", tech.List)
		api.Get("/technologies/{slug}", tech.Get)

		// Relation endpoints
		api.Get("/trends/{slug}/technologies", tech.ListByTrend)
		api.Get("/sdgs/{code}/technologies", tech.ListBySDG)
		api.Get("/tags/{slug}/technologies", tech.ListByTag)
		api.Get("/organizations/{slug}/technologies", tech.ListByOrganization)

		// Preferences
		api.Post("/preferences", prefs.Save)
		api.Get("/preferences/{user_id}", prefs.Get)
	})

	// Admin API
	r.Route("/api/admin", func(a chi.Router) {
		a.Post("/login", admin.Login)

		a.Group(func(pr chi.Router) {
			pr.Use(AuthRequired(d.Auth))

			pr.Get("/me", admin.Me)

			// Technologies
			pr.Post("/technologies", adminTech.Create)
			pr.Put("/technologies/{slug}", adminTech.Update)
			pr.Delete("/technologies/{slug}", adminTech.Delete)

			// Trends
			pr.Post("/trends", adminCatalog.CreateTrend)
			pr.Put("/trends/{slug}", adminCatalog.UpdateTrend)
			pr.Delete("/trends/{slug}", adminCatalog.DeleteTrend)

			// Tags
			pr.Post("/tags", adminCatalog.CreateTag)
			pr.Put("/tags/{slug}", adminCatalog.UpdateTag)
			pr.Delete("/tags/{slug}", adminCatalog.DeleteTag)

			// Organizations
			pr.Post("/organizations", adminOrg.Create)
			pr.Put("/organizations/{slug}", adminOrg.Update)
			pr.Delete("/organizations/{slug}", adminOrg.Delete)

			// Metrics
			pr.Post("/metrics", adminMetrics.Create)
			pr.Put("/metrics/{id}", adminMetrics.Update)
			pr.Delete("/metrics/{id}", adminMetrics.Delete)
		})
	})

	return r
}
