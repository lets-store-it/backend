package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

type StorageGroupRepository struct {
	Queries *database.Queries
}

func toStorageGroup(group database.StorageSpace) (*models.StorageGroup, error) {
	id := uuidFromPgx(group.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	unitID := uuidFromPgx(group.UnitID)
	if unitID == nil {
		return nil, errors.New("unit_id is nil")
	}

	return &models.StorageGroup{
		ID:       *id,
		UnitID:   *unitID,
		ParentID: *uuidFromPgx(group.ParentID),
		Name:     group.Name,
		Alias:    group.Alias.String,
	}, nil
}

func (r *StorageGroupRepository) GetStorageGroupByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	group, err := r.Queries.GetStorageGroupById(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return toStorageGroup(group)
}

func (r *StorageGroupRepository) GetStorageGroups(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	groups, err := r.Queries.GetOrganizationStorageGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*models.StorageGroup, len(groups))
	for i, group := range groups {
		result[i], err = toStorageGroup(group)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *StorageGroupRepository) CreateStorageGroup(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	var parentIDPgx pgtype.UUID
	if parentID != nil {
		parentIDPgx = pgtype.UUID{Bytes: *parentID, Valid: true}
	}

	group, err := r.Queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    pgtype.UUID{Bytes: orgID, Valid: true},
		UnitID:   pgtype.UUID{Bytes: unitID, Valid: true},
		ParentID: parentIDPgx,
		Name:     name,
		Alias:    pgtype.Text{String: alias, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return toStorageGroup(group)
}

func (r *StorageGroupRepository) DeleteStorageGroup(ctx context.Context, id uuid.UUID) error {
	return r.Queries.DeleteStorageGroup(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *StorageGroupRepository) UpdateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	updatedGroup, err := r.Queries.UpdateStorageGroup(ctx, database.UpdateStorageGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: pgtype.Text{String: group.Alias, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return toStorageGroup(updatedGroup)
}

func (r *StorageGroupRepository) IsStorageGroupExistsForOrganization(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error) {
	return r.Queries.IsStorageGroupExistsForOrganization(ctx, database.IsStorageGroupExistsForOrganizationParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: groupID, Valid: true},
	})
}
