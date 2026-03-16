package app

import (
	"net/http"
	"time"

	"github.com/Anushervon0550/RadarTcell/internal/httpapi"
	"github.com/Anushervon0550/RadarTcell/internal/ports"
	"github.com/Anushervon0550/RadarTcell/internal/repository/postgres"
	"github.com/Anushervon0550/RadarTcell/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Options struct {
	AdminUser            string
	AdminPassword        string
	AdminAuthMode        string
	AdminLoginRateLimit  int
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
	Storage              ports.StorageService
	Logger               *zap.Logger
	EnableSwagger        bool
}

type publicServices struct {
	catalog     ports.CatalogService
	technology  ports.TechnologyService
	preferences ports.PreferencesService
}

type adminServices struct {
	technology   ports.AdminTechnologyService
	trend        ports.AdminTrendService
	tag          ports.AdminTagService
	organization ports.AdminOrganizationService
	metric       ports.AdminMetricService
	sdg          ports.AdminSDGService
	users        ports.AdminUserService
	i18n         ports.AdminI18nService
}

func BuildRouter(db *pgxpool.Pool, opt Options) (http.Handler, error) {
	pub := buildPublicServices(db, opt)
	adm := buildAdminServices(db, opt)
	authService, err := buildAuthService(db, opt)
	if err != nil {
		return nil, err
	}

	return httpapi.NewRouter(composeRouterDeps(db, opt, pub, adm, authService)), nil
}

func buildPublicServices(db *pgxpool.Pool, opt Options) publicServices {
	catalogRepo := postgres.NewCatalogRepo(db)
	techRepo := postgres.NewTechnologyRepo(db)
	prefsRepo := postgres.NewPreferencesRepo(db)

	return publicServices{
		catalog:     service.NewCatalogService(catalogRepo, opt.Cache, opt.CatalogCacheTTL),
		technology:  service.NewTechnologyService(techRepo, opt.Cache, opt.TechnologyCacheTTL),
		preferences: service.NewPreferencesService(prefsRepo),
	}
}

func buildAdminServices(db *pgxpool.Pool, opt Options) adminServices {
	adminTechRepo := postgres.NewAdminTechnologyRepo(db)
	adminTrendRepo := postgres.NewAdminTrendRepo(db)
	adminTagRepo := postgres.NewAdminTagRepo(db)
	adminOrgRepo := postgres.NewAdminOrganizationRepo(db)
	adminMetricRepo := postgres.NewAdminMetricRepo(db)
	adminI18nRepo := postgres.NewAdminI18nRepo(db)
	adminSDGRepo := postgres.NewAdminSDGRepo(db)
	adminUsersRepo := postgres.NewAdminUsersRepo(db)

	return adminServices{
		technology:   service.NewAdminTechnologyService(adminTechRepo, opt.Cache),
		trend:        service.NewAdminTrendService(adminTrendRepo, opt.Cache),
		tag:          service.NewAdminTagService(adminTagRepo, opt.Cache),
		organization: service.NewAdminOrganizationService(adminOrgRepo, opt.Cache),
		metric:       service.NewAdminMetricService(adminMetricRepo, opt.Cache),
		sdg:          service.NewAdminSDGService(adminSDGRepo, opt.Cache),
		users:        service.NewAdminUsersService(adminUsersRepo),
		i18n:         service.NewAdminI18nService(adminI18nRepo, opt.Cache),
	}
}

func buildAuthService(db *pgxpool.Pool, opt Options) (ports.AuthService, error) {
	authRepo := postgres.NewAuthRepo(db)
	return service.NewAuthService(authRepo, opt.AdminUser, opt.AdminPassword, opt.JWTSecret, opt.AdminAuthMode, opt.JWTTTL)
}

func composeRouterDeps(db *pgxpool.Pool, opt Options, pub publicServices, adm adminServices, auth ports.AuthService) httpapi.RouterDeps {
	return httpapi.RouterDeps{
		DB:                db,
		Catalog:           pub.catalog,
		Technology:        pub.technology,
		Preferences:       pub.preferences,
		Auth:              auth,
		AdminTechnology:   adm.technology,
		AdminTrend:        adm.trend,
		AdminTag:          adm.tag,
		AdminOrganization: adm.organization,
		AdminMetric:       adm.metric,
		AdminSDG:          adm.sdg,
		AdminUsers:        adm.users,
		AdminI18n:         adm.i18n,
		LoginRateLimit:    opt.AdminLoginRateLimit,
		Storage:           opt.Storage,
		Logger:            opt.Logger,
		EnableSwagger:     opt.EnableSwagger,
		CORS: httpapi.CORSConfig{
			AllowedOrigins:   opt.CORSAllowedOrigins,
			AllowedHeaders:   opt.CORSAllowedHeaders,
			AllowedMethods:   opt.CORSAllowedMethods,
			AllowCredentials: opt.CORSAllowCredentials,
		},
		CSRF: httpapi.CSRFConfig{TrustedOrigins: opt.CSRFTrustedOrigins},
	}
}
