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
	catalogRepo := postgres.NewCatalogRepo(db)
	techRepo := postgres.NewTechnologyRepo(db)
	prefsRepo := postgres.NewPreferencesRepo(db)

	catalogService := service.NewCatalogService(catalogRepo)
	techService := service.NewTechnologyService(techRepo)
	prefsService := service.NewPreferencesService(prefsRepo)

	adminTechRepo := postgres.NewAdminTechnologyRepo(db)
	adminTechService := service.NewAdminTechnologyService(adminTechRepo)

	authService, err := service.NewAuthService(opt.AdminUser, opt.AdminPassword, opt.JWTSecret, opt.JWTTTL)
	if err != nil {
		return nil, err
	}

	return httpapi.NewRouter(httpapi.RouterDeps{
		DB:              db,
		Catalog:         catalogService,
		Technology:      techService,
		Preferences:     prefsService,
		Auth:            authService,
		AdminTechnology: adminTechService,
	}), nil
}
