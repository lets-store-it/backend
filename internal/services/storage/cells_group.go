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
	"go.opentelemetry.io/otel/trace"
)

func (s *StorageService) CreateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateCellsGroup", func(ctx context.Context, span trace.Span) (*models.CellsGroup, error) {
		span.SetAttributes(
			attribute.String("org.id", group.OrgID.String()),
			attribute.String("storage_group.id", group.ID.String()),
			attribute.String("unit.id", group.UnitID.String()),
			attribute.String("cells_group.name", group.Name),
			attribute.String("cells_group.alias", group.Alias),
		)

		if err := s.validateName(group.Name); err != nil {
			return nil, err
		}
		if err := s.validateAlias(group.Alias); err != nil {
			return nil, err
		}

		cellsGroup, err := s.queries.CreateCellsGroup(ctx, sqlc.CreateCellsGroupParams{
			OrgID:          database.PgUUID(group.OrgID),
			StorageGroupID: database.PgUUIDPtr(group.StorageGroupID),
			UnitID:         database.PgUUID(group.UnitID),
			Name:           group.Name,
			Alias:          group.Alias,
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		model := toCellsGroupModel(cellsGroup)
		err = s.audit.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionCreate,
			TargetObjectType: models.ObjectTypeCellsGroup,
			TargetObjectID:   model.ID,
			PrechangeState:   nil,
			PostchangeState:  model,
		})
		if err != nil {
			return nil, err
		}

		return toCellsGroupModel(cellsGroup), nil
	})
}

func (s *StorageService) GetCellsGroups(ctx context.Context, orgID uuid.UUID) ([]*models.CellsGroup, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCellsGroups", func(ctx context.Context, span trace.Span) ([]*models.CellsGroup, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
		)

		groups, err := s.queries.GetCellsGroups(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := make([]*models.CellsGroup, len(groups))
		for i, group := range groups {
			result[i] = toCellsGroupModel(group)
		}

		return result, nil
	})
}

func (s *StorageService) GetCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.CellsGroup, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetCellsGroup", func(ctx context.Context, span trace.Span) (*models.CellsGroup, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cells_group.id", id.String()),
		)

		group, err := s.queries.GetCellsGroupById(ctx, sqlc.GetCellsGroupByIdParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toCellsGroupModel(group), nil
	})
}

func (s *StorageService) UpdateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateCellsGroup", func(ctx context.Context, span trace.Span) (*models.CellsGroup, error) {
		span.SetAttributes(
			attribute.String("org.id", group.OrgID.String()),
			attribute.String("cells_group.id", group.ID.String()),
		)

		if err := s.validateName(group.Name); err != nil {
			return nil, err
		}
		if err := s.validateAlias(group.Alias); err != nil {
			return nil, err
		}

		beforeUpdate, err := s.GetCellsGroup(ctx, group.OrgID, group.ID)
		if err != nil {
			return nil, err
		}

		updatedGroup, err := s.queries.UpdateCellsGroup(ctx, sqlc.UpdateCellsGroupParams{
			ID:     database.PgUUID(group.ID),
			OrgID:  database.PgUUID(group.OrgID),
			UnitID: database.PgUUID(group.UnitID),
			Name:   group.Name,
			Alias:  group.Alias,
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}
		model := toCellsGroupModel(updatedGroup)

		err = s.audit.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeCellsGroup,
			TargetObjectID:   group.ID,
			PrechangeState:   beforeUpdate,
			PostchangeState:  model,
		})
		if err != nil {
			return nil, err
		}

		return toCellsGroupModel(updatedGroup), nil
	})
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteCellsGroup", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("cells_group.id", id.String()),
		)

		err := s.queries.DeleteCellsGroup(ctx, sqlc.DeleteCellsGroupParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		return nil
	})
}
