package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
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
	AdminSDG          ports.AdminSDGService
	CORS              CORSConfig
	CSRF              CSRFConfig
	AdminI18n         ports.AdminI18nService
	Storage           ports.StorageService
	Logger            *zap.Logger
	EnableSwagger     bool
}

func NewRouter(d RouterDeps) http.Handler {
	prefs := NewPreferencesHandler(d.Preferences)
	admin := NewAdminHandler(d.Auth)
	adminTech := NewAdminTechnologyHandler(d.AdminTechnology)
	adminCatalog := NewAdminCatalogHandler(d.AdminTrend, d.AdminTag)
	adminOrg := NewAdminOrganizationHandler(d.AdminOrganization)
	adminMetrics := NewAdminMetricsHandler(d.AdminMetric)
	adminSDG := NewAdminSDGHandler(d.AdminSDG)
	adminI18n := NewAdminI18nHandler(d.AdminI18n)

	var upload *UploadHandler
	if d.Storage != nil {
		upload = NewUploadHandler(d.Storage)
	}

	catalog := NewCatalogHandler(d.Catalog)
	tech := NewTechnologyHandler(d.Technology)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(StructuredLogger(d.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(CORS(d.CORS))
	r.Use(CSRF(d.CSRF))

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

	// OpenAPI file + Swagger UI (optional)
	if d.EnableSwagger {
		r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "docs/openapi.yaml")
		})

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/openapi.yaml"),
		))
	}

	// Metrics
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

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

			// Upload (optional, если Storage настроен)
			if upload != nil {
				pr.Post("/upload", upload.Upload)
			}

			// Technologies
			pr.Get("/technologies", adminTech.List)
			pr.Get("/technologies/{slug}", adminTech.Get)
			pr.Post("/technologies", adminTech.Create)
			pr.Put("/technologies/{slug}", adminTech.Update)
			pr.Delete("/technologies/{slug}", adminTech.Delete)

			// Trends
			pr.Get("/trends", adminCatalog.ListTrends)
			pr.Get("/trends/{slug}", adminCatalog.GetTrend)
			pr.Post("/trends", adminCatalog.CreateTrend)
			pr.Put("/trends/{slug}", adminCatalog.UpdateTrend)
			pr.Delete("/trends/{slug}", adminCatalog.DeleteTrend)

			// Tags
			pr.Get("/tags", adminCatalog.ListTags)
			pr.Get("/tags/{slug}", adminCatalog.GetTag)
			pr.Post("/tags", adminCatalog.CreateTag)
			pr.Put("/tags/{slug}", adminCatalog.UpdateTag)
			pr.Delete("/tags/{slug}", adminCatalog.DeleteTag)

			// Organizations
			pr.Get("/organizations", adminOrg.List)
			pr.Get("/organizations/{slug}", adminOrg.Get)
			pr.Post("/organizations", adminOrg.Create)
			pr.Put("/organizations/{slug}", adminOrg.Update)
			pr.Delete("/organizations/{slug}", adminOrg.Delete)

			// Metrics
			pr.Get("/metrics", adminMetrics.List)
			pr.Get("/metrics/{id}", adminMetrics.Get)
			pr.Post("/metrics", adminMetrics.Create)
			pr.Put("/metrics/{id}", adminMetrics.Update)
			pr.Delete("/metrics/{id}", adminMetrics.Delete)

			pr.Get("/sdgs", adminSDG.List)
			pr.Get("/sdgs/{code}", adminSDG.Get)
			pr.Post("/sdgs", adminSDG.Create)
			pr.Put("/sdgs/{code}", adminSDG.Update)
			pr.Delete("/sdgs/{code}", adminSDG.Delete)

			// I18n
			pr.Put("/i18n/trends/{slug}", adminI18n.UpsertTrend)
			pr.Get("/i18n/trends/{slug}", adminI18n.GetTrend)
			pr.Delete("/i18n/trends/{slug}", adminI18n.DeleteTrend)

			pr.Put("/i18n/technologies/{slug}", adminI18n.UpsertTechnology)
			pr.Get("/i18n/technologies/{slug}", adminI18n.GetTechnology)
			pr.Delete("/i18n/technologies/{slug}", adminI18n.DeleteTechnology)

			pr.Put("/i18n/metrics/{id}", adminI18n.UpsertMetric)
			pr.Get("/i18n/metrics/{id}", adminI18n.GetMetric)
			pr.Delete("/i18n/metrics/{id}", adminI18n.DeleteMetric)
		})
	})

	return r
}
