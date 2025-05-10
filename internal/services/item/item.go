package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ItemService struct {
	storageService *storage.StorageService
	queries        *sqlc.Queries
	pgxPool        *pgxpool.Pool
	tracer         trace.Tracer
}

type ItemServiceConfig struct {
	Queries        *sqlc.Queries
	PGXPool        *pgxpool.Pool
	StorageService *storage.StorageService
}

func New(config ItemServiceConfig) *ItemService {
	return &ItemService{
		queries:        config.Queries,
		pgxPool:        config.PGXPool,
		storageService: config.StorageService,
		tracer:         otel.GetTracerProvider().Tracer("item-service"),
	}
}

func (s *ItemService) CreateItem(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", orgID.String()),
		attribute.String("item.name", item.Name),
	)

	if item.Description != nil {
		span.SetAttributes(attribute.String("item.description", *item.Description))
	}

	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	createdItem, err := qtx.CreateItem(ctx, sqlc.CreateItemParams{
		OrgID:       database.PgUUID(orgID),
		Name:        item.Name,
		Description: database.PgTextPtr(item.Description),
	})
	item.ID = database.UUIDFromPgx(createdItem.ID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create item")
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	if item.Variants != nil {
		createdVariants := make([]models.ItemVariant, len(*item.Variants))

		for i, variant := range *item.Variants {
			var article string
			if variant.Article != nil {
				article = *variant.Article
			}

			createdVariant, err := qtx.CreateItemVariant(ctx, sqlc.CreateItemVariantParams{
				ItemID:  createdItem.ID,
				Name:    variant.Name,
				Article: database.PgText(article),
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

func (s *ItemService) GetItemsAll(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "GetAll")
	defer span.End()

	span.SetAttributes(attribute.String("org.id", orgID.String()))

	results, err := s.queries.GetItems(ctx, database.PgUUID(orgID))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get items")
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	itemsModels := make([]*models.Item, len(results))

	for i, item := range results {
		variants, err := s.queries.GetItemVariants(ctx, sqlc.GetItemVariantsParams{
			OrgID:  database.PgUUID(orgID),
			ItemID: database.PgUUID(uuid.UUID(item.ID.Bytes)),
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

func (s *ItemService) GetItemByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "GetByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", orgID.String()),
		attribute.String("item.id", id.String()),
	)

	item, err := s.queries.GetItem(ctx, sqlc.GetItemParams{
		ID:    database.PgUUID(id),
		OrgID: database.PgUUID(orgID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			span.SetStatus(codes.Error, "item not found")
			return nil, services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item")
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	variants, err := s.queries.GetItemVariants(ctx, sqlc.GetItemVariantsParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item variants")
		return nil, fmt.Errorf("failed to get item variants: %w", err)
	}

	instances, err := s.queries.GetItemInstancesForItem(ctx, sqlc.GetItemInstancesForItemParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(id),
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

func (s *ItemService) UpdateItem(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	ctx, span := s.tracer.Start(ctx, "Update")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", orgID.String()),
		attribute.String("item.id", item.ID.String()),
		attribute.String("item.name", item.Name),
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
		return nil, services.ErrNotFoundError
	}

	existingItem, err := s.GetItemByID(ctx, orgID, item.ID)
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
				err := s.queries.DeleteItemVariant(ctx, sqlc.DeleteItemVariantParams{
					ItemID: database.PgUUID(item.ID),
					ID:     database.PgUUID(v.ID),
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

	var description string
	if item.Description != nil {
		description = *item.Description
	}

	_, err = qtx.UpdateItem(ctx, sqlc.UpdateItemParams{
		OrgID:       database.PgUUID(orgID),
		ID:          database.PgUUID(item.ID),
		Name:        item.Name,
		Description: database.PgText(description),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update item")
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	if item.Variants != nil {
		for _, variant := range *item.Variants {
			var article string
			if variant.Article != nil {
				article = *variant.Article
			}

			_, err = qtx.UpdateItemVariant(ctx, sqlc.UpdateItemVariantParams{
				ItemID:  database.PgUUID(item.ID),
				Name:    variant.Name,
				Article: database.PgText(article),
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

func (s *ItemService) DeleteItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", orgID.String()),
		attribute.String("item.id", id.String()),
	)

	err := s.queries.DeleteItem(ctx, sqlc.DeleteItemParams{
		ID:    database.PgUUID(id),
		OrgID: database.PgUUID(orgID),
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

	span.SetAttributes(
		attribute.String("org.id", orgID.String()),
		attribute.String("item.id", id.String()),
	)

	exists, err := s.queries.IsItemExists(ctx, sqlc.IsItemExistsParams{
		ID:    database.PgUUID(id),
		OrgID: database.PgUUID(orgID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to check item existence")
		return false, fmt.Errorf("failed to check item existence: %w", err)
	}

	span.SetStatus(codes.Ok, "item existence checked successfully")
	return exists, nil
}

func (s *ItemService) CreateItemVariant(ctx context.Context, orgID uuid.UUID, variant *models.ItemVariant) (*models.ItemVariant, error) {
	ctx, span := s.tracer.Start(ctx, "CreateItemVariant")
	defer span.End()

	createdVariant, err := s.queries.CreateItemVariant(ctx, sqlc.CreateItemVariantParams{
		OrgID:   database.PgUUID(orgID),
		ItemID:  database.PgUUID(variant.ItemID),
		Name:    variant.Name,
		Article: database.PgTextPtr(variant.Article),
		Ean13:   database.PgInt4Ptr(variant.EAN13),
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			span.SetStatus(codes.Error, "unique violation")
			return nil, services.ErrDuplicationError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create item variant")
		return nil, fmt.Errorf("failed to create item variant: %w", err)
	}

	span.SetStatus(codes.Ok, "item variant created successfully")
	return toItemVariantModel(createdVariant)
}

func (s *ItemService) DeleteItemVariant(ctx context.Context, orgID uuid.UUID, id uuid.UUID, variantId uuid.UUID) error {
	ctx, span := s.tracer.Start(ctx, "DeleteItemVariant")
	defer span.End()

	err := s.queries.DeleteItemVariant(ctx, sqlc.DeleteItemVariantParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(id),
		ID:     database.PgUUID(variantId),
	})
	if err != nil {
		if database.IsNotFound(err) {
			span.SetStatus(codes.Error, "item variant not found")
			return services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to delete item variant")
		return fmt.Errorf("failed to delete item variant: %w", err)
	}

	span.SetStatus(codes.Ok, "item variant deleted successfully")
	return nil
}

func (s *ItemService) GetItemVariantsAll(ctx context.Context, orgID uuid.UUID, id uuid.UUID) ([]*models.ItemVariant, error) {
	ctx, span := s.tracer.Start(ctx, "GetItemVariants")
	defer span.End()

	variants, err := s.queries.GetItemVariants(ctx, sqlc.GetItemVariantsParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(id),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item variants")
		return nil, fmt.Errorf("failed to get item variants: %w", err)
	}

	variantsModels := make([]*models.ItemVariant, len(variants))
	for i, variant := range variants {
		variantsModels[i], err = toItemVariantModel(variant)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to convert item variant")
			return nil, fmt.Errorf("failed to convert item variant: %w", err)
		}
	}

	span.SetStatus(codes.Ok, "item variants retrieved successfully")
	return variantsModels, nil
}

func (s *ItemService) GetItemVariantById(ctx context.Context, orgID uuid.UUID, id uuid.UUID, variantId uuid.UUID) (*models.ItemVariant, error) {
	ctx, span := s.tracer.Start(ctx, "GetItemVariantById")
	defer span.End()

	variant, err := s.queries.GetItemVariantById(ctx, sqlc.GetItemVariantByIdParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(id),
		ID:     database.PgUUID(variantId),
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item variant by id")
		return nil, fmt.Errorf("failed to get item variant by id: %w", err)
	}

	return toItemVariantModel(variant)
}

func (s *ItemService) UpdateItemVariant(ctx context.Context, orgID uuid.UUID, variant *models.ItemVariant) (*models.ItemVariant, error) {
	ctx, span := s.tracer.Start(ctx, "UpdateItemVariant")
	defer span.End()

	updatedVariant, err := s.queries.UpdateItemVariant(ctx, sqlc.UpdateItemVariantParams{
		OrgID:   database.PgUUID(orgID),
		ItemID:  database.PgUUID(variant.ItemID),
		ID:      database.PgUUID(variant.ID),
		Name:    variant.Name,
		Article: database.PgTextPtr(variant.Article),
		Ean13:   database.PgInt4Ptr(variant.EAN13),
	})
	if err != nil {
		if database.IsUniqueViolation(err) {
			span.SetStatus(codes.Error, "unique violation")
			return nil, services.ErrDuplicationError
		}
		if database.IsNotFound(err) {
			span.SetStatus(codes.Error, "item variant not found")
			return nil, services.ErrNotFoundError
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to update item variant")
		return nil, fmt.Errorf("failed to update item variant: %w", err)
	}

	span.SetStatus(codes.Ok, "item variant updated successfully")
	return toItemVariantModel(updatedVariant)
}

func (s *ItemService) CreateItemInstance(ctx context.Context, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	ctx, span := s.tracer.Start(ctx, "CreateInstance")
	defer span.End()

	span.SetAttributes(
		attribute.String("org.id", itemInstance.OrgID.String()),
		attribute.String("item.id", itemInstance.ItemID.String()),
		attribute.String("variant.id", itemInstance.VariantID.String()),
		attribute.String("cell.id", itemInstance.CellID.String()),
	)

	createdInstance, err := s.queries.CreateItemInstance(ctx, sqlc.CreateItemInstanceParams{
		OrgID:     database.PgUUID(itemInstance.OrgID),
		ItemID:    database.PgUUID(itemInstance.ItemID),
		VariantID: database.PgUUID(itemInstance.VariantID),
		CellID:    database.PgUUIDPtr(itemInstance.CellID),
		Status:    string(models.ItemInstanceStatusAvailable),
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

func (s *ItemService) GetItemInstances(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID) ([]*models.ItemInstance, error) {
	ctx, span := s.tracer.Start(ctx, "GetItemInstances")
	defer span.End()

	instances, err := s.queries.GetItemInstancesForItem(ctx, sqlc.GetItemInstancesForItemParams{
		OrgID:  database.PgUUID(orgID),
		ItemID: database.PgUUID(itemID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item instances")
		return nil, fmt.Errorf("failed to get item instances: %w", err)
	}

	return toItemInstances(instances), nil
}

func (s *ItemService) GetItemInstanceFull(ctx context.Context, orgID uuid.UUID, instanceID uuid.UUID) (*models.ItemInstance, error) {
	ctx, span := s.tracer.Start(ctx, "GetItemInstanceFull")
	defer span.End()

	instanceDb, err := s.queries.GetItemInstance(ctx, sqlc.GetItemInstanceParams{
		ID:    database.PgUUID(instanceID),
		OrgID: database.PgUUID(orgID),
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item instance")
		return nil, fmt.Errorf("failed to get item instance: %w", err)
	}
	instance := toItemInstance(instanceDb)

	if instance.CellID != nil {
		cell, err := s.storageService.GetCellFull(ctx, orgID, *instance.CellID)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get cell")
			return nil, fmt.Errorf("failed to get cell: %w", err)
		}
		instance.Cell = cell
	}

	variant, err := s.GetItemVariantById(ctx, orgID, instance.ItemID, instance.VariantID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get item variant")
		return nil, fmt.Errorf("failed to get item variant: %w", err)
	}
	instance.Variant = variant

	span.SetStatus(codes.Ok, "item instance retrieved successfully")

	return instance, nil
}
