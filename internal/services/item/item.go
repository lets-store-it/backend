package item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/telemetry"
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

	auditService *audit.AuditService
}

type ItemServiceConfig struct {
	Queries        *sqlc.Queries
	PGXPool        *pgxpool.Pool
	StorageService *storage.StorageService
	AuditService   *audit.AuditService
}

func New(config ItemServiceConfig) *ItemService {
	return &ItemService{
		queries:        config.Queries,
		pgxPool:        config.PGXPool,
		storageService: config.StorageService,
		tracer:         otel.GetTracerProvider().Tracer("item-service"),
		auditService:   config.AuditService,
	}
}

func (s *ItemService) CreateItem(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateItem", func(ctx context.Context, span trace.Span) (*models.Item, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("item.name", item.Name),
		)

		if item.Description != nil {
			span.SetAttributes(attribute.String("item.description", *item.Description))
		}

		createdItem, err := s.queries.CreateItem(ctx, sqlc.CreateItemParams{
			OrgID:       database.PgUUID(orgID),
			Name:        item.Name,
			Description: database.PgTextPtr(item.Description),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		item.ID = database.UUIDFromPgx(createdItem.ID)
		item.Variants = []*models.ItemVariant{}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionCreate,
			TargetObjectType: models.ObjectTypeItem,
			TargetObjectID:   item.ID,
			PostchangeState:  item,
		})

		if err != nil {
			return nil, err
		}
		return item, nil
	})
}

func (s *ItemService) GetItemsAll(ctx context.Context, orgID uuid.UUID) ([]*models.Item, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemsAll", func(ctx context.Context, span trace.Span) ([]*models.Item, error) {
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
				return nil, fmt.Errorf("failed to get item variants: %w", err)
			}
			itemModel := toItemModel(toItemModelParams{
				item:     item,
				variants: variants,
			})
			itemsModels[i] = itemModel
		}

		span.SetStatus(codes.Ok, "items retrieved successfully")
		return itemsModels, nil
	})
}

func (s *ItemService) GetItemByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Item, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemByID", func(ctx context.Context, span trace.Span) (*models.Item, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("item.id", id.String()),
		)

		item, err := s.queries.GetItemById(ctx, sqlc.GetItemByIdParams{
			ID:    database.PgUUID(id),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		variants, err := s.GetItemVariantsAll(ctx, orgID, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get item variants: %w", err)
		}

		instances, err := s.GetItemInstances(ctx, orgID, id)

		if err != nil {
			return nil, fmt.Errorf("failed to get item instances: %w", err)
		}

		result := toItemModel(toItemModelParams{
			item: item,
		})
		result.Instances = instances
		result.Variants = variants
		span.SetStatus(codes.Ok, "item retrieved successfully")
		return result, nil
	})
}

func (s *ItemService) UpdateItem(ctx context.Context, orgID uuid.UUID, item *models.Item) (*models.Item, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateItem", func(ctx context.Context, span trace.Span) (*models.Item, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("item.id", item.ID.String()),
			attribute.String("item.name", item.Name),
		)
		if item.Description != nil {
			span.SetAttributes(attribute.String("description", *item.Description))
		}

		existingItem, err := s.GetItemByID(ctx, orgID, item.ID)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		_, err = s.queries.UpdateItem(ctx, sqlc.UpdateItemParams{
			OrgID:       database.PgUUID(orgID),
			ID:          database.PgUUID(item.ID),
			Name:        item.Name,
			Description: database.PgTextPtr(item.Description),
		})
		updatedItemModel, err := s.GetItemByID(ctx, orgID, item.ID)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeItem,
			TargetObjectID:   item.ID,
			PrechangeState:   existingItem,
			PostchangeState:  updatedItemModel,
		})

		return updatedItemModel, nil
	})
}

func (s *ItemService) DeleteItem(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteItem", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("item.id", id.String()),
		)

		itemBeforeDelete, err := s.GetItemByID(ctx, orgID, id)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.queries.DeleteItem(ctx, sqlc.DeleteItemParams{
			ID:    database.PgUUID(id),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionDelete,
			TargetObjectType: models.ObjectTypeItem,
			TargetObjectID:   id,
			PrechangeState:   itemBeforeDelete,
		})

		return nil
	})
}

