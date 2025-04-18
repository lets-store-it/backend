package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/config"
	"github.com/let-store-it/backend/generated/database"
)

// Connection wraps database connection and queries
type Connection struct {
	Pool    *pgxpool.Pool
	Queries *database.Queries
}

// Close closes the database connection
func (c *Connection) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}

// InitDatabaseOrDie initializes database connection and queries
func InitDatabaseOrDie(ctx context.Context, cfg *config.Config) (*Connection, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return &Connection{
		Pool:    pool,
		Queries: database.New(pool),
	}, nil
}
