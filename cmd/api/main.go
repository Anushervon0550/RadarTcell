package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Anushervon0550/RadarTcell/internal/app"
	"github.com/Anushervon0550/RadarTcell/internal/cache"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	corsAllowedOrigins := splitEnvList("CORS_ALLOWED_ORIGINS")
	corsAllowedHeaders := splitEnvList("CORS_ALLOWED_HEADERS")
	corsAllowedMethods := splitEnvList("CORS_ALLOWED_METHODS")
	corsAllowCredentials := strings.EqualFold(strings.TrimSpace(os.Getenv("CORS_ALLOW_CREDENTIALS")), "true")
	csrfTrustedOrigins := splitEnvList("CSRF_TRUSTED_ORIGINS")

	redisAddr := strings.TrimSpace(os.Getenv("REDIS_ADDR"))
	redisPassword := strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))
	redisDB := envInt("REDIS_DB", 0)
	catalogCacheTTL := time.Duration(envInt("CATALOG_CACHE_TTL_SECONDS", 0)) * time.Second
	technologyCacheTTL := time.Duration(envInt("TECHNOLOGY_CACHE_TTL_SECONDS", 0)) * time.Second

	var cacheClient ports.Cache
	if redisAddr != "" {
		cacheClient = cache.NewRedisCache(redisAddr, redisPassword, redisDB)
	}

	// admin env (для JWT)
	adminUser := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASSWORD")
	jwtSecret := os.Getenv("JWT_SECRET")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// DB pool
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("db pool create error: %v", err)
	}
	defer pool.Close()

	// DB ping
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		log.Fatalf("db ping error: %v", err)
	}

	// Build router (composition root)
	router, err := app.BuildRouter(pool, app.Options{
		AdminUser:            adminUser,
		AdminPassword:        adminPass,
		JWTSecret:            jwtSecret,
		JWTTTL:               8 * time.Hour,
		CORSAllowedOrigins:   corsAllowedOrigins,
		CORSAllowedHeaders:   corsAllowedHeaders,
		CORSAllowedMethods:   corsAllowedMethods,
		CORSAllowCredentials: corsAllowCredentials,
		CSRFTrustedOrigins:   csrfTrustedOrigins,
		Cache:                cacheClient,
		CatalogCacheTTL:      catalogCacheTTL,
		TechnologyCacheTTL:   technologyCacheTTL,
	})
	if err != nil {
		log.Fatalf("app build error: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + appPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http server starting on :%s", appPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	log.Println("bye")
}

func splitEnvList(key string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return nil
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
