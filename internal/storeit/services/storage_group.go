package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrStorageGroupNotFound = errors.New("storage group not found")
)

type StorageGroupRepository interface {
	GetStorageGroupByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error)
	GetStorageGroups(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error)
	CreateStorageGroup(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error)
	DeleteStorageGroup(ctx context.Context, id uuid.UUID) error
	UpdateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error)
	IsStorageGroupExistsForOrganization(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error)
}

type StorageGroupService struct {
	repo StorageGroupRepository
}

func NewStorageGroupService(repo StorageGroupRepository) *StorageGroupService {
	return &StorageGroupService{
		repo: repo,
	}
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
	return s.repo.CreateStorageGroup(ctx, orgID, unitID, parentID, name, alias)
}

func (s *StorageGroupService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	return s.repo.GetStorageGroups(ctx, orgID)
}

func (s *StorageGroupService) GetByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	group, err := s.repo.GetStorageGroupByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrStorageGroupNotFound
	}
	return group, nil
}

func (s *StorageGroupService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteStorageGroup(ctx, id)
}

func (s *StorageGroupService) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	return s.repo.UpdateStorageGroup(ctx, group)
}

func (s *StorageGroupService) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.StorageGroup, error) {
	group, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		group.Name = name
	}
	if alias, ok := updates["alias"].(string); ok {
		group.Alias = alias
	}
	if parentID, ok := updates["parent_id"].(uuid.UUID); ok {
		group.ParentID = parentID
	}

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return s.repo.UpdateStorageGroup(ctx, group)
}

func (s *StorageGroupService) IsStorageGroupExistsForOrganization(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error) {
	return s.repo.IsStorageGroupExistsForOrganization(ctx, orgID, groupID)
}
