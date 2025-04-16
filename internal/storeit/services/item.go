package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/repositories"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type ItemService struct {
	repo *repositories.ItemRepository
}

func NewItemService(repo *repositories.ItemRepository) *ItemService {
	return &ItemService{
		repo: repo,
	}
}

func (s *ItemService) Create(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	item.ID = orgID
	return s.repo.CreateItemWithVariants(ctx, item)
}

func (s *ItemService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	return s.repo.GetItems(ctx, orgID)
}

func (s *ItemService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	return s.repo.GetItem(ctx, orgID, id)
}

func (s *ItemService) Update(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	exists, err := s.repo.IsItemExists(ctx, orgID, item.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrItemNotFound
	}

	// Get existing variants to determine which ones to delete
	existingItem, err := s.repo.GetItem(ctx, orgID, item.ID)
	if err != nil {
		return nil, err
	}

	// Create a map of variant IDs that will remain after the update
	remainingVariants := make(map[uuid.UUID]bool)
	if item.Variants != nil {
		for _, v := range *item.Variants {
			remainingVariants[v.ID] = true
		}
	}

	// Mark variants for deletion that are not in the update
	if existingItem.Variants != nil {
		for _, v := range *existingItem.Variants {
			if !remainingVariants[v.ID] {
				if err := s.repo.DeleteItemVariant(ctx, item.ID, v.ID); err != nil {
					return nil, fmt.Errorf("failed to delete variant: %w", err)
				}
			}
		}
	}

	return s.repo.UpdateItem(ctx, item)
}

func (s *ItemService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteItem(ctx, id)
}

func (s *ItemService) IsItemExists(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (bool, error) {
	return s.repo.IsItemExists(ctx, orgID, id)
}
