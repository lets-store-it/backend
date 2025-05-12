package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *StorageService) GetCells(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCells", func(ctx context.Context, span trace.Span) ([]*models.Cell, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cells_group.id", cellsGroupID.String()),
		)

		cells, err := s.queries.GetCells(ctx, sqlc.GetCellsParams{
			OrgID:        database.PgUUID(orgID),
			CellsGroupID: database.PgUUID(cellsGroupID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := make([]*models.Cell, len(cells))
		for i, cell := range cells {
			result[i] = toCellModel(cell)
		}

		return result, nil
	})
}

func (s *StorageService) GetCellByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Cell, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCellByID", func(ctx context.Context, span trace.Span) (*models.Cell, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cell.id", id.String()),
		)

		cell, err := s.queries.GetCellById(ctx, sqlc.GetCellByIdParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := toCellModel(cell)
		return result, nil
	})
}

func (s *StorageService) CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateCell", func(ctx context.Context, span trace.Span) (*models.Cell, error) {
		if err := s.validateAlias(cell.Alias); err != nil {
			return nil, err
		}

		span.SetAttributes(
			attribute.String("org.id", cell.OrgID.String()),
			attribute.String("cells_group.id", cell.CellsGroupID.String()),
			attribute.String("cell.alias", cell.Alias),
			attribute.Int("cell.row", cell.Row),
			attribute.Int("cell.level", cell.Level),
			attribute.Int("cell.position", cell.Position),
		)

		createdCell, err := s.queries.CreateCell(ctx, sqlc.CreateCellParams{
			OrgID:        database.PgUUID(cell.OrgID),
			CellsGroupID: database.PgUUID(cell.CellsGroupID),
			Alias:        cell.Alias,
			Row:          int32(cell.Row),
			Level:        int32(cell.Level),
			Position:     int32(cell.Position),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := toCellModel(createdCell)
		return result, nil
	})
}

func (s *StorageService) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateCell", func(ctx context.Context, span trace.Span) (*models.Cell, error) {
		if cell == nil {
			return nil, services.ErrValidationError
		}

		if err := s.validateAlias(cell.Alias); err != nil {
			return nil, err
		}

		span.SetAttributes(
			attribute.String("cell.id", cell.ID.String()),
			attribute.String("org.id", cell.OrgID.String()),
			attribute.String("cells_group.id", cell.CellsGroupID.String()),
			attribute.String("cell.alias", cell.Alias),
			attribute.Int("cell.row", cell.Row),
			attribute.Int("cell.level", cell.Level),
			attribute.Int("cell.position", cell.Position),
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
			return nil, services.MapDbErrorToService(err)
		}

		result := toCellModel(updatedCell)
		return result, nil
	})
}

func (s *StorageService) DeleteCell(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteCell", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cell.id", id.String()),
		)

		err := s.queries.DeleteCell(ctx, sqlc.DeleteCellParams{
			ID:    database.PgUUID(id),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}

func (s *StorageService) GetCellFull(ctx context.Context, orgID uuid.UUID, cellID uuid.UUID) (*models.Cell, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCellFull", func(ctx context.Context, span trace.Span) (*models.Cell, error) {
		cellDb, err := s.queries.GetCellById(ctx, sqlc.GetCellByIdParams{
			ID:    database.PgUUID(cellID),
			OrgID: database.PgUUID(orgID),
		})

		if err != nil {
			if database.IsNotFound(err) {
				span.SetStatus(codes.Error, "cell not found")
				return nil, services.ErrNotFoundError
			}

			return nil, services.MapDbErrorToService(err)
		}

		cell := toCellModel(cellDb)

		cellPath, err := s.GetCellPath(ctx, orgID, cellID)
		if err != nil {
			return nil, err
		}

		cell.Path = &cellPath
		return cell, nil
	})
}

func (s *StorageService) GetCellPath(ctx context.Context, orgID uuid.UUID, cellID uuid.UUID) ([]models.CellPathSegment, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCellPath", func(ctx context.Context, span trace.Span) ([]models.CellPathSegment, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cell.id", cellID.String()),
		)

		segments, err := s.queries.GetCellPath(ctx, sqlc.GetCellPathParams{
			ID:    database.PgUUID(cellID),
			OrgID: database.PgUUID(orgID),
		})

		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toCellPathModel(segments), nil
	})
}
