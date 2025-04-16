package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/config"
)

func InitDatabaseOrDie(ctx context.Context, config *config.Config) (*pgxpool.Pool, error) {
	pgxConfig, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.Database.User, config.Database.Password, config.Database.Name, config.Database.Host, config.Database.Port))
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
