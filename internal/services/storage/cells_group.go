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
	"go.opentelemetry.io/otel/trace"
)

func toCellsGroupModel(group sqlc.CellsGroup) *models.CellsGroup {
	id := database.UuidFromPgx(group.ID)
	orgID := database.UuidFromPgx(group.OrgID)
	unitID := database.UuidFromPgx(group.UnitID)
	storageGroupID := database.UuidPtrFromPgx(group.StorageGroupID)

	return &models.CellsGroup{
		ID:             id,
		OrgID:          orgID,
		UnitID:         unitID,
		StorageGroupID: storageGroupID,
		Name:           group.Name,
		Alias:          group.Alias,
		CreatedAt:      group.CreatedAt.Time,
	}
}

func (s *StorageService) CreateCellsGroup(ctx context.Context, group *models.CellsGroup, name string, alias string) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "CreateCellsGroup",
		trace.WithAttributes(
			attribute.String("org_id", group.OrgID.String()),
			attribute.String("storage_group_id", group.ID.String()),
			attribute.String("unit_id", group.UnitID.String()),
			attribute.String("name", name),
			attribute.String("alias", alias),
		),
	)
	defer span.End()

	if err := s.validateName(name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}
	if err := s.validateAlias(alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "alias validation failed")
		return nil, err
	}

	cellsGroup, err := s.queries.CreateCellsGroup(ctx, sqlc.CreateCellsGroupParams{
		OrgID:          database.PgUUID(group.OrgID),
		StorageGroupID: database.PgUuidPtr(group.StorageGroupID),
		UnitID:         database.PgUUID(group.UnitID),
		Name:           name,
		Alias:          alias,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create cells group")
		return nil, fmt.Errorf("failed to create cells group: %w", err)
	}

	span.SetStatus(codes.Ok, "cells group created successfully")
	return toCellsGroupModel(cellsGroup), nil
}

func (s *StorageService) GetCellsGroups(ctx context.Context, orgID uuid.UUID) ([]*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetCellsGroups",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
		),
	)
	defer span.End()

	groups, err := s.queries.GetCellsGroups(ctx, database.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cells groups")
		return nil, fmt.Errorf("failed to get cells groups: %w", err)
	}

	result := make([]*models.CellsGroup, len(groups))
	for i, group := range groups {
		result[i] = toCellsGroupModel(group)
	}

	span.SetStatus(codes.Ok, "cells groups retrieved successfully")
	return result, nil
}

func (s *StorageService) GetCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetCellsGroup",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("group_id", id.String()),
		),
	)
	defer span.End()

	group, err := s.queries.GetCellsGroup(ctx, sqlc.GetCellsGroupParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cells group")
		return nil, fmt.Errorf("failed to get cells group: %w", err)
	}

	span.SetStatus(codes.Ok, "cells group retrieved successfully")
	return toCellsGroupModel(group), nil
}

func (s *StorageService) UpdateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateCellsGroup",
		trace.WithAttributes(
			attribute.String("org_id", group.OrgID.String()),
			attribute.String("group_id", group.ID.String()),
		),
	)
	defer span.End()

	if err := s.validateName(group.Name); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "name validation failed")
		return nil, err
	}
	if err := s.validateAlias(group.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "alias validation failed")
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
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update cells group")
		return nil, fmt.Errorf("failed to update cells group: %w", err)
	}

	span.SetStatus(codes.Ok, "cells group updated successfully")
	return toCellsGroupModel(updatedGroup), nil
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteCellsGroup",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("group_id", id.String()),
		),
	)
	defer span.End()

	err := s.queries.DeleteCellsGroup(ctx, sqlc.DeleteCellsGroupParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete cells group")
		return fmt.Errorf("failed to delete cells group: %w", err)
	}

	span.SetStatus(codes.Ok, "cells group deleted successfully")
	return nil
}
