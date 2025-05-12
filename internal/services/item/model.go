package item

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toItemVariantModel(variant sqlc.ItemVariant) *models.ItemVariant {
	return &models.ItemVariant{
		ID:        database.UUIDFromPgx(variant.ID),
		ItemID:    database.UUIDFromPgx(variant.ItemID),
		Name:      variant.Name,
		Article:   database.PgTextPtrFromPgx(variant.Article),
		EAN13:     database.PgInt64PtrFromPgx(variant.Ean13),
		CreatedAt: variant.CreatedAt.Time,
		DeletedAt: database.PgTimePtrFromPgx(variant.DeletedAt),
	}
}

type toItemModelParams struct {
	item      sqlc.Item
	variants  []sqlc.ItemVariant
	instances []sqlc.ItemInstance
}

func toItemModel(params toItemModelParams) *models.Item {
	itemModel := &models.Item{
		ID:          database.UUIDFromPgx(params.item.ID),
		Name:        params.item.Name,
		Description: database.PgTextPtrFromPgx(params.item.Description),
	}

	itemVariants := make([]*models.ItemVariant, len(params.variants))
	for i, variant := range params.variants {
		itemVariants[i] = toItemVariantModel(variant)
	}
	if len(itemVariants) > 0 {
		itemModel.Variants = itemVariants
	}

	itemInstances := make([]*models.ItemInstance, len(params.instances))
	for i, instance := range params.instances {
		itemInstances[i] = toItemInstanceModel(instance)
	}
	if len(itemInstances) > 0 {
		itemModel.Instances = itemInstances
	}

	return itemModel
}

func toItemInstance(instance sqlc.ItemInstance) *models.ItemInstance {
	return &models.ItemInstance{
		ID:               database.UUIDFromPgx(instance.ID),
		OrgID:            database.UUIDFromPgx(instance.OrgID),
		ItemID:           database.UUIDFromPgx(instance.ItemID),
		VariantID:        database.UUIDFromPgx(instance.VariantID),
		CellID:           database.UUIDPtrFromPgx(instance.CellID),
		Status:           models.ItemInstanceStatus(instance.Status),
		AffectedByTaskID: database.UUIDPtrFromPgx(instance.AffectedByTaskID),
	}
}

func toItemInstanceModel(instance sqlc.ItemInstance) *models.ItemInstance {
	return &models.ItemInstance{
		ID:        database.UUIDFromPgx(instance.ID),
		OrgID:     database.UUIDFromPgx(instance.OrgID),
		ItemID:    database.UUIDFromPgx(instance.ItemID),
		VariantID: database.UUIDFromPgx(instance.VariantID),
		CellID:    database.UUIDPtrFromPgx(instance.CellID),
		Status:    models.ItemInstanceStatus(instance.Status),
	}
}

func toItemInstancesModel(instances []sqlc.ItemInstance) []*models.ItemInstance {
	instancesModels := make([]*models.ItemInstance, len(instances))
	for i, instance := range instances {
		instancesModels[i] = toItemInstance(instance)
	}
	return instancesModels
}
