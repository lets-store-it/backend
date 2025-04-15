package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/repositories"
)

var (
	ErrStorageGroupNotFound = errors.New("storage group not found")
)

type StorageGroupService struct {
	repo *repositories.StorageGroupRepository
}

func NewStorageGroupService(repo *repositories.StorageGroupRepository) *StorageGroupService {
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
	if err := s.validateStorageGroupData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return s.repo.CreateStorageGroup(ctx, orgID, unitID, parentID, name, alias)
}

func (s *StorageGroupService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	return s.repo.GetStorageGroups(ctx, orgID)
}

func (s *StorageGroupService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	group, err := s.repo.GetStorageGroup(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, ErrStorageGroupNotFound
	}
	return group, nil
}

func (s *StorageGroupService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteStorageGroup(ctx, orgID, id)
}

func (s *StorageGroupService) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return s.repo.UpdateStorageGroup(ctx, group)
}

func (s *StorageGroupService) IsStorageGroupExists(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error) {
	return s.repo.IsStorageGroupExists(ctx, orgID, groupID)
}
