package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/generated/database"
	db "github.com/let-store-it/backend/internal/storeit/database"
	"github.com/let-store-it/backend/internal/storeit/handlers"
	"github.com/let-store-it/backend/internal/storeit/services"
	"github.com/let-store-it/backend/internal/storeit/services/yandex"
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

func main() {
	config := config.GetConfigOrDie()

	dbCtx := context.Background()
	conn, err := db.InitDatabaseOrDie(dbCtx, config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	queries := database.New(conn)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize item layers
	itemService := services.NewItemService(queries, conn)
	itemUseCase := usecases.NewItemUseCase(itemService)

	// Initialize auth layers
	authService := services.NewAuthService(queries, conn)
	yandexOAuthService := yandex.NewYandexOAuthService(config.YandexOAuth.ClientID, config.YandexOAuth.ClientSecret)
	authUseCase := usecases.NewAuthUseCase(authService, yandexOAuthService)

	// Initialize organization layers
	orgService := services.NewOrganizationService(queries, conn)
	orgUseCase := usecases.NewOrganizationUseCase(orgService, authService)
	orgUnitUseCase := usecases.NewOrganizationUnitUseCase(orgService, authUseCase)

	// Initialize storage group layers
	storageGroupService := services.NewStorageService(queries)
	storageGroupUseCase := usecases.NewStorageUseCase(storageGroupService, orgService, authService)
	// Initialize handlers
	handler := handlers.NewRestApiImplementation(orgUseCase, orgUnitUseCase, storageGroupUseCase, itemUseCase, authUseCase)

	server, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add organization ID middleware
	orgIDMiddleware := handlers.NewOrganizationIDMiddleware(orgUseCase, authUseCase)
	e.Any("/*", echo.WrapHandler(orgIDMiddleware.WithOrganizationID(server)))

	go func() {
		if err := e.Start(config.Server.ListenAddress); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
