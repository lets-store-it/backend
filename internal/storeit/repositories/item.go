package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

type ItemRepository struct {
	queries *database.Queries
	dbConn  *pgx.Conn
}

func NewItemRepository(queries *database.Queries, dbConn *pgx.Conn) *ItemRepository {
	return &ItemRepository{
		queries: queries,
		dbConn:  dbConn,
	}
}

func toItem(item database.Item) (*models.Item, error) {
	id := uuidFromPgx(item.ID)
	if id == nil {
		return nil, errors.New("id is nil")
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
		return nil, errors.New("id is nil")
	}
	itemID := uuidFromPgx(variant.ItemID)
	if itemID == nil {
		return nil, errors.New("item_id is nil")
	}

	var article *string
	if variant.Article.Valid {
		article = &variant.Article.String
	}

	var ean13 *int64
	if variant.Ean13.Valid {
		inInt64 := int64(variant.Ean13.Int32)
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

func (r *ItemRepository) CreateItemWithVariants(ctx context.Context, item *models.Item) (*models.Item, error) {
	tx, err := r.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	var description pgtype.Text
	if item.Description != nil {
		description = pgtype.Text{String: *item.Description, Valid: true}
	} else {
		description = pgtype.Text{Valid: false}
	}

	createdItem, err := qtx.CreateItem(ctx, database.CreateItemParams{
		Name:        item.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	if item.Variants != nil {
		createdVariants := make([]models.ItemVariant, len(*item.Variants))

		for i, variant := range *item.Variants {
			var article pgtype.Text
			if variant.Article != nil {
				article = pgtype.Text{String: *variant.Article, Valid: true}
			} else {
				article = pgtype.Text{Valid: false}
			}

			var ean13 pgtype.Int4
			if variant.EAN13 != nil {
				ean13 = pgtype.Int4{Int32: int32(*variant.EAN13), Valid: true}
			} else {
				ean13 = pgtype.Int4{Valid: false}
			}

			createdVariant, err := qtx.CreateItemVariant(ctx, database.CreateItemVariantParams{
				ItemID:  createdItem.ID,
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
			})

			if err != nil {
				return nil, err
			}

			createdVariantModel, err := toItemVariant(createdVariant)
			if err != nil {
				return nil, err
			}
			createdVariants[i] = *createdVariantModel
		}
		item.Variants = &createdVariants
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) GetItems(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	results, err := r.queries.GetActiveItems(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}

	itemsModels := make([]*models.Item, len(results))

	for i, item := range results {
		variants, err := r.queries.GetItemVariants(ctx, pgtype.UUID{Bytes: item.ID.Bytes, Valid: true})
		if err != nil {
			return nil, err
		}
		itemVariants := make([]models.ItemVariant, len(variants))
		for j, variant := range variants {
			itemVariant, err := toItemVariant(variant)
			if err != nil {
				return nil, err
			}
			itemVariants[j] = *itemVariant
		}
		itemModel, err := toItem(item)
		if err != nil {
			return nil, err
		}
		itemModel.Variants = &itemVariants
		itemsModels[i] = itemModel
	}
	return itemsModels, nil
}

func (r *ItemRepository) GetItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	item, err := r.queries.GetItem(ctx, database.GetItemParams{
		ID:    pgtype.UUID{Bytes: id, Valid: true},
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	variants, err := r.queries.GetItemVariants(ctx, pgtype.UUID{Bytes: item.ID.Bytes, Valid: true})
	if err != nil {
		return nil, err
	}

	itemVariants := make([]models.ItemVariant, len(variants))
	for j, variant := range variants {
		itemVariant, err := toItemVariant(variant)
		if err != nil {
			return nil, err
		}
		itemVariants[j] = *itemVariant
	}

	itemModel, err := toItem(item)
	if err != nil {
		return nil, err
	}

	itemModel.Variants = &itemVariants

	return itemModel, nil
}

func (r *ItemRepository) UpdateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	tx, err := r.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := r.queries.WithTx(tx)

	var description pgtype.Text
	if item.Description != nil {
		description = pgtype.Text{String: *item.Description, Valid: true}
	} else {
		description = pgtype.Text{Valid: false}
	}

	_, err = qtx.UpdateItem(ctx, database.UpdateItemParams{
		ID:          pgtype.UUID{Bytes: item.ID, Valid: true},
		Name:        item.Name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	if item.Variants != nil {
		for _, variant := range *item.Variants {
			var article pgtype.Text
			if variant.Article != nil {
				article = pgtype.Text{String: *variant.Article, Valid: true}
			} else {
				article = pgtype.Text{Valid: false}
			}

			var ean13 pgtype.Int4
			if variant.EAN13 != nil {
				ean13 = pgtype.Int4{Int32: int32(*variant.EAN13), Valid: true}
			} else {
				ean13 = pgtype.Int4{Valid: false}
			}

			_, err = qtx.UpdateItemVariant(ctx, database.UpdateItemVariantParams{
				ItemID:  pgtype.UUID{Bytes: item.ID, Valid: true},
				Name:    variant.Name,
				Article: article,
				Ean13:   ean13,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *ItemRepository) DeleteItem(ctx context.Context, id uuid.UUID) error {
	err := r.queries.DeleteItem(ctx, database.DeleteItemParams{
		ID: pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepository) IsItemExists(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (bool, error) {
	exists, err := r.queries.IsItemExists(ctx, database.IsItemExistsParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	return exists, err
}

func (r *ItemRepository) DeleteItemVariant(ctx context.Context, itemID uuid.UUID, variantID uuid.UUID) error {
	return r.queries.DeleteItemVariant(ctx, database.DeleteItemVariantParams{
		ItemID: pgtype.UUID{Bytes: itemID, Valid: true},
		ID:     pgtype.UUID{Bytes: variantID, Valid: true},
	})
}
