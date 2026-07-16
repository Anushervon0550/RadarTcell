package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Anushervon0550/RadarTcell/internal/app"
	"github.com/Anushervon0550/RadarTcell/internal/cache"
	"github.com/Anushervon0550/RadarTcell/internal/logging"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
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

	securityHSTS := strings.EqualFold(strings.TrimSpace(os.Getenv("SECURITY_HSTS")), "true")
	securityHSTSMaxAge := envInt("SECURITY_HSTS_MAX_AGE", 31536000)
	securityCSP := strings.TrimSpace(os.Getenv("SECURITY_CSP"))
	securityFrameAncestors := strings.TrimSpace(os.Getenv("SECURITY_FRAME_ANCESTORS"))
	serveFrontend := strings.EqualFold(strings.TrimSpace(os.Getenv("SERVE_FRONTEND")), "true")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
		AdminUser:              adminUser,
		AdminPassword:          adminPass,
		AdminAuthMode:          adminAuthMode,
		AdminLoginRateLimit:    adminLoginRateLimit,
		JWTSecret:              jwtSecret,
		JWTTTL:                 time.Duration(jwtTTLHours) * time.Hour,
		CORSAllowedOrigins:     corsAllowedOrigins,
		CORSAllowedHeaders:     corsAllowedHeaders,
		CORSAllowedMethods:     corsAllowedMethods,
		CORSAllowCredentials:   corsAllowCredentials,
		CSRFTrustedOrigins:     csrfTrustedOrigins,
		Cache:                  cacheClient,
		CatalogCacheTTL:        catalogCacheTTL,
		TechnologyCacheTTL:     technologyCacheTTL,
		Logger:                 logger,
		EnableSwagger:          swaggerEnabled,
		SecurityHSTS:           securityHSTS,
		SecurityHSTSMaxAge:     securityHSTSMaxAge,
		SecurityCSP:            securityCSP,
		SecurityFrameAncestors: securityFrameAncestors,
	})
	if err != nil {
		logger.Fatal("app build error", zap.Error(err))
	}
	if serveFrontend {
		router = withFrontend(router)
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
func withFrontend(h http.Handler) http.Handler {
	frontendDir := http.Dir("./web")
	fileServer := http.FileServer(frontendDir)
	// Оборачиваем статический файл-сервер, чтобы добавлять no-cache заголовки.
	// Без этого браузер агрессивно кеширует app.js/styles.css и не подхватывает изменения.
	static := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		fileServer.ServeHTTP(w, r)
	})
	mux := http.NewServeMux()
	mux.Handle("/api/", h)
	mux.Handle("/swagger/", h)
	mux.Handle("/openapi.yaml", h)
	mux.Handle("/healthz", h)
	mux.Handle("/readyz", h)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		cleanPath := path.Clean("/" + r.URL.Path)
		rel := strings.TrimPrefix(cleanPath, "/")
		if rel == "" || rel == "." {
			static.ServeHTTP(w, r)
			return
		}

		if _, err := os.Stat(filepath.Join("./web", rel)); err == nil {
			static.ServeHTTP(w, r)
			return
		}

		if path.Ext(rel) != "" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./web/index.html")
	})
	return mux
}
