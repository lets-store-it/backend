package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *StorageService) validateStorageGroupData(name, alias string) error {
	if err := s.validateName(name); err != nil {
		return err
	}
	return s.validateAlias(alias)
}

func (s *StorageService) CreateStorageGroup(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "CreateStorageGroup",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("unit_id", unitID.String()),
			attribute.String("name", name),
			attribute.String("alias", alias),
		),
	)
	defer span.End()

	if err := s.validateStorageGroupData(name, alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if parentID != nil {
		span.SetAttributes(attribute.String("parent_id", parentID.String()))
	}

	var parentUUID pgtype.UUID
	if parentID != nil {
		parentUUID = database.PgUUID(*parentID)
	}

	group, err := s.queries.CreateStorageGroup(ctx, sqlc.CreateStorageGroupParams{
		OrgID:    database.PgUUID(orgID),
		UnitID:   database.PgUUID(unitID),
		ParentID: parentUUID,
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
	ctx, span := s.tracer.Start(ctx, "GetAllStorageGroups",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
		),
	)
	defer span.End()

	groups, err := s.queries.GetStorageGroups(ctx, database.PgUUID(orgID))
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
	ctx, span := s.tracer.Start(ctx, "GetStorageGroupByID",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("group_id", id.String()),
		),
	)
	defer span.End()

	group, err := s.queries.GetStorageGroup(ctx, sqlc.GetStorageGroupParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
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
	ctx, span := s.tracer.Start(ctx, "DeleteStorageGroup",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("group_id", id.String()),
		),
	)
	defer span.End()

	err := s.queries.DeleteStorageGroup(ctx, sqlc.DeleteStorageGroupParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
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
	ctx, span := s.tracer.Start(ctx, "UpdateStorageGroup",
		trace.WithAttributes(
			attribute.String("org_id", group.OrgID.String()),
			attribute.String("group_id", group.ID.String()),
		),
	)
	defer span.End()

	if group == nil {
		span.SetStatus(codes.Error, "invalid storage group: nil")
		return nil, ErrInvalidStorageGroup
	}

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateStorageGroup(ctx, sqlc.UpdateStorageGroupParams{
		ID:     database.PgUUID(group.ID),
		OrgID:  database.PgUUID(group.OrgID),
		UnitID: database.PgUUID(group.UnitID),
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
