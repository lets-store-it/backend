package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/utils"
)

var (
	ErrItemNotFound        = errors.New("item not found")
	ErrInvalidItem         = errors.New("invalid item")
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidItemID       = errors.New("invalid item ID")
	ErrInvalidVariant      = errors.New("invalid variant")
)

type ItemService struct {
	storageService *storage.StorageService
	queries        *database.Queries
	pgxPool        *pgxpool.Pool
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool, storageService *storage.StorageService) *ItemService {
	return &ItemService{
		queries:        queries,
		pgxPool:        pgxPool,
		storageService: storageService,
	}
}

func (s *ItemService) Create(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if item == nil {
		return nil, ErrInvalidItem
	}

	item.ID = orgID
	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	createdItem, err := qtx.CreateItem(ctx, database.CreateItemParams{
		Name:        item.Name,
		Description: utils.PgTextPtr(item.Description),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	if item.Variants != nil {
		createdVariants := make([]models.ItemVariant, len(*item.Variants))

		for i, variant := range *item.Variants {
			createdVariant, err := qtx.CreateItemVariant(ctx, database.CreateItemVariantParams{
				ItemID:  createdItem.ID,
				Name:    variant.Name,
				Article: utils.PgTextPtr(variant.Article),
				Ean13:   pgtype.Int4{Int32: int32(*variant.EAN13), Valid: variant.EAN13 != nil},
			})

			if err != nil {
				return nil, fmt.Errorf("failed to create item variant: %w", err)
			}

			createdVariantModel, err := toItemVariantModel(createdVariant)
			if err != nil {
				return nil, fmt.Errorf("failed to convert created variant: %w", err)
			}
			createdVariants[i] = *createdVariantModel
		}
		item.Variants = &createdVariants
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return item, nil
}

func (s *ItemService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	results, err := s.queries.GetItems(ctx, utils.PgUUID(orgID))
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	itemsModels := make([]*models.Item, len(results))

	for i, item := range results {
		variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
			OrgID:  utils.PgUUID(orgID),
			ItemID: utils.PgUUID(uuid.UUID(item.ID.Bytes)),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get item variants: %w", err)
		}
		itemVariants := make([]models.ItemVariant, len(variants))
		for j, variant := range variants {
			itemVariant, err := toItemVariantModel(variant)
			if err != nil {
				return nil, fmt.Errorf("failed to convert variant: %w", err)
			}
			itemVariants[j] = *itemVariant
		}
		itemModel, err := toItemModel(item)
		if err != nil {
			return nil, fmt.Errorf("failed to convert item: %w", err)
		}
		itemModel.Variants = &itemVariants
		itemsModels[i] = itemModel
	}
	return itemsModels, nil
}

func (s *ItemService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidItemID
	}

	item, err := s.queries.GetItem(ctx, database.GetItemParams{
		ID:    utils.PgUUID(id),
		OrgID: utils.PgUUID(orgID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
		OrgID:  utils.PgUUID(orgID),
		ItemID: utils.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item variants: %w", err)
	}

	// Get item instances
	instances, err := s.queries.GetItemInstancesForItem(ctx, database.GetItemInstancesForItemParams{
		OrgID:  utils.PgUUID(orgID),
		ItemID: utils.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item instances: %w", err)
	}

	return toFullItemModel(ctx, toFullItemModelParams{
		item:           item,
		variants:       variants,
		instances:      instances,
		storageService: s.storageService,
		orgID:          orgID,
	})
}

func (s *ItemService) Update(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if item == nil {
		return nil, ErrInvalidItem
	}
	if item.ID == uuid.Nil {
		return nil, ErrInvalidItemID
	}

	exists, err := s.IsItemExists(ctx, orgID, item.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check item existence: %w", err)
	}
	if !exists {
		return nil, ErrItemNotFound
	}

	// Get existing variants to determine which ones to delete
	existingItem, err := s.GetByID(ctx, orgID, item.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing item: %w", err)
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
				err := s.queries.DeleteItemVariant(ctx, database.DeleteItemVariantParams{
					ItemID: utils.PgUUID(item.ID),
					ID:     utils.PgUUID(v.ID),
				})
				if err != nil {
					return nil, fmt.Errorf("failed to delete item variant: %w", err)
				}
			}
		}
	}

	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	_, err = qtx.UpdateItem(ctx, database.UpdateItemParams{
		ID:          utils.PgUUID(item.ID),
		Name:        item.Name,
		Description: utils.PgTextPtr(item.Description),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	if item.Variants != nil {
		for _, variant := range *item.Variants {
			_, err = qtx.UpdateItemVariant(ctx, database.UpdateItemVariantParams{
				ItemID:  utils.PgUUID(item.ID),
				Name:    variant.Name,
				Article: utils.PgTextPtr(variant.Article),
				Ean13:   pgtype.Int4{Int32: int32(*variant.EAN13), Valid: variant.EAN13 != nil},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update item variant: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return item, nil
}

func (s *ItemService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidItemID
	}

	err := s.queries.DeleteItem(ctx, database.DeleteItemParams{
		ID: utils.PgUUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrItemNotFound
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (s *ItemService) IsItemExists(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (bool, error) {
	if orgID == uuid.Nil {
		return false, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return false, ErrInvalidItemID
	}

	exists, err := s.queries.IsItemExists(ctx, database.IsItemExistsParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return false, fmt.Errorf("failed to check item existence: %w", err)
	}
	return exists, nil
}

func (s *ItemService) CreateInstance(ctx context.Context, orgID uuid.UUID, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	panic("unimplemented")
}
