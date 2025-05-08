package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/item"
)

type ItemUseCase struct {
	service *item.ItemService
}

func (uc *ItemUseCase) validateOrganizationAccess(ctx context.Context, itemID uuid.UUID) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	if itemID != uuid.Nil {
		exists, err := uc.service.IsItemExists(ctx, orgID, itemID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to check item ownership: %w", err)
		}
		if !exists {
			return uuid.Nil, item.ErrItemNotFound
		}
	}

	return orgID, nil
}

type ItemUseCaseConfig struct {
	Service *item.ItemService
}

func New(config ItemUseCaseConfig) *ItemUseCase {
	return &ItemUseCase{
		service: config.Service,
	}
}

func (uc *ItemUseCase) Create(ctx context.Context, item *models.Item) (*models.Item, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return uc.service.Create(ctx, orgID, item)
}

func (uc *ItemUseCase) GetAll(ctx context.Context) ([]*models.Item, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return uc.service.GetAll(ctx, orgID)
}

func (uc *ItemUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	return uc.service.GetByID(ctx, orgID, id)
}

func (uc *ItemUseCase) Update(ctx context.Context, orgId uuid.UUID, item *models.Item) (*models.Item, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, item.ID)
	if err != nil {
		return nil, err
	}

	return uc.service.Update(ctx, orgID, item)
}

func (uc *ItemUseCase) Patch(ctx context.Context, orgId uuid.UUID, id uuid.UUID, updates map[string]interface{}) (*models.Item, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get the existing item first
	item, err := uc.service.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		item.Name = name
	}
	if description, ok := updates["description"].(*string); ok {
		item.Description = description
	}

	// Handle variants update
	if variants, ok := updates["variants"].([]interface{}); ok {
		// Create a map of existing variants for quick lookup
		existingVariants := make(map[uuid.UUID]bool)
		if item.Variants != nil {
			for _, v := range *item.Variants {
				existingVariants[v.ID] = true
			}
		}

		// Process new/updated variants
		newVariants := make([]models.ItemVariant, 0, len(variants))
		for _, v := range variants {
			variant, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			var variantID uuid.UUID
			if id, ok := variant["id"].(string); ok {
				variantID, err = uuid.Parse(id)
				if err != nil {
					continue
				}
			} else {
				variantID = uuid.New() // Generate new ID for new variants
			}

			name, _ := variant["name"].(string)
			article, _ := variant["article"].(*string)
			var ean13 *int
			if e, ok := variant["ean13"].(float64); ok {
				e64 := int(e)
				ean13 = &e64
			}

			newVariant := models.ItemVariant{
				ID:      variantID,
				ItemID:  item.ID,
				Name:    name,
				Article: article,
				EAN13:   ean13,
			}
			newVariants = append(newVariants, newVariant)
			delete(existingVariants, variantID) // Remove from map to track which ones to delete
		}

		// If variants array is provided (even if empty), update the item's variants
		item.Variants = &newVariants

		// Any remaining variants in the map should be marked for deletion
		// The service layer should handle the deletion of variants not present in the update
	}

	// Update the item
	return uc.service.Update(ctx, orgID, item)
}

func (uc *ItemUseCase) Delete(ctx context.Context, orgId uuid.UUID, id uuid.UUID) error {
	return uc.service.Delete(ctx, orgId, id)
}
