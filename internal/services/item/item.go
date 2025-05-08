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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	tracer         trace.Tracer
}

func New(queries *database.Queries, pgxPool *pgxpool.Pool, storageService *storage.StorageService) *ItemService {
	return &ItemService{
		queries:        queries,
		pgxPool:        pgxPool,
		storageService: storageService,
		tracer:         otel.GetTracerProvider().Tracer("item-service"),
	}
}

func (s *ItemService) Create(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "Create")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if item == nil {
		span.SetStatus(codes.Error, "invalid item")
		return nil, ErrInvalidItem
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("name", item.Name),
	)
	if item.Description != nil {
		span.SetAttributes(attribute.String("description", *item.Description))
	}

	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	createdItem, err := qtx.CreateItem(ctx, database.CreateItemParams{
		OrgID:       utils.PgUUID(orgID),
		Name:        item.Name,
		Description: utils.PgTextPtr(item.Description),
	})
	item.ID = *utils.UuidFromPgx(createdItem.ID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create item")
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
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to create item variant")
				return nil, fmt.Errorf("failed to create item variant: %w", err)
			}

			createdVariantModel, err := toItemVariantModel(createdVariant)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to convert created variant")
				return nil, fmt.Errorf("failed to convert created variant: %w", err)
			}
			createdVariants[i] = *createdVariantModel
		}
		item.Variants = &createdVariants
	}

	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "item created successfully")
	return item, nil
}

func (s *ItemService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "GetAll")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}

	span.SetAttributes(attribute.String("org_id", orgID.String()))

	results, err := s.queries.GetItems(ctx, utils.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get items")
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	itemsModels := make([]*models.Item, len(results))

	for i, item := range results {
		variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
			OrgID:  utils.PgUUID(orgID),
			ItemID: utils.PgUUID(uuid.UUID(item.ID.Bytes)),
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get item variants")
			return nil, fmt.Errorf("failed to get item variants: %w", err)
		}
		itemVariants := make([]models.ItemVariant, len(variants))
		for j, variant := range variants {
			itemVariant, err := toItemVariantModel(variant)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to convert variant")
				return nil, fmt.Errorf("failed to convert variant: %w", err)
			}
			itemVariants[j] = *itemVariant
		}
		itemModel, err := toItemModel(item)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert item")
			return nil, fmt.Errorf("failed to convert item: %w", err)
		}
		itemModel.Variants = &itemVariants
		itemsModels[i] = itemModel
	}

	span.SetStatus(codes.Ok, "items retrieved successfully")
	return itemsModels, nil
}

