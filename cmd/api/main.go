package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"radar-tcell/internal/db"
	"radar-tcell/internal/handlers"
)

func main() {
	// 1) Читаем порт из ENV
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// 2) Читаем строку подключения к БД
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is empty")
	}

	// 3) Подключаемся к Postgres (пул соединений)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPostgresPool(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	// 4) Роутер (стандартный)
	mux := http.NewServeMux()

	// 5) Healthcheck
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 6) Handlers API (у них есть доступ к БД)
	h := handlers.New(pool)

	// По ТЗ: GET /api/trends, GET /api/technologies
	mux.HandleFunc("/api/trends", h.GetTrends)
	mux.HandleFunc("/api/technologies", h.GetTechnologies)

	// 7) HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("listening on http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}
