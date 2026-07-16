package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// CatalogRepo aggregates public read-only catalog queries.
// Methods are split per-domain across separate files
// (trends.go, tags.go, sdgs.go, organizations.go, metrics.go).
type CatalogRepo struct {
	db *pgxpool.Pool
}

func NewCatalogRepo(db *pgxpool.Pool) *CatalogRepo {
	return &CatalogRepo{db: db}
}
