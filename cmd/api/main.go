package main

import (
	"context"
	"fmt"
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
	csrfTrustedOrigins = addSwaggerLocalOrigin(csrfTrustedOrigins, appPort, swaggerEnabled)

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
	adminAuthMode := strings.TrimSpace(os.Getenv("ADMIN_AUTH_MODE"))
	if adminAuthMode == "" {
		adminAuthMode = "db_then_env"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtTTLHours := envInt("JWT_TTL_HOURS", 8)
	if jwtTTLHours <= 0 {
		jwtTTLHours = 8
	}
	adminLoginRateLimit := envInt("ADMIN_LOGIN_RATE_LIMIT", 10)
	if adminLoginRateLimit <= 0 {
		adminLoginRateLimit = 10
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// storage (MinIO/S3) - после создания ctx
	minioEndpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	minioAccessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	minioSecretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	minioBucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	minioPublicURL := strings.TrimSpace(os.Getenv("MINIO_PUBLIC_URL"))
	minioUseSSL := strings.EqualFold(strings.TrimSpace(os.Getenv("MINIO_USE_SSL")), "true")
	minioPublicRead := envBool("MINIO_PUBLIC_READ", false)

	var storageClient ports.StorageService
	if minioEndpoint != "" && minioBucket != "" {
		st, err := storage.NewMinioStorage(minioEndpoint, minioAccessKey, minioSecretKey, minioBucket, minioPublicURL, minioUseSSL, minioPublicRead)
		if err != nil {
			logger.Fatal("minio storage init error", zap.Error(err))
		}
		if err := st.EnsureBucket(ctx); err != nil {
			logger.Warn("minio ensure bucket failed", zap.Error(err))
		}
		storageClient = st
	}

	// DB pool
	poolCfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Fatal("db config parse error", zap.Error(err))
	}
	databaseMaxConns := envInt("DATABASE_MAX_CONNS", 20)
	if databaseMaxConns <= 0 {
		databaseMaxConns = 20
	}
	poolCfg.MaxConns = int32(databaseMaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
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
		AdminAuthMode:        adminAuthMode,
		AdminLoginRateLimit:  adminLoginRateLimit,
		JWTSecret:            jwtSecret,
		JWTTTL:               time.Duration(jwtTTLHours) * time.Hour,
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

func envBool(key string, def bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
	}
}

func addSwaggerLocalOrigin(origins []string, appPort string, swaggerEnabled bool) []string {
	if !swaggerEnabled {
		return origins
	}
	appPort = strings.TrimSpace(appPort)
	if appPort == "" {
		appPort = "8080"
	}
	local := []string{
		fmt.Sprintf("http://localhost:%s", appPort),
		fmt.Sprintf("http://127.0.0.1:%s", appPort),
	}
	for _, v := range local {
		if !containsFoldTrim(origins, v) {
			origins = append(origins, v)
		}
	}
	return origins
}

func containsFoldTrim(items []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, it := range items {
		if strings.EqualFold(strings.TrimSpace(it), target) {
			return true
		}
	}
	return false
}