// Item Variants

func (s *ItemService) CreateItemVariant(ctx context.Context, orgID uuid.UUID, variant *models.ItemVariant) (*models.ItemVariant, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateItemVariant", func(ctx context.Context, span trace.Span) (*models.ItemVariant, error) {
		createdVariant, err := s.queries.CreateItemVariant(ctx, sqlc.CreateItemVariantParams{
			OrgID:   database.PgUUID(orgID),
			ItemID:  database.PgUUID(variant.ItemID),
			Name:    variant.Name,
			Article: database.PgTextPtr(variant.Article),
			Ean13:   database.PgInt8Ptr(variant.EAN13),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		variantModel := toItemVariantModel(createdVariant)
		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionCreate,
			TargetObjectType: models.ObjectTypeItemVariant,
			TargetObjectID:   variantModel.ID,
			PostchangeState:  variantModel,
		})
		if err != nil {
			return nil, err
		}
		return variantModel, nil
	})
}

func (s *ItemService) DeleteItemVariant(ctx context.Context, orgID uuid.UUID, id uuid.UUID, variantId uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteItemVariant", func(ctx context.Context, span trace.Span) error {
		variantBeforeDelete, err := s.GetItemVariantById(ctx, orgID, id, variantId)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.queries.DeleteItemVariant(ctx, sqlc.DeleteItemVariantParams{
			OrgID:  database.PgUUID(orgID),
			ItemID: database.PgUUID(id),
			ID:     database.PgUUID(variantId),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionDelete,
			TargetObjectType: models.ObjectTypeItemVariant,
			TargetObjectID:   variantId,
			PrechangeState:   variantBeforeDelete,
		})

		return nil
	})
}

func (s *ItemService) GetItemVariantsAll(ctx context.Context, orgID uuid.UUID, id uuid.UUID) ([]*models.ItemVariant, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemVariantsAll", func(ctx context.Context, span trace.Span) ([]*models.ItemVariant, error) {
		variants, err := s.queries.GetItemVariants(ctx, sqlc.GetItemVariantsParams{
			OrgID:  database.PgUUID(orgID),
			ItemID: database.PgUUID(id),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		variantsModels := make([]*models.ItemVariant, len(variants))
		for i, variant := range variants {
			variantsModels[i] = toItemVariantModel(variant)
		}

		return variantsModels, nil
	})
}

func (s *ItemService) GetItemVariantById(ctx context.Context, orgID uuid.UUID, id uuid.UUID, variantId uuid.UUID) (*models.ItemVariant, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemVariantById", func(ctx context.Context, span trace.Span) (*models.ItemVariant, error) {
		variant, err := s.queries.GetItemVariantById(ctx, sqlc.GetItemVariantByIdParams{
			OrgID:  database.PgUUID(orgID),
			ItemID: database.PgUUID(id),
			ID:     database.PgUUID(variantId),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		return toItemVariantModel(variant), nil
	})
}

func (s *ItemService) UpdateItemVariant(ctx context.Context, orgID uuid.UUID, variant *models.ItemVariant) (*models.ItemVariant, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateItemVariant", func(ctx context.Context, span trace.Span) (*models.ItemVariant, error) {
		variantBeforeUpdate, err := s.GetItemVariantById(ctx, orgID, variant.ItemID, variant.ID)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		updatedVariant, err := s.queries.UpdateItemVariant(ctx, sqlc.UpdateItemVariantParams{
			OrgID:   database.PgUUID(orgID),
			ItemID:  database.PgUUID(variant.ItemID),
			ID:      database.PgUUID(variant.ID),
			Name:    variant.Name,
			Article: database.PgTextPtr(variant.Article),
			Ean13:   database.PgInt8Ptr(variant.EAN13),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		updatedVariantModel := toItemVariantModel(updatedVariant)
		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeItemVariant,
			TargetObjectID:   variant.ID,
			PrechangeState:   variantBeforeUpdate,
			PostchangeState:  updatedVariantModel,
		})

		return updatedVariantModel, nil
	})
}

