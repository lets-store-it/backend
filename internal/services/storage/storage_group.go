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

func (s *StorageService) validateStorageGroupData(name, alias string) error {
	if err := s.validateName(name); err != nil {
		return err
	}
	return s.validateAlias(alias)
}

func (s *StorageService) CreateStorageGroup(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "CreateStorageGroup")
	defer span.End()

	if err := s.validateStorageGroupData(name, alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("unit_id", unitID.String()),
		attribute.String("name", name),
		attribute.String("alias", alias),
	)
	if parentID != nil {
		span.SetAttributes(attribute.String("parent_id", parentID.String()))
	}

	group, err := s.queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    utils.PgUUID(orgID),
		UnitID:   utils.PgUUID(unitID),
		ParentID: utils.NullablePgUUID(parentID),
		Name:     name,
		Alias:    alias,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create storage group")
		return nil, fmt.Errorf("failed to create storage group: %w", err)
	}

	result, err := toStorageGroup(group)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert storage group")
		return nil, err
	}

	span.SetStatus(codes.Ok, "storage group created successfully")
	return result, nil
}

func (s *StorageService) GetAllStorageGroups(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllStorageGroups")
	defer span.End()

	span.SetAttributes(attribute.String("org_id", orgID.String()))

	groups, err := s.queries.GetStorageGroups(ctx, utils.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get storage groups")
		return nil, fmt.Errorf("failed to get storage groups: %w", err)
	}

	result := make([]*models.StorageGroup, len(groups))
	for i, group := range groups {
		result[i], err = toStorageGroup(group)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert storage group")
			return nil, fmt.Errorf("failed to convert storage group: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "storage groups retrieved successfully")
	return result, nil
}

func (s *StorageService) GetStorageGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetStorageGroupByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("group_id", id.String()),
	)

	group, err := s.queries.GetStorageGroup(ctx, database.GetStorageGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get storage group")
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	model, err := toStorageGroup(group)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert storage group")
		return nil, err
	}
	if model == nil {
		span.SetStatus(codes.Error, "storage group not found")
		return nil, ErrStorageGroupNotFound
	}

	span.SetStatus(codes.Ok, "storage group retrieved successfully")
	return model, nil
}

func (s *StorageService) DeleteStorageGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteStorageGroup")
	defer span.End()

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("group_id", id.String()),
	)

	err := s.queries.DeleteStorageGroup(ctx, database.DeleteStorageGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete storage group")
		return fmt.Errorf("failed to delete storage group: %w", err)
	}

	span.SetStatus(codes.Ok, "storage group deleted successfully")
	return nil
}

func (s *StorageService) UpdateStoragrGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateStorageGroup")
	defer span.End()

	if group == nil {
		span.SetStatus(codes.Error, "invalid storage group: nil")
		return nil, ErrInvalidStorageGroup
	}

	span.SetAttributes(
		attribute.String("group_id", group.ID.String()),
		attribute.String("name", group.Name),
		attribute.String("alias", group.Alias),
	)

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateStorageGroup(ctx, database.UpdateStorageGroupParams{
		ID:     utils.PgUUID(group.ID),
		OrgID:  utils.PgUUID(group.OrgID),
		UnitID: utils.PgUUID(group.UnitID),
		Name:   group.Name,
		Alias:  group.Alias,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update storage group")
		return nil, fmt.Errorf("failed to update storage group: %w", err)
	}

	result, err := toStorageGroup(updatedGroup)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert storage group")
		return nil, err
	}

	span.SetStatus(codes.Ok, "storage group updated successfully")
	return result, nil
}
