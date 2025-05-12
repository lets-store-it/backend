package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/telemetry"
	"github.com/let-store-it/backend/internal/utils"
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

func (s *StorageService) CreateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	return telemetry.WithTrace(
		ctx,
		s.tracer,
		"CreateStorageGroup",
		func(ctx context.Context, span trace.Span) (*models.StorageGroup, error) {
			span.SetAttributes(
				attribute.String("org.id", group.OrgID.String()),
				attribute.String("unit.id", group.UnitID.String()),
				attribute.String("storage_group.name", group.Name),
				attribute.String("storage_group.alias", group.Alias),
				attribute.String("storage_group.parent_id", utils.SafeUUIDString(group.ParentID)),
			)

			if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
				return nil, fmt.Errorf("validation failed: %w", err)
			}

			sqlGroup, err := s.queries.CreateStorageGroup(ctx, sqlc.CreateStorageGroupParams{
				OrgID:    database.PgUUID(group.OrgID),
				UnitID:   database.PgUUID(group.UnitID),
				ParentID: database.PgUUIDPtr(group.ParentID),
				Name:     group.Name,
				Alias:    group.Alias,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create storage group: %w", err)
			}

			result := toStorageGroupModel(sqlGroup)
			return result, nil
		},
	)
}

func (s *StorageService) GetAllStorageGroups(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllStorageGroups",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
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
		result[i] = toStorageGroupModel(group)
	}

	span.SetStatus(codes.Ok, "storage groups retrieved successfully")
	return result, nil
}

func (s *StorageService) GetStorageGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	ctx, span := s.tracer.Start(ctx, "GetStorageGroupByID",
		trace.WithAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("storage_group.id", id.String()),
		),
	)
	defer span.End()

	group, err := s.queries.GetStorageGroupById(ctx, sqlc.GetStorageGroupByIdParams{
		OrgID: database.PgUUID(orgID),
		ID:    database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get storage group")
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	model := toStorageGroupModel(group)

	span.SetStatus(codes.Ok, "storage group retrieved successfully")
	return model, nil
}

func (s *StorageService) DeleteStorageGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteStorageGroup", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("storage_group.id", id.String()),
		)

		err := s.queries.DeleteStorageGroup(ctx, sqlc.DeleteStorageGroupParams{
			OrgID: database.PgUUID(orgID),
			ID:    database.PgUUID(id),
		})
		if err != nil {
			return fmt.Errorf("failed to delete storage group: %w", err)
		}

		return nil
	})
}

func (s *StorageService) UpdateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateStorageGroup", func(ctx context.Context, span trace.Span) (*models.StorageGroup, error) {
		if group == nil {
			return nil, common.ErrValidationError
		}
		span.SetAttributes(
			attribute.String("org.id", group.OrgID.String()),
			attribute.String("storage_group.id", group.ID.String()),
			attribute.String("unit.id", group.UnitID.String()),
			attribute.String("storage_group.name", group.Name),
			attribute.String("storage_group.alias", group.Alias),
		)

		if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
			return nil, err
		}

		updatedGroup, err := s.queries.UpdateStorageGroup(ctx, sqlc.UpdateStorageGroupParams{
			ID:     database.PgUUID(group.ID),
			OrgID:  database.PgUUID(group.OrgID),
			UnitID: database.PgUUID(group.UnitID),
			Name:   group.Name,
			Alias:  group.Alias,
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		result := toStorageGroupModel(updatedGroup)
		return result, nil
	})
}
