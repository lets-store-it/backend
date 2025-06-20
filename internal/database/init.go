package database

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/config"
	database "github.com/let-store-it/backend/generated/sqlc"
)

type Connection struct {
	Pool    *pgxpool.Pool
	Queries *database.Queries
}

func (c *Connection) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}

func InitDatabaseOrDie(ctx context.Context, cfg *config.Config) (*Connection, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %v", err)
	}

	// Add OpenTelemetry instrumentation
	poolConfig.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	// Record database stats with OpenTelemetry
	if err := otelpgx.RecordStats(pool); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to record database stats: %v", err)
	}

	return &Connection{
		Pool:    pool,
		Queries: database.New(pool),
	}, nil
}
