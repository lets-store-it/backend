package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/internal/database"
)

func main() {
	schemaPath := flag.String("schema", "", "Path to schema.sql file (required)")
	dryRun := flag.Bool("dry-run", false, "Verify schema without executing it")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *schemaPath == "" {
		slog.Error("Schema path is required")
		flag.Usage()
		os.Exit(1)
	}

	absSchemaPath, err := filepath.Abs(*schemaPath)
	if err != nil {
		slog.Error("Failed to resolve schema path", "error", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absSchemaPath); os.IsNotExist(err) {
		slog.Error("Schema file does not exist", "path", absSchemaPath)
		os.Exit(1)
	}

	schemaSQL, err := os.ReadFile(absSchemaPath)
	if err != nil {
		slog.Error("Failed to read schema file", "error", err, "path", absSchemaPath)
		os.Exit(1)
	}

	if *dryRun {
		slog.Info("Dry run - schema file read successfully", "path", absSchemaPath)
		return
	}

	cfg := config.GetConfigOrDie()

	conn, err := database.InitDatabaseOrDie(context.Background(), cfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Pool.Exec(context.Background(), string(schemaSQL))
	if err != nil {
		slog.Error("Failed to execute schema", "error", err)
		os.Exit(1)
	}

	slog.Info("Database initialized successfully", "schema_path", absSchemaPath)
}