func (s *ItemService) CreateItemInstance(ctx context.Context, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "CreateItemInstance", func(ctx context.Context, span trace.Span) (*models.ItemInstance, error) {
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
			Status:    sqlc.ItemInstanceStatus(models.ItemInstanceStatusAvailable),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		model := toItemInstance(createdInstance)
		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionCreate,
			TargetObjectType: models.ObjectTypeItemInstance,
			TargetObjectID:   model.ID,
			PostchangeState:  model,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create audit log: %w", err)
		}

		result, err := s.GetItemInstanceFull(ctx, itemInstance.OrgID, database.UUIDFromPgx(createdInstance.ID))
		if err != nil {
			return nil, err
		}

		return result, nil
	})
}

func (s *ItemService) GetItemInstances(ctx context.Context, orgID uuid.UUID, itemID uuid.UUID) ([]*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemInstances", func(ctx context.Context, span trace.Span) ([]*models.ItemInstance, error) {
		instances, err := s.queries.GetItemInstancesForItem(ctx, sqlc.GetItemInstancesForItemParams{
			OrgID:  database.PgUUID(orgID),
			ItemID: database.PgUUID(itemID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		instancesModels := make([]*models.ItemInstance, len(instances))
		for i, instance := range instances {
			instancesModels[i] = toItemInstance(instance)
			if instancesModels[i].CellID != nil {
				cell, err := s.storageService.GetCellFull(ctx, orgID, database.UUIDFromPgx(instance.CellID))
				if err != nil {
					return nil, err
				}
				instancesModels[i].Cell = cell
			}
			variant, err := s.GetItemVariantById(ctx, orgID, itemID, database.UUIDFromPgx(instance.VariantID))
			if err != nil {
				return nil, err
			}
			instancesModels[i].Variant = variant
		}

		return instancesModels, nil
	})
}

func (s *ItemService) GetItemInstanceFull(ctx context.Context, orgID uuid.UUID, instanceID uuid.UUID) (*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemInstanceFull", func(ctx context.Context, span trace.Span) (*models.ItemInstance, error) {
		instanceDb, err := s.queries.GetItemInstance(ctx, sqlc.GetItemInstanceParams{
			ID:    database.PgUUID(instanceID),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}
		instance := toItemInstance(instanceDb)

		if instance.CellID != nil {
			cell, err := s.storageService.GetCellFull(ctx, orgID, *instance.CellID)
			if err != nil {
				return nil, err
			}
			instance.Cell = cell
		}

		variant, err := s.GetItemVariantById(ctx, orgID, instance.ItemID, instance.VariantID)
		if err != nil {
			return nil, err
		}
		instance.Variant = variant

		item, err := s.GetItemByID(ctx, orgID, instance.ItemID)
		if err != nil {
			return nil, err
		}

		instance.Item = item
		return instance, nil
	})
}

func (s *ItemService) SetItemInstanceStatus(ctx context.Context, itemInstance *models.ItemInstance) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "SetItemInstanceStatus", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", itemInstance.OrgID.String()),
			attribute.String("item.id", itemInstance.ItemID.String()),
			attribute.String("variant.id", itemInstance.VariantID.String()),
			attribute.String("instance.id", itemInstance.ID.String()),
		)

		instanceBeforeUpdate, err := s.GetItemInstanceById(ctx, itemInstance.OrgID, itemInstance.ID)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.queries.SetItemInstanceTaskStatus(ctx, sqlc.SetItemInstanceTaskStatusParams{
			OrgID:            database.PgUUID(itemInstance.OrgID),
			ID:               database.PgUUID(itemInstance.ID),
			Status:           sqlc.ItemInstanceStatus(itemInstance.Status),
			AffectedByTaskID: database.PgUUIDPtr(itemInstance.AffectedByTaskID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		updatedInstance, err := s.GetItemInstanceById(ctx, itemInstance.OrgID, itemInstance.ID)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeItemInstance,
			TargetObjectID:   itemInstance.ID,
			PrechangeState:   instanceBeforeUpdate,
			PostchangeState:  updatedInstance,
		})
		if err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}
		return nil
	})
}

