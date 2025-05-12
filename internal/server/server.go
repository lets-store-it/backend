package server

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/handlers"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/employee"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/services/tasks"
	"github.com/let-store-it/backend/internal/services/tvboard"
	"github.com/let-store-it/backend/internal/services/yandex"
	"github.com/let-store-it/backend/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	auditUC "github.com/let-store-it/backend/internal/usecases/audit"
	authUC "github.com/let-store-it/backend/internal/usecases/auth"
	itemUC "github.com/let-store-it/backend/internal/usecases/item"
	organizationUC "github.com/let-store-it/backend/internal/usecases/organization"
	storageUC "github.com/let-store-it/backend/internal/usecases/storage"
	taskUC "github.com/let-store-it/backend/internal/usecases/task"
	tvboardUC "github.com/let-store-it/backend/internal/usecases/tv_board"
)

// Server represents the main server instance and its dependencies
type Server struct {
	echo         *echo.Echo
	config       *config.Config
	auditService *audit.AuditService
}

// New creates and configures a new server instance
func New(cfg *config.Config, queries *sqlc.Queries, pool *pgxpool.Pool) (*Server, error) {
	// Initialize telemetry
	if err := telemetry.InitTelemetry(context.Background(), cfg.ServiceName); err != nil {
		return nil, err
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Server.CorsOrigins,
		// AllowMethods: []string{
		// 	echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH, echo.OPTIONS,
		// },
		// AllowHeaders: []string{
		// 	echo.HeaderOrigin,
		// 	echo.HeaderContentType,
		// 	echo.HeaderAccept,
		// 	echo.HeaderAuthorization,
		// 	"X-Organization-Id",
		// 	"X-Api-Key",
		// 	"X-Requested-With",
		// 	"Access-Control-Allow-Origin",
		// },
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Initialize services
	employeeService := employee.New(employee.EmployeeServiceConfig{
		Queries: queries,
		PGXPool: pool,
	})
	storageGroupService, err := storage.New(&storage.StorageServiceConfig{
		Queries: queries,
		PGXPool: pool,
	})
	if err != nil {
		return nil, err
	}

	auditService, err := audit.New(audit.AuditServiceConfig{
		Queries:         queries,
		KafkaEnabled:    cfg.Kafka.Enabled,
		KafkaBrokers:    cfg.Kafka.GetBrokersList(),
		KafkaTopic:      cfg.Kafka.AuditTopic,
		PGXPool:         pool,
		EmployeeService: employeeService,
	})
	if err != nil {
		return nil, err
	}

	authService := auth.New(auth.AuthServiceConfig{
		Queries:      queries,
		PGXPool:      pool,
		AuditService: auditService,
	})

	itemService := item.New(item.ItemServiceConfig{
		Queries:        queries,
		PGXPool:        pool,
		StorageService: storageGroupService,
		AuditService:   auditService,
	})

	yandexOAuthService := yandex.NewYandexOAuthService(yandex.YandexOAuthServiceConfig{
		ClientID:     cfg.YandexOAuth.ClientID,
		ClientSecret: cfg.YandexOAuth.ClientSecret,
	})
	orgService := organization.New(organization.OrganizationServiceConfig{
		Queries:      queries,
		PGXPool:      pool,
		AuditService: auditService,
	})

	taskService := tasks.New(tasks.TaskServiceConfig{
		Queries:         queries,
		PGXPool:         pool,
		Auth:            authService,
		Org:             orgService,
		ItemService:     itemService,
		StorageService:  storageGroupService,
		EmployeeService: employeeService,
	})
	tvBoardService := tvboard.New(tvboard.TvBoardServiceConfig{
		Queries: queries,
		PGXPool: pool,
	})

	// Initialize use cases
	itemUseCase := itemUC.New(itemUC.ItemUseCaseConfig{
		Service:     itemService,
		AuthService: authService,
	})
	authUseCase := authUC.New(authUC.AuthUseCaseConfig{
		AuthService:        authService,
		YandexOAuthService: yandexOAuthService,
		EmployeeService:    employeeService,
	})
	orgUseCase := organizationUC.New(organizationUC.OrganizationUseCaseConfig{
		Service:     orgService,
		AuthService: authService,
	})

	storageUseCase := storageUC.New(storageUC.StorageUseCaseConfig{
		StorageService: storageGroupService,
		OrgService:     orgService,
		AuthService:    authService,
	})
	auditUseCase := auditUC.New(auditUC.AuditUseCaseConfig{
		AuthService:  authService,
		AuditService: auditService,
	})
	taskUseCase := taskUC.New(taskUC.TaskUseCaseConfig{
		TaskService: taskService,
		AuthService: authService,
		OrgService:  orgService,
	})
	tvBoardUseCase := tvboardUC.New(tvboardUC.TvBoardUseCaseConfig{
		TvBoardService:      tvBoardService,
		OrganizationService: orgService,
		AuthService:         authService,
	})

	// Initialize auth middleware
	e.Use(echo.WrapMiddleware(handlers.WithOrganizationID))
	e.Use(echo.WrapMiddleware(handlers.WithSetCookieFromContext))

	// Initialize API handlers
	handler := handlers.NewRestApiImplementation(
		orgUseCase,
		orgUseCase,
		storageUseCase,
		itemUseCase,
		authUseCase,
		auditUseCase,
		taskUseCase,
		tvBoardUseCase,
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
	e.POST("/auth/test", func(c echo.Context) error {
		return handler.TestAuth(c)
	})
	e.Any("/*", echo.WrapHandler(server))

	return &Server{
		echo:         e,
		config:       cfg,
		auditService: auditService,
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
	if err := s.auditService.Close(); err != nil {
		return err
	}
	return s.echo.Shutdown(ctx)
}
