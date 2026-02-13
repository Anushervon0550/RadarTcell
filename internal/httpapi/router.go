package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Anushervon0550/RadarTcell/internal/repository"
)

type Deps struct {
	DB *pgxpool.Pool
}

func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	repo := repository.New(d.DB)

	// health endpoints
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 800*time.Millisecond)
		defer cancel()

		if err := d.DB.Ping(ctx); err != nil {
			http.Error(w, "db not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	// API v1
	r.Route("/api", func(api chi.Router) {
		api.Get("/trends", func(w http.ResponseWriter, r *http.Request) {
			items, err := repo.ListTrends(r.Context())
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, items)
		})

		api.Get("/sdgs", func(w http.ResponseWriter, r *http.Request) {
			items, err := repo.ListSDGs(r.Context())
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, items)
		})

		api.Get("/tags", func(w http.ResponseWriter, r *http.Request) {
			items, err := repo.ListTags(r.Context())
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, items)
		})

		api.Get("/organizations", func(w http.ResponseWriter, r *http.Request) {
			items, err := repo.ListOrganizations(r.Context())
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, items)
		})

		api.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
			items, err := repo.ListMetrics(r.Context())
			if err != nil {
				writeError(w, 500, err.Error())
				return
			}
			writeJSON(w, 200, items)
		})
	})

	return r
}
