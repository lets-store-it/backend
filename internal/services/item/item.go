package item

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

var (
	ErrItemNotFound        = errors.New("item not found")
	ErrInvalidItem         = errors.New("invalid item")
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidItemID       = errors.New("invalid item ID")
	ErrInvalidVariant      = errors.New("invalid variant")
)

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}

func toItem(item database.Item) (*models.Item, error) {
	id := uuidFromPgx(item.ID)
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

func toItemVariant(variant database.ItemVariant) (*models.ItemVariant, error) {
	id := uuidFromPgx(variant.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert variant: %w", ErrInvalidVariant)
	}
	itemID := uuidFromPgx(variant.ItemID)
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

type ItemService struct {
	queries *database.Queries
	pgxPool *pgxpool.Pool
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool) *ItemService {
	return &ItemService{
		queries: queries,
		pgxPool: pgxPool,
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

	var description pgtype.Text
	if item.Description != nil {
		description = pgtype.Text{String: *item.Description, Valid: true}
	}

	createdItem, err := qtx.CreateItem(ctx, database.CreateItemParams{
		Name:        item.Name,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	if item.Variants != nil {
		createdVariants := make([]models.ItemVariant, len(*item.Variants))

		for i, variant := range *item.Variants {
			var article pgtype.Text
			if variant.Article != nil {
				article = pgtype.Text{String: *variant.Article, Valid: true}
			}

			var ean13 pgtype.Int4
			if variant.EAN13 != nil {
				ean13 = pgtype.Int4{Int32: int32(*variant.EAN13), Valid: true}
			}

			createdVariant, err := qtx.CreateItemVariant(ctx, database.CreateItemVariantParams{
				ItemID:  createdItem.ID,
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
			})

			if err != nil {
				return nil, fmt.Errorf("failed to create item variant: %w", err)
			}

			createdVariantModel, err := toItemVariant(createdVariant)
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

	results, err := s.queries.GetItems(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	itemsModels := make([]*models.Item, len(results))

	for i, item := range results {
		variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
			OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
			ItemID: pgtype.UUID{Bytes: item.ID.Bytes, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get item variants: %w", err)
		}
		itemVariants := make([]models.ItemVariant, len(variants))
		for j, variant := range variants {
			itemVariant, err := toItemVariant(variant)
			if err != nil {
				return nil, fmt.Errorf("failed to convert variant: %w", err)
			}
			itemVariants[j] = *itemVariant
		}
		itemModel, err := toItem(item)
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
		ID:    pgtype.UUID{Bytes: id, Valid: true},
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		ItemID: pgtype.UUID{Bytes: item.ID.Bytes, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item variants: %w", err)
	}

	itemVariants := make([]models.ItemVariant, len(variants))
	for j, variant := range variants {
		itemVariant, err := toItemVariant(variant)
		if err != nil {
			return nil, fmt.Errorf("failed to convert variant: %w", err)
		}
		itemVariants[j] = *itemVariant
	}

	itemModel, err := toItem(item)
	if err != nil {
		return nil, fmt.Errorf("failed to convert item: %w", err)
	}

	itemModel.Variants = &itemVariants

	return itemModel, nil
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
					ItemID: pgtype.UUID{Bytes: item.ID, Valid: true},
					ID:     pgtype.UUID{Bytes: v.ID, Valid: true},
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

	var description pgtype.Text
	if item.Description != nil {
		description = pgtype.Text{String: *item.Description, Valid: true}
	}

	_, err = qtx.UpdateItem(ctx, database.UpdateItemParams{
		ID:          pgtype.UUID{Bytes: item.ID, Valid: true},
		Name:        item.Name,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	if item.Variants != nil {
		for _, variant := range *item.Variants {
			var article pgtype.Text
			if variant.Article != nil {
				article = pgtype.Text{String: *variant.Article, Valid: true}
			}

			var ean13 pgtype.Int4
			if variant.EAN13 != nil {
				ean13 = pgtype.Int4{Int32: int32(*variant.EAN13), Valid: true}
			}

			_, err = qtx.UpdateItemVariant(ctx, database.UpdateItemVariantParams{
				ItemID:  pgtype.UUID{Bytes: item.ID, Valid: true},
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
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
		ID: pgtype.UUID{Bytes: id, Valid: true},
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
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("failed to check item existence: %w", err)
	}
	return exists, nil
}

func (s *ItemService) CreateInstance(ctx context.Context, orgID uuid.UUID, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	panic("unimplemented")
}
