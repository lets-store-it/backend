package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/handlers"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/services/yandex"
	"github.com/let-store-it/backend/internal/telemetry"
	"github.com/let-store-it/backend/internal/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server represents the main server instance and its dependencies
type Server struct {
	echo   *echo.Echo
	config *config.Config
}

// New creates and configures a new server instance
func New(cfg *config.Config, queries *database.Queries, pool *pgxpool.Pool) (*Server, error) {
	// Initialize telemetry
	if err := telemetry.InitTelemetry(context.Background()); err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost",
			"https://store-it.ru",
			"https://www.store-it.ru",
			"http://store-it.ru",
			"http://www.store-it.ru",
		},
		AllowMethods: []string{
			echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH, echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Organization-Id",
			"X-Api-Key",
			"X-Requested-With",
			"Access-Control-Allow-Origin",
		},
		AllowCredentials: true,
	}))
	// Initialize services
	storageGroupService, err := storage.New(&storage.StorageServiceConfig{
		Queries: queries,
	})
	if err != nil {
		return nil, err
	}
	authService := auth.New(queries, pool)

	auditService, err := audit.New(audit.AuditServiceConfig{
		Queries:      queries,
		KafkaEnabled: cfg.Kafka.Enabled,
		KafkaBrokers: cfg.Kafka.GetBrokersList(),
		KafkaTopic:   cfg.Kafka.AuditTopic,
		PGXPool:      pool,
		Auth:         authService,
	})
	if err != nil {
		return nil, err
	}
	defer auditService.Close()

	itemService := item.New(queries, pool, storageGroupService)
	yandexOAuthService := yandex.NewYandexOAuthService(cfg.YandexOAuth.ClientID, cfg.YandexOAuth.ClientSecret)
	orgService := organization.New(queries, pool)

	// Initialize use cases
	itemUseCase := usecases.NewItemUseCase(itemService)
	authUseCase := usecases.NewAuthUseCase(authService, yandexOAuthService)
	orgUseCase := usecases.NewOrganizationUseCase(orgService, authService, auditService)
	orgUnitUseCase := usecases.NewOrganizationUnitUseCase(orgService, authUseCase)
	storageGroupUseCase := usecases.NewStorageUseCase(storageGroupService, orgService, authService)
	auditUseCase := usecases.NewAuditUseCase(authService, auditService)

	// Initialize auth middleware
	e.Use(echo.WrapMiddleware(handlers.WithOrganizationID))

	// Initialize API handlers
	handler := handlers.NewRestApiImplementation(
		orgUseCase,
		orgUnitUseCase,
		storageGroupUseCase,
		itemUseCase,
		authUseCase,
		auditUseCase,
	)

	// Setup API server with global telemetry providers
	server, err := api.NewServer(handler, handler,
		api.WithMeterProvider(telemetry.GetMeterProvider()),
		api.WithTracerProvider(telemetry.GetTracerProvider()),
	)
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
	if err := telemetry.Shutdown(ctx); err != nil {
		return err
	}
	return s.echo.Shutdown(ctx)
}