func (s *ItemService) SetInstanceCell(ctx context.Context, orgID uuid.UUID, instanceID uuid.UUID, cellID *uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "SetInstanceCell", func(ctx context.Context, span trace.Span) error {
		instanceBeforeUpdate, err := s.GetItemInstanceById(ctx, orgID, instanceID)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.queries.SetItemInstanceCell(ctx, sqlc.SetItemInstanceCellParams{
			OrgID:  database.PgUUID(orgID),
			ID:     database.PgUUID(instanceID),
			CellID: database.PgUUIDPtr(cellID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		updatedInstance, err := s.GetItemInstanceById(ctx, orgID, instanceID)
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeItemInstance,
			TargetObjectID:   instanceID,
			PrechangeState:   instanceBeforeUpdate,
			PostchangeState:  updatedInstance,
		})

		return nil
	})
}

func (s *ItemService) GetItemInstanceById(ctx context.Context, orgID uuid.UUID, instanceID uuid.UUID) (*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemInstanceById", func(ctx context.Context, span trace.Span) (*models.ItemInstance, error) {
		instance, err := s.GetItemInstanceFull(ctx, orgID, instanceID)
		if err != nil {
			return nil, err
		}

		return instance, nil
	})
}

func (s *ItemService) UpdateItemInstance(ctx context.Context, orgID uuid.UUID, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "UpdateItemInstance", func(ctx context.Context, span trace.Span) (*models.ItemInstance, error) {
		instanceBeforeUpdate, err := s.GetItemInstanceById(ctx, orgID, itemInstance.ID)
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		instance, err := s.queries.UpdateItemInstance(ctx, sqlc.UpdateItemInstanceParams{
			OrgID:     database.PgUUID(orgID),
			ID:        database.PgUUID(itemInstance.ID),
			CellID:    database.PgUUIDPtr(itemInstance.CellID),
			VariantID: database.PgUUID(itemInstance.VariantID),
		})
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		model := toItemInstance(instance)

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionUpdate,
			TargetObjectType: models.ObjectTypeItemInstance,
			TargetObjectID:   itemInstance.ID,
			PrechangeState:   instanceBeforeUpdate,
			PostchangeState:  model,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create audit log: %w", err)
		}

		return model, nil
	})
}

func (s *ItemService) DeleteItemInstance(ctx context.Context, orgID uuid.UUID, instanceID uuid.UUID) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "DeleteItemInstance", func(ctx context.Context, span trace.Span) error {

		instanceBeforeDelete, err := s.queries.GetItemInstance(ctx, sqlc.GetItemInstanceParams{
			ID:    database.PgUUID(instanceID),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.queries.DeleteItemInstance(ctx, sqlc.DeleteItemInstanceParams{
			ID:    database.PgUUID(instanceID),
			OrgID: database.PgUUID(orgID),
		})
		if err != nil {
			return services.MapDbErrorToService(err)
		}

		err = s.auditService.CreateObjectChange(ctx, &models.ObjectChangeCreate{
			Action:           models.ObjectChangeActionDelete,
			TargetObjectType: models.ObjectTypeItemInstance,
			TargetObjectID:   instanceID,
			PrechangeState:   instanceBeforeDelete,
		})
		if err != nil {
			return fmt.Errorf("failed to create audit log: %w", err)
		}

		return nil
	})
}

func (s *ItemService) GetItemInstancesAll(ctx context.Context, orgID uuid.UUID) ([]*models.ItemInstance, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetItemInstancesAll", func(ctx context.Context, span trace.Span) ([]*models.ItemInstance, error) {
		instances, err := s.queries.GetItemInstancesAll(ctx, database.PgUUID(orgID))
		if err != nil {
			return nil, services.MapDbErrorToService(err)
		}

		instancesModels := make([]*models.ItemInstance, len(instances))
		for i, instance := range instances {
			instancesModels[i] = toItemInstance(instance)
			if instancesModels[i].CellID != nil {
				cell, err := s.storageService.GetCellFull(ctx, orgID, database.UUIDFromPgx(instance.CellID))
				if err != nil {
					return nil, err
				}
				instancesModels[i].Cell = cell
			}
			variant, err := s.GetItemVariantById(ctx, orgID, database.UUIDFromPgx(instance.ItemID), database.UUIDFromPgx(instance.VariantID))
			if err != nil {
				return nil, err
			}
			instancesModels[i].Variant = variant

			item, err := s.GetItemByID(ctx, orgID, database.UUIDFromPgx(instance.ItemID))
			if err != nil {
				return nil, err
			}
			instancesModels[i].Item = item

		}

		return instancesModels, nil
	})
}
