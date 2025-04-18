package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/handlers"
	"github.com/let-store-it/backend/internal/storeit/services"
	"github.com/let-store-it/backend/internal/storeit/services/yandex"
	"github.com/let-store-it/backend/internal/storeit/telemetry"
	"github.com/let-store-it/backend/internal/storeit/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server represents the main server instance and its dependencies
type Server struct {
	echo   *echo.Echo
	config *config.Config
}

// New creates and configures a new server instance
func New(cfg *config.Config, queries *database.Queries, pool *pgxpool.Pool) (*Server, error) {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	itemService := services.NewItemService(queries, pool)
	authService := services.NewAuthService(queries, pool)
	yandexOAuthService := yandex.NewYandexOAuthService(cfg.YandexOAuth.ClientID, cfg.YandexOAuth.ClientSecret)
	orgService := services.NewOrganizationService(queries, pool)
	storageGroupService := services.NewStorageService(queries)

	// Initialize use cases
	itemUseCase := usecases.NewItemUseCase(itemService)
	authUseCase := usecases.NewAuthUseCase(authService, yandexOAuthService)
	orgUseCase := usecases.NewOrganizationUseCase(orgService, authService)
	orgUnitUseCase := usecases.NewOrganizationUnitUseCase(orgService, authUseCase)
	storageGroupUseCase := usecases.NewStorageUseCase(storageGroupService, orgService, authService)

	// Initialize auth middleware
	authMiddleware := handlers.NewAuthMiddleware(authUseCase, "storeit_session", []string{"/auth", "/metrics"})
	e.Use(echo.WrapMiddleware(handlers.WithOrganizationID))
	e.Use(echo.WrapMiddleware(authMiddleware.Process))

	// Initialize API handlers
	handler := handlers.NewRestApiImplementation(
		orgUseCase,
		orgUnitUseCase,
		storageGroupUseCase,
		itemUseCase,
		authUseCase,
	)

	// Setup telemetry
	server, err := setupAPI(handler)
	if err != nil {
		return nil, err
	}

	// Setup routes
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.Any("/*", echo.WrapHandler(server))

	return &Server{
		echo:   e,
		config: cfg,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	return s.echo.Start(s.config.Server.ListenAddress)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func setupAPI(h api.Handler) (*api.Server, error) {
	meterProvider, meterShutdown, err := telemetry.NewMeterProvider()
	if err != nil {
		return nil, err
	}

	tracerProvider, tracerShutdown, err := telemetry.NewTracerProvider()
	if err != nil {
		meterShutdown(context.Background())
		return nil, err
	}

	server, err := api.NewServer(h,
		api.WithMeterProvider(meterProvider),
		api.WithTracerProvider(tracerProvider),
	)
	if err != nil {
		meterShutdown(context.Background())
		tracerShutdown(context.Background())
		return nil, err
	}

	return server, nil
}
