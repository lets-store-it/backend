package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/let-store-it/backend/config"
	db "github.com/let-store-it/backend/internal/storeit/database"
	"github.com/let-store-it/backend/internal/storeit/server"
)

func main() {
	// Load configuration
	cfg := config.GetConfigOrDie()

	// Initialize database connection
	dbCtx := context.Background()
	conn, err := db.InitDatabaseOrDie(dbCtx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	// Create server instance
	srv, err := server.New(cfg, conn.Queries, conn.Pool)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Println("Server shutdown complete")
}
