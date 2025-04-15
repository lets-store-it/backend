package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/evevseev/storeit/backend/config"
	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/generated/database"
	"github.com/evevseev/storeit/backend/internal/storeit/handlers"
	"github.com/evevseev/storeit/backend/internal/storeit/repositories"
	"github.com/evevseev/storeit/backend/internal/storeit/services"
	"github.com/evevseev/storeit/backend/internal/storeit/usecases"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfigOrDie()

	ctx := context.Background()
	conn, err := pgx.Connect(ctx,
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
			config.Database.User, config.Database.Password, config.Database.Name, config.Database.Host, config.Database.Port))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	e := echo.New()
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize database queries
	queries := database.New(conn)

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

	// Initialize handlers
	handler := handlers.NewRestApiImplementation(orgUseCase, orgUnitUseCase)

	server, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Add organization ID middleware
	e.Any("/*", echo.WrapHandler(handlers.WithOrganizationID(server)))

	go func() {
		if err := e.Start(config.Server.ListenAddress); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

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
