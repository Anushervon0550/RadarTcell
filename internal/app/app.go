package app

import (
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/httpapi"
	"github.com/Anushervon0550/RadarTcell/internal/repository/postgres"
	"github.com/Anushervon0550/RadarTcell/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	AdminUser     string
	AdminPassword string
	JWTSecret     string
	JWTTTL        time.Duration
}

func BuildRouter(db *pgxpool.Pool, opt Options) (http.Handler, error) {
	// public repos
	catalogRepo := postgres.NewCatalogRepo(db)
	techRepo := postgres.NewTechnologyRepo(db)
	prefsRepo := postgres.NewPreferencesRepo(db)

	// public services
	catalogService := service.NewCatalogService(catalogRepo)
	techService := service.NewTechnologyService(techRepo)
	prefsService := service.NewPreferencesService(prefsRepo)

	// admin repos/services
	adminTechRepo := postgres.NewAdminTechnologyRepo(db)
	adminTechService := service.NewAdminTechnologyService(adminTechRepo)

	adminTrendRepo := postgres.NewAdminTrendRepo(db)
	adminTagRepo := postgres.NewAdminTagRepo(db)

	adminTrendService := service.NewAdminTrendService(adminTrendRepo)
	adminTagService := service.NewAdminTagService(adminTagRepo)

	adminOrgRepo := postgres.NewAdminOrganizationRepo(db)
	adminOrgService := service.NewAdminOrganizationService(adminOrgRepo)

	adminMetricRepo := postgres.NewAdminMetricRepo(db)
	adminMetricService := service.NewAdminMetricService(adminMetricRepo)

	// ✅ SDG admin repo/service (вот тут, ДО RouterDeps)
	adminSDGRepo := postgres.NewAdminSDGRepo(db)
	adminSDGService := service.NewAdminSDGService(adminSDGRepo)

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
		AdminSDG:          adminSDGService, // ✅ добавили
	}), nil
}
