package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *StorageService) GetCells(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	ctx, span := s.tracer.Start(ctx, "GetCells")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("cells_group_id", cellsGroupID.String()),
	)

	cells, err := s.queries.GetCells(ctx, sqlc.GetCellsParams{
		OrgID:        database.PgUUID(orgID),
		CellsGroupID: database.PgUUID(cellsGroupID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cells")
		return nil, fmt.Errorf("failed to get cells: %w", err)
	}

	result := make([]*models.Cell, len(cells))
	for i, cell := range cells {
		result[i], err = toCell(cell)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert cell")
			return nil, fmt.Errorf("failed to convert cell: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "cells retrieved successfully")
	return result, nil
}

func (s *StorageService) GetCellByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Cell, error) {
	ctx, span := s.tracer.Start(ctx, "GetCellByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("cell_id", id.String()),
	)

	cell, err := s.queries.GetCell(ctx, sqlc.GetCellParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cell")
		return nil, fmt.Errorf("failed to get cell: %w", err)
	}

	result, err := toCell(cell)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cell")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cell retrieved successfully")
	return result, nil
}

func (s *StorageService) CreateCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, alias string, row int, level int, position int) (*models.Cell, error) {
	ctx, span := s.tracer.Start(ctx, "CreateCell")
	defer span.End()

	if err := s.validateAlias(alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("cells_group_id", cellsGroupID.String()),
		attribute.String("alias", alias),
		attribute.Int("row", row),
		attribute.Int("level", level),
		attribute.Int("position", position),
	)

	cell, err := s.queries.CreateCell(ctx, sqlc.CreateCellParams{
		OrgID:        database.PgUUID(orgID),
		CellsGroupID: database.PgUUID(cellsGroupID),
		Alias:        alias,
		Row:          int32(row),
		Level:        int32(level),
		Position:     int32(position),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create cell")
		return nil, fmt.Errorf("failed to create cell: %w", err)
	}

	result, err := toCell(cell)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cell")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cell created successfully")
	return result, nil
}

func (s *StorageService) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateCell")
	defer span.End()

	if cell == nil {
		span.SetStatus(codes.Error, "invalid cell: nil")
		return nil, ErrInvalidCell
	}

	if err := s.validateAlias(cell.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	span.SetAttributes(
		attribute.String("cell_id", cell.ID.String()),
		attribute.String("org_id", cell.OrgID.String()),
		attribute.String("cells_group_id", cell.CellsGroupID.String()),
		attribute.String("alias", cell.Alias),
		attribute.Int("row", cell.Row),
		attribute.Int("level", cell.Level),
		attribute.Int("position", cell.Position),
	)

	updatedCell, err := s.queries.UpdateCell(ctx, sqlc.UpdateCellParams{
		ID:           database.PgUUID(cell.ID),
		OrgID:        database.PgUUID(cell.OrgID),
		CellsGroupID: database.PgUUID(cell.CellsGroupID),
		Alias:        cell.Alias,
		Row:          int32(cell.Row),
		Level:        int32(cell.Level),
		Position:     int32(cell.Position),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update cell")
		return nil, fmt.Errorf("failed to update cell: %w", err)
	}

	result, err := toCell(updatedCell)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cell")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cell updated successfully")
	return result, nil
}

func (s *StorageService) DeleteCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteCell")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("cells_group_id", cellsGroupID.String()),
		attribute.String("cell_id", id.String()),
	)

	err := s.queries.DeleteCell(ctx, sqlc.DeleteCellParams{
		ID:           database.PgUUID(id),
		OrgID:        database.PgUUID(orgID),
		CellsGroupID: database.PgUUID(cellsGroupID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete cell")
		return fmt.Errorf("failed to delete cell: %w", err)
	}

	span.SetStatus(codes.Ok, "cell deleted successfully")
	return nil
}

func (s *StorageService) GetCellPath(ctx context.Context, orgID uuid.UUID, cellID uuid.UUID) ([]models.CellPathSegment, error) {
	ctx, span := s.tracer.Start(ctx, "GetCellPath")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("cell_id", cellID.String()),
	)

	segments, err := s.queries.GetCellPath(ctx, sqlc.GetCellPathParams{
		ID:    database.PgUUID(cellID),
		OrgID: database.PgUUID(orgID),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cell path")
		return nil, fmt.Errorf("failed to get cell path: %w", err)
	}

	result := make([]models.CellPathSegment, len(segments))
	for i, segment := range segments {
		result[i] = models.CellPathSegment{
			ID:         database.UuidFromPgx(segment.ID),
			Name:       segment.Name,
			ObjectType: models.CellPathObjectType(segment.Type),
			Alias:      segment.Alias,
		}
	}

	span.SetStatus(codes.Ok, "cell path retrieved successfully")
	return result, nil
}
