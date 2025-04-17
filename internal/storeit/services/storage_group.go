package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrStorageGroupNotFound = errors.New("storage group not found")
)

type StorageGroupService struct {
	queries *database.Queries
}

func NewStorageGroupService(queries *database.Queries) *StorageGroupService {
	return &StorageGroupService{
		queries: queries,
	}
}

func toStorageGroup(group database.StorageGroup) (*models.StorageGroup, error) {
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
		ParentID: uuidFromPgx(group.ParentID),
		Name:     group.Name,
		Alias:    group.Alias,
	}, nil
}

func (s *StorageGroupService) validateStorageGroupData(name string, alias string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("storage group name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("storage group name is too long (max 100 characters)")
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("storage group alias cannot be empty")
	}
	if len(alias) > 100 {
		return fmt.Errorf("storage group alias is too long (max 100 characters)")
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("storage group alias can only contain letters, numbers, and hyphens (no spaces)")
	}
	return nil
}

func (s *StorageGroupService) Create(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	var parentIDPgx pgtype.UUID
	if parentID != nil {
		parentIDPgx = pgtype.UUID{Bytes: *parentID, Valid: true}
	} else {
		parentIDPgx = pgtype.UUID{Valid: false}
	}

	group, err := s.queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    pgtype.UUID{Bytes: orgID, Valid: true},
		UnitID:   pgtype.UUID{Bytes: unitID, Valid: true},
		ParentID: parentIDPgx,
		Name:     name,
		Alias:    alias,
	})
	if err != nil {
		return nil, err
	}

	return toStorageGroup(group)
}

func (s *StorageGroupService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	groups, err := s.queries.GetActiveStorageGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
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

func (s *StorageGroupService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	group, err := s.queries.GetStorageGroup(ctx, database.GetStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	model, err := toStorageGroup(group)

	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, ErrStorageGroupNotFound
	}
	return model, nil
}

func (s *StorageGroupService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.queries.DeleteStorageGroup(ctx, database.DeleteStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (s *StorageGroupService) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	updatedGroup, err := s.queries.UpdateStorageGroup(ctx, database.UpdateStorageGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, err
	}
	return toStorageGroup(updatedGroup)
}

func (s *StorageGroupService) IsStorageGroupExists(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error) {
	return s.queries.IsStorageGroupExists(ctx, database.IsStorageGroupExistsParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: groupID, Valid: true},
	})
}