func (s *ItemService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid item ID")
		return nil, ErrInvalidItemID
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("item_id", id.String()),
	)

	item, err := s.queries.GetItem(ctx, database.GetItemParams{
		ID:    utils.PgUUID(id),
		OrgID: utils.PgUUID(orgID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "item not found")
			return nil, ErrItemNotFound
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item")
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	variants, err := s.queries.GetItemVariants(ctx, database.GetItemVariantsParams{
		OrgID:  utils.PgUUID(orgID),
		ItemID: utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item variants")
		return nil, fmt.Errorf("failed to get item variants: %w", err)
	}

	instances, err := s.queries.GetItemInstancesForItem(ctx, database.GetItemInstancesForItemParams{
		OrgID:  utils.PgUUID(orgID),
		ItemID: utils.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item instances")
		return nil, fmt.Errorf("failed to get item instances: %w", err)
	}

	result, err := toFullItemModel(ctx, toFullItemModelParams{
		item:           item,
		variants:       variants,
		instances:      instances,
		storageService: s.storageService,
		orgID:          orgID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert item model")
		return nil, err
	}

	span.SetStatus(codes.Ok, "item retrieved successfully")
	return result, nil
}

func (s *ItemService) Update(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "Update")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if item == nil {
		span.SetStatus(codes.Error, "invalid item")
		return nil, ErrInvalidItem
	}
	if item.ID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid item ID")
		return nil, ErrInvalidItemID
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("item_id", item.ID.String()),
		attribute.String("name", item.Name),
	)
	if item.Description != nil {
		span.SetAttributes(attribute.String("description", *item.Description))
	}

	exists, err := s.IsItemExists(ctx, orgID, item.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to check item existence")
		return nil, fmt.Errorf("failed to check item existence: %w", err)
	}
	if !exists {
		span.SetStatus(codes.Error, "item not found")
		return nil, ErrItemNotFound
	}

	existingItem, err := s.GetByID(ctx, orgID, item.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get existing item")
		return nil, fmt.Errorf("failed to get existing item: %w", err)
	}

	remainingVariants := make(map[uuid.UUID]bool)
	if item.Variants != nil {
		for _, v := range *item.Variants {
			remainingVariants[v.ID] = true
		}
	}

	if existingItem.Variants != nil {
		for _, v := range *existingItem.Variants {
			if !remainingVariants[v.ID] {
				err := s.queries.DeleteItemVariant(ctx, database.DeleteItemVariantParams{
					ItemID: utils.PgUUID(item.ID),
					ID:     utils.PgUUID(v.ID),
				})
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, "failed to delete item variant")
					return nil, fmt.Errorf("failed to delete item variant: %w", err)
				}
			}
		}
	}

	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	_, err = qtx.UpdateItem(ctx, database.UpdateItemParams{
		OrgID:       utils.PgUUID(orgID),
		ID:          utils.PgUUID(item.ID),
		Name:        item.Name,
		Description: utils.PgTextPtr(item.Description),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update item")
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
				span.RecordError(err)
				span.SetStatus(codes.Error, "failed to update item variant")
				return nil, fmt.Errorf("failed to update item variant: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	span.SetStatus(codes.Ok, "item updated successfully")
	return item, nil
}

func (s *ItemService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid item ID")
		return ErrInvalidItemID
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("item_id", id.String()),
	)

	err := s.queries.DeleteItem(ctx, database.DeleteItemParams{
		ID:    utils.PgUUID(id),
		OrgID: utils.PgUUID(orgID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete item")
		return fmt.Errorf("failed to delete item: %w", err)
	}

	span.SetStatus(codes.Ok, "item deleted successfully")
	return nil
}

func (s *ItemService) IsItemExists(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "IsItemExists")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return false, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		span.SetStatus(codes.Error, "invalid item ID")
		return false, ErrInvalidItemID
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("item_id", id.String()),
	)

	exists, err := s.queries.IsItemExists(ctx, database.IsItemExistsParams{
		ID:    utils.PgUUID(id),
		OrgID: utils.PgUUID(orgID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to check item existence")
		return false, fmt.Errorf("failed to check item existence: %w", err)
	}

	span.SetStatus(codes.Ok, "item existence checked successfully")
	return exists, nil
}

func (s *ItemService) CreateInstance(ctx context.Context, orgID uuid.UUID, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	ctx, span := s.tracer.Start(ctx, "CreateInstance")
	defer span.End()

	if orgID == uuid.Nil {
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if itemInstance == nil {
		span.SetStatus(codes.Error, "invalid item instance")
		return nil, ErrInvalidItem
	}

	span.SetAttributes(
		attribute.String("org_id", orgID.String()),
		attribute.String("item_id", itemInstance.ItemID.String()),
		attribute.String("variant_id", itemInstance.VariantID.String()),
		attribute.String("cell_id", itemInstance.CellID.String()),
	)

	createdInstance, err := s.queries.CreateItemInstance(ctx, database.CreateItemInstanceParams{
		OrgID:     utils.PgUUID(orgID),
		ItemID:    utils.PgUUID(itemInstance.ItemID),
		VariantID: utils.PgUUID(itemInstance.VariantID),
		CellID:    utils.PgUUID(itemInstance.CellID),
		Status:    string(itemInstance.Status),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create item instance")
		return nil, fmt.Errorf("failed to create item instance: %w", err)
	}

	result, err := toItemInstanceModel(createdInstance)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert item instance")
		return nil, err
	}

	span.SetStatus(codes.Ok, "item instance created successfully")
	return result, nil
}
