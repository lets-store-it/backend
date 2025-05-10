package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/storage"
)

func toItemModel(item sqlc.Item) (*models.Item, error) {
	return &models.Item{
		ID:          database.UUIDFromPgx(item.ID),
		Name:        item.Name,
		Description: database.PgTextPtrFromPgx(item.Description),
	}, nil
}

func toItemVariantModel(variant sqlc.ItemVariant) (*models.ItemVariant, error) {
	return &models.ItemVariant{
		ID:        database.UUIDFromPgx(variant.ID),
		ItemID:    database.UUIDFromPgx(variant.ItemID),
		Name:      variant.Name,
		Article:   database.PgTextPtrFromPgx(variant.Article),
		EAN13:     database.PgInt32PtrFromPgx(variant.Ean13),
		CreatedAt: variant.CreatedAt.Time,
		DeletedAt: database.PgTimePtrFromPgx(variant.DeletedAt),
	}, nil
}

func toItemInstanceModel(instance sqlc.ItemInstance) (*models.ItemInstance, error) {
	return &models.ItemInstance{
		ID:        database.UUIDFromPgx(instance.ID),
		OrgID:     database.UUIDFromPgx(instance.OrgID),
		ItemID:    database.UUIDFromPgx(instance.ItemID),
		VariantID: database.UUIDFromPgx(instance.VariantID),
		CellID:    database.UUIDFromPgx(instance.CellID),
		Status:    models.ItemInstanceStatus(instance.Status),
	}, nil
}

type toFullItemModelParams struct {
	item           sqlc.Item
	variants       []sqlc.ItemVariant
	instances      []sqlc.ItemInstance
	storageService *storage.StorageService
	orgID          uuid.UUID
}

func toFullItemModel(ctx context.Context, params toFullItemModelParams) (*models.Item, error) {
	itemModel, err := toItemModel(params.item)
	if err != nil {
		return nil, fmt.Errorf("failed to convert base item: %w", err)
	}

	itemVariants := make([]models.ItemVariant, len(params.variants))
	for i, variant := range params.variants {
		variantModel, err := toItemVariantModel(variant)
		if err != nil {
			return nil, fmt.Errorf("failed to convert variant: %w", err)
		}
		itemVariants[i] = *variantModel
	}
	itemModel.Variants = &itemVariants

	itemInstances := make([]models.ItemInstance, len(params.instances))
	for i, instance := range params.instances {
		instanceModel, err := toItemInstanceModel(instance)
		if err != nil {
			return nil, fmt.Errorf("failed to convert instance: %w", err)
		}

		var instanceVariant *models.ItemVariant
		for _, variant := range itemVariants {
			if variant.ID == instanceModel.VariantID {
				instanceVariant = &variant
				break
			}
		}

		if instanceVariant == nil {
			return nil, errors.New("failed to find variant for instance")
		}

		instanceModel.Variant = instanceVariant

		if instance.CellID.Valid {
			cell, err := params.storageService.GetCellByID(ctx, params.orgID, instanceModel.CellID)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					continue
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
