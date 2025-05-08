package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *StorageService) CreateCellsGroup(ctx context.Context, group models.StorageGroup) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "CreateCellsGroup")
	defer span.End()

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	span.SetAttributes(
		attribute.String("org_id", group.OrgID.String()),
		attribute.String("storage_group_id", group.ID.String()),
		attribute.String("name", group.Name),
		attribute.String("alias", group.Alias),
	)

	createdGroup, err := s.queries.CreateCellsGroup(ctx, database.CreateCellsGroupParams{
		OrgID:          utils.PgUUID(group.OrgID),
		StorageGroupID: utils.PgUUID(group.ID),
		UnitID:         utils.PgUUID(group.UnitID),
		Name:           group.Name,
		Alias:          group.Alias,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create cells group")
		return nil, fmt.Errorf("failed to create cells group: %w", err)
	}

	result, err := toCellsGroup(createdGroup)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cells group")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cells group created successfully")
	return result, nil
}

func (s *StorageService) GetAllCellsGroups(ctx context.Context, orgID uuid.UUID) ([]*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllCellsGroups")
	defer span.End()

	span.SetAttributes(attribute.String("org_id", orgID.String()))

	groups, err := s.queries.GetCellsGroups(ctx, utils.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cells groups")
		return nil, fmt.Errorf("failed to get cells groups: %w", err)
	}

	result := make([]*models.CellsGroup, len(groups))
	for i, group := range groups {
		result[i], err = toCellsGroup(group)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert cells group")
			return nil, fmt.Errorf("failed to convert cells group: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "cells groups retrieved successfully")
	return result, nil
}

func (s *StorageService) GetCellsGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetCellsGroupByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("group_id", id.String()),
	)

	group, err := s.queries.GetCellsGroup(ctx, database.GetCellsGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get cells group")
		return nil, fmt.Errorf("failed to get cells group: %w", err)
	}

	result, err := toCellsGroup(group)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cells group")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cells group retrieved successfully")
	return result, nil
}

func (s *StorageService) UpdateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateCellsGroup")
	defer span.End()

	if group == nil {
		span.SetStatus(codes.Error, "invalid cells group: nil")
		return nil, ErrInvalidCellsGroup
	}

	span.SetAttributes(
		attribute.String("group_id", group.ID.String()),
		attribute.String("org_id", group.OrgID.String()),
		attribute.String("name", group.Name),
		attribute.String("alias", group.Alias),
	)

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateCellsGroup(ctx, database.UpdateCellsGroupParams{
		ID:    utils.PgUUID(group.ID),
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update cells group")
		return nil, fmt.Errorf("failed to update cells group: %w", err)
	}

	result, err := toCellsGroup(updatedGroup)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert cells group")
		return nil, err
	}

	span.SetStatus(codes.Ok, "cells group updated successfully")
	return result, nil
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteCellsGroup")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("group_id", id.String()),
	)

	err := s.queries.DeleteCellsGroup(ctx, database.DeleteCellsGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete cells group")
		return fmt.Errorf("failed to delete cells group: %w", err)
	}

	span.SetStatus(codes.Ok, "cells group deleted successfully")
	return nil
}
