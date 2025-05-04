package item

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/utils"
)

func toItemModel(item database.Item) (*models.Item, error) {
	id := utils.UuidFromPgx(item.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert item: %w", ErrInvalidItemID)
	}

	var description *string
	if item.Description.Valid {
		description = &item.Description.String
	}

	return &models.Item{
		ID:          *id,
		Name:        item.Name,
		Description: description,
	}, nil
}

func toItemVariantModel(variant database.ItemVariant) (*models.ItemVariant, error) {
	id := utils.UuidFromPgx(variant.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert variant: %w", ErrInvalidVariant)
	}
	itemID := utils.UuidFromPgx(variant.ItemID)
	if itemID == nil {
		return nil, fmt.Errorf("failed to convert variant: %w", ErrInvalidItemID)
	}

	var article *string
	if variant.Article.Valid {
		article = &variant.Article.String
	}

	var ean13 *int
	if variant.Ean13.Valid {
		inInt64 := int(variant.Ean13.Int32)
		ean13 = &inInt64
	}

	var deletedAt *time.Time
	if variant.DeletedAt.Valid {
		deletedAt = &variant.DeletedAt.Time
	}

	return &models.ItemVariant{
		ID:        *id,
		ItemID:    *itemID,
		Name:      variant.Name,
		Article:   article,
		EAN13:     ean13,
		CreatedAt: variant.CreatedAt.Time,
		DeletedAt: deletedAt,
	}, nil
}

func toItemInstanceModel(instance database.ItemInstance) (*models.ItemInstance, error) {
	id := utils.UuidFromPgx(instance.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert instance: invalid instance ID")
	}

	orgID := utils.UuidFromPgx(instance.OrgID)
	if orgID == nil {
		return nil, fmt.Errorf("failed to convert instance: invalid organization ID")
	}

	itemID := utils.UuidFromPgx(instance.ItemID)
	if itemID == nil {
		return nil, fmt.Errorf("failed to convert instance: invalid item ID")
	}

	variantID := utils.UuidFromPgx(instance.VariantID)
	if variantID == nil {
		return nil, fmt.Errorf("failed to convert instance: invalid variant ID")
	}

	cellID := utils.UuidFromPgx(instance.CellID)
	if cellID == nil {
		return nil, fmt.Errorf("failed to convert instance: invalid cell ID")
	}

	return &models.ItemInstance{
		ID:        *id,
		OrgID:     *orgID,
		ItemID:    *itemID,
		VariantID: *variantID,
		CellID:    *cellID,
		Status:    models.ItemInstanceStatus(instance.Status),
	}, nil
}

type toFullItemModelParams struct {
	item           database.Item
	variants       []database.ItemVariant
	instances      []database.ItemInstance
	storageService *storage.StorageService
	orgID          uuid.UUID
}

func toFullItemModel(ctx context.Context, params toFullItemModelParams) (*models.Item, error) {
	// Convert base item
	itemModel, err := toItemModel(params.item)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base item: %w", err)
	}

	// Convert variants
	itemVariants := make([]models.ItemVariant, len(params.variants))
	for i, variant := range params.variants {
		variantModel, err := toItemVariantModel(variant)
		if err != nil {
			return nil, fmt.Errorf("failed to convert variant: %w", err)
		}
		itemVariants[i] = *variantModel
	}
	itemModel.Variants = &itemVariants

	// Convert instances
	itemInstances := make([]models.ItemInstance, len(params.instances))
	for i, instance := range params.instances {
		instanceModel, err := toItemInstanceModel(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to convert instance: %w", err)
		}

		// Find the variant for this instance
		var instanceVariant *models.ItemVariant
		for _, variant := range itemVariants {
			if variant.ID == instanceModel.VariantID {
				instanceVariant = &variant
				break
			}
		}

		if instanceVariant == nil {
			return nil, fmt.Errorf("failed to find variant for instance: %w", ErrInvalidVariant)
		}

		instanceModel.Variant = instanceVariant

		// If instance has a cell, get cell information
		if instance.CellID.Valid {
			cell, err := params.storageService.GetCellByID(ctx, params.orgID, instanceModel.CellID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					continue // Skip if cell not found
				}
				return nil, fmt.Errorf("failed to get cell: %w", err)
			}

			instanceModel.Cell = cell

			cellPath, err := params.storageService.GetCellPath(ctx, params.orgID, cell.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get cell path: %w", err)
			}
			instanceModel.Cell.Path = &cellPath
		}

		itemInstances[i] = *instanceModel
	}
	itemModel.Instances = &itemInstances

	return itemModel, nil
}
