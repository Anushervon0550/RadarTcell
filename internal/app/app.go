package app

import (
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/httpapi"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/Anushervon0550/RadarTcell/internal/repository/postgres"
	"github.com/Anushervon0550/RadarTcell/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	AdminUser            string
	AdminPassword        string
	JWTSecret            string
	JWTTTL               time.Duration
	CORSAllowedOrigins   []string
	CORSAllowedHeaders   []string
	CORSAllowedMethods   []string
	CORSAllowCredentials bool
	CSRFTrustedOrigins   []string
	Cache                ports.Cache
	CatalogCacheTTL      time.Duration
	TechnologyCacheTTL   time.Duration
}

func BuildRouter(db *pgxpool.Pool, opt Options) (http.Handler, error) {
	// public repos
	catalogRepo := postgres.NewCatalogRepo(db)
	techRepo := postgres.NewTechnologyRepo(db)
	prefsRepo := postgres.NewPreferencesRepo(db)

	// public services
	catalogService := service.NewCatalogService(catalogRepo, opt.Cache, opt.CatalogCacheTTL)
	techService := service.NewTechnologyService(techRepo, opt.Cache, opt.TechnologyCacheTTL)
	prefsService := service.NewPreferencesService(prefsRepo)

	// admin repos/services
	adminTechRepo := postgres.NewAdminTechnologyRepo(db)
	adminTechService := service.NewAdminTechnologyService(adminTechRepo, opt.Cache)

	adminTrendRepo := postgres.NewAdminTrendRepo(db)
	adminTagRepo := postgres.NewAdminTagRepo(db)

	adminTrendService := service.NewAdminTrendService(adminTrendRepo, opt.Cache)
	adminTagService := service.NewAdminTagService(adminTagRepo, opt.Cache)

	adminOrgRepo := postgres.NewAdminOrganizationRepo(db)
	adminOrgService := service.NewAdminOrganizationService(adminOrgRepo, opt.Cache)

	adminMetricRepo := postgres.NewAdminMetricRepo(db)
	adminMetricService := service.NewAdminMetricService(adminMetricRepo, opt.Cache)

	// ✅ SDG admin repo/service (вот тут, ДО RouterDeps)
	adminSDGRepo := postgres.NewAdminSDGRepo(db)
	adminSDGService := service.NewAdminSDGService(adminSDGRepo, opt.Cache)

	// auth
	authService, err := service.NewAuthService(opt.AdminUser, opt.AdminPassword, opt.JWTSecret, opt.JWTTTL)
	if err != nil {
		return nil, err
	}

	// router deps
	return httpapi.NewRouter(httpapi.RouterDeps{
		DB:                db,
		Catalog:           catalogService,
		Technology:        techService,
		Preferences:       prefsService,
		Auth:              authService,
		AdminTechnology:   adminTechService,
		AdminTrend:        adminTrendService,
		AdminTag:          adminTagService,
		AdminOrganization: adminOrgService,
		AdminMetric:       adminMetricService,
		AdminSDG:          adminSDGService,
		CORS: httpapi.CORSConfig{
			AllowedOrigins:   opt.CORSAllowedOrigins,
			AllowedHeaders:   opt.CORSAllowedHeaders,
			AllowedMethods:   opt.CORSAllowedMethods,
			AllowCredentials: opt.CORSAllowCredentials,
		},
		CSRF: httpapi.CSRFConfig{
			TrustedOrigins: opt.CSRFTrustedOrigins,
		},
	}), nil
}
