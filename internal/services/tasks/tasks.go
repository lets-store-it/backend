package tasks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

type TaskService struct {
	queries *database.Queries
	pgxpool *pgxpool.Pool
}

func New(queries *database.Queries, pgxpool *pgxpool.Pool) *TaskService {
	return &TaskService{queries: queries, pgxpool: pgxpool}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	description := pgtype.Text{}
	if task.Description != nil {
		description.Set(*task.Description)
	}

	tx, err := s.pgxpool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	createdTask, err := qtx.CreateTask(ctx, database.CreateTaskParams{
		OrgID:       pgtype.UUID{Bytes: task.OrgID, Valid: true},
		UnitID:      pgtype.UUID{Bytes: task.UnitID, Valid: true},
		Type:        string(task.Type),
		Name:        task.Name,
		Description: description,
	})

	return err
}
