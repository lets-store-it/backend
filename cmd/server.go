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
	"github.com/let-store-it/backend/internal/storeit/repositories"
	"github.com/let-store-it/backend/internal/storeit/services"
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
	defer conn.Close(dbCtx)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize organization layers
	orgRepo := &repositories.OrganizationRepository{
		Queries: queries,
	}
	orgService := services.NewOrganizationService(orgRepo)
	orgUseCase := usecases.NewOrganizationUseCase(orgService)

	// Initialize organization unit layers
	orgUnitRepo := &repositories.OrganizationUnitRepository{
		Queries: queries,
	}
	orgUnitService := services.NewOrganizationUnitService(orgUnitRepo)
	orgUnitUseCase := usecases.NewOrganizationUnitUseCase(orgUnitService, orgService)

	// Initialize storage group layers
	storageGroupRepo := &repositories.StorageGroupRepository{
		Queries: queries,
	}
	storageGroupService := services.NewStorageGroupService(storageGroupRepo)
	storageGroupUseCase := usecases.NewStorageGroupUseCase(storageGroupService, orgService)

	// Initialize item layers
	itemRepo := repositories.NewItemRepository(queries, conn)
	itemService := services.NewItemService(itemRepo)
	itemUseCase := usecases.NewItemUseCase(itemService)

	// Initialize handlers
	handler := handlers.NewRestApiImplementation(orgUseCase, orgUnitUseCase, storageGroupUseCase, itemUseCase)

	server, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add organization ID middleware
	orgIDMiddleware := handlers.NewOrganizationIDMiddleware(orgUseCase)
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
