package app

import (
	"net/http"

	"github.com/Anushervon0550/RadarTcell/internal/httpapi"
	"github.com/Anushervon0550/RadarTcell/internal/repository/postgres"
	"github.com/Anushervon0550/RadarTcell/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildRouter(db *pgxpool.Pool) http.Handler {
	// repositories
	catalogRepo := postgres.NewCatalogRepo(db)
	techRepo := postgres.NewTechnologyRepo(db)

	// services
	catalogService := service.NewCatalogService(catalogRepo)
	techService := service.NewTechnologyService(techRepo)

	// http router
	prefsRepo := postgres.NewPreferencesRepo(db)
	prefsService := service.NewPreferencesService(prefsRepo)

	return httpapi.NewRouter(httpapi.RouterDeps{
		DB:          db,
		Catalog:     catalogService,
		Technology:  techService,
		Preferences: prefsService,
	})

}
