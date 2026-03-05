package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/app"
	"github.com/Anushervon0550/RadarTcell/internal/cache"
	"github.com/Anushervon0550/RadarTcell/internal/logging"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BenchmarkGetTechnologies(b *testing.B) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		b.Skip("DATABASE_URL not set, skipping benchmark")
	}

	logger, _ := logging.NewLogger("test")
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		b.Fatalf("db pool: %v", err)
	}
	defer pool.Close()

	var cacheClient ports.Cache
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		cacheClient = cache.NewRedisCache(redisAddr, "", 0)
	}

	router, err := app.BuildRouter(pool, app.Options{
		AdminUser:     "admin",
		AdminPassword: "admin123",
		JWTSecret:     "test-secret",
		JWTTTL:        8 * time.Hour,
		Cache:         cacheClient,
		Logger:        logger,
		EnableSwagger: false,
	})
	if err != nil {
		b.Fatalf("build router: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/technologies?limit=100", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			b.Fatalf("expected 200, got %d", rec.Code)
		}
	}
}

func TestGetTechnologiesLatency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping latency test in short mode")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}

	logger, _ := logging.NewLogger("test")
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("db pool: %v", err)
	}
	defer pool.Close()

	router, err := app.BuildRouter(pool, app.Options{
		AdminUser:     "admin",
		AdminPassword: "admin123",
		JWTSecret:     "test-secret",
		JWTTTL:        8 * time.Hour,
		Logger:        logger,
	})
	if err != nil {
		t.Fatalf("build router: %v", err)
	}

	// Прогрев
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/technologies?limit=100", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	// Измеряем latency
	const iterations = 100
	var total time.Duration

	for i := 0; i < iterations; i++ {
		start := time.Now()
		req := httptest.NewRequest(http.MethodGet, "/api/technologies?limit=100", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		elapsed := time.Since(start)
		total += elapsed

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
	}

	avg := total / iterations
	fmt.Printf("Average latency: %v\n", avg)

	const threshold = 200 * time.Millisecond
	if avg > threshold {
		t.Errorf("Average latency %v exceeds threshold %v", avg, threshold)
	}
}
