package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/usecases"
	"github.com/let-store-it/backend/internal/utils"
)

type ItemUseCase struct {
	service *item.ItemService
	authService *auth.AuthService
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
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.Create(ctx, validateResult.OrgID, item)
}

func (uc *ItemUseCase) GetAll(ctx context.Context) ([]*models.Item, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetByID(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) Update(ctx context.Context, orgId uuid.UUID, item *models.Item) (*models.Item, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.Update(ctx, validateResult.OrgID, item)
}

func (uc *ItemUseCase) Patch(ctx context.Context, orgId uuid.UUID, id uuid.UUID, updates map[string]interface{}) (*models.Item, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	// Get the existing item first
	item, err := uc.service.GetByID(ctx, validateResult.OrgID, id)
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
	return uc.service.Update(ctx, validateResult.OrgID, item)
}

func (uc *ItemUseCase) Delete(ctx context.Context, orgId uuid.UUID, id uuid.UUID) error {
	return uc.service.Delete(ctx, orgId, id)
}
