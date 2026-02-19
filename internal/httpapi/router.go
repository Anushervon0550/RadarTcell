package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RouterDeps struct {
	DB      *pgxpool.Pool
	Catalog ports.CatalogService
}

func NewRouter(d RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	catalog := NewCatalogHandler(d.Catalog)

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
	})

	return r
}
