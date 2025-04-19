package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/let-store-it/backend/config"
	dbLayer "github.com/let-store-it/backend/internal/storeit/database"
	"github.com/let-store-it/backend/internal/storeit/server"
)

func main() {
	cfg := config.GetConfigOrDie()

	conn, err := dbLayer.InitDatabaseOrDie(context.Background(), cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create server instance
	srv, err := server.New(cfg, conn.Queries, conn.Pool)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("Server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
