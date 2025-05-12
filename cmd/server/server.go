package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/server"
)

func main() {
	cfg := config.GetConfigOrDie()

	conn, err := database.InitDatabaseOrDie(context.Background(), cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	srv, err := server.New(cfg, conn.Queries, conn.Pool)
	if err != nil {
		slog.Error("Failed to create server", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
