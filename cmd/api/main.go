package main

import (
	"context"
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
	"github.com/Anushervon0550/RadarTcell/internal/logging"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/Anushervon0550/RadarTcell/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	logger, err := logging.NewLogger(os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	logger.Info("starting app")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Fatal("DATABASE_URL is required")
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
	swaggerEnabled := strings.EqualFold(strings.TrimSpace(os.Getenv("SWAGGER_ENABLED")), "true")

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

	// storage (MinIO/S3) - после создания ctx
	minioEndpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	minioAccessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	minioSecretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	minioBucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	minioPublicURL := strings.TrimSpace(os.Getenv("MINIO_PUBLIC_URL"))
	minioUseSSL := strings.EqualFold(strings.TrimSpace(os.Getenv("MINIO_USE_SSL")), "true")

	var storageClient ports.StorageService
	if minioEndpoint != "" && minioBucket != "" {
		st, err := storage.NewMinioStorage(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, minioPublicURL, minioUseSSL)
		if err != nil {
			logger.Fatal("minio storage init error", zap.Error(err))
		}
		if err := st.EnsureBucket(ctx); err != nil {
			logger.Warn("minio ensure bucket failed", zap.Error(err))
		}
		storageClient = st
	}

	// DB pool
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		logger.Fatal("db pool create error", zap.Error(err))
	}
	defer pool.Close()

	// DB ping
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		logger.Fatal("db ping error", zap.Error(err))
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
		Storage:              storageClient,
		Logger:               logger,
		EnableSwagger:        swaggerEnabled,
	})
	if err != nil {
		logger.Fatal("app build error", zap.Error(err))
	}

	srv := &http.Server{
		Addr:              ":" + appPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		logger.Info("http server starting", zap.String("addr", ":"+appPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http server error", zap.Error(err))
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Warn("graceful shutdown error", zap.Error(err))
	}
	logger.Info("bye")
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
