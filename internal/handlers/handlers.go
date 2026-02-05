package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type Handlers struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Handlers {
	return &Handlers{pool: pool}
}
