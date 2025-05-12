package item

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/usecases"
)

type ItemUseCase struct {
	service      *item.ItemService
	authService  *auth.AuthService
	auditService *audit.AuditService
}

type ItemUseCaseConfig struct {
	Service      *item.ItemService
	AuthService  *auth.AuthService
	AuditService *audit.AuditService
}

func New(config ItemUseCaseConfig) *ItemUseCase {
	return &ItemUseCase{
		service:      config.Service,
		authService:  config.AuthService,
		auditService: config.AuditService,
	}
}

func (uc *ItemUseCase) CreateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	createdItem, err := uc.service.CreateItem(ctx, validateResult.OrgID, item)
	if err != nil {
		return nil, err
	}

	// Fetch the full item with variants
	fullItem, err := uc.service.GetItemByID(ctx, validateResult.OrgID, createdItem.ID)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(fullItem)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeItem,
		TargetObjectID:     fullItem.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return fullItem, nil
}

func (uc *ItemUseCase) GetItemsAll(ctx context.Context) ([]*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemsAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) GetItemById(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemByID(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	// Get the current state for audit
	currentItem, err := uc.service.GetItemByID(ctx, validateResult.OrgID, item.ID)
	if err != nil {
		return nil, err
	}

	prechangeState, err := json.Marshal(currentItem)
	if err != nil {
		return nil, err
	}

	updatedItem, err := uc.service.UpdateItem(ctx, validateResult.OrgID, item)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(updatedItem)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeItem,
		TargetObjectID:     updatedItem.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (uc *ItemUseCase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelWorker)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	// Get the current state for audit
	currentItem, err := uc.service.GetItemByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	prechangeState, err := json.Marshal(currentItem)
	if err != nil {
		return err
	}

	err = uc.service.DeleteItem(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionDelete,
		TargetObjectTypeId: models.ObjectTypeItem,
		TargetObjectID:     id,
		PrechangeState:     prechangeState,
		PostchangeState:    nil,
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *ItemUseCase) CreateItemVariant(ctx context.Context, variant *models.ItemVariant) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	createdVariant, err := uc.service.CreateItemVariant(ctx, validateResult.OrgID, variant)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(createdVariant)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeItemVariant,
		TargetObjectID:     createdVariant.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return createdVariant, nil
}

func (uc *ItemUseCase) GetItemVariantById(ctx context.Context, id uuid.UUID, variantId uuid.UUID) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemVariantById(ctx, validateResult.OrgID, id, variantId)
}

func (uc *ItemUseCase) GetItemVariants(ctx context.Context, id uuid.UUID) ([]*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemVariantsAll(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItemVariant(ctx context.Context, variant *models.ItemVariant) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	// Get the current state for audit
	currentVariant, err := uc.service.GetItemVariantById(ctx, validateResult.OrgID, variant.ItemID, variant.ID)
	if err != nil {
		return nil, err
	}

	prechangeState, err := json.Marshal(currentVariant)
	if err != nil {
		return nil, err
	}

	updatedVariant, err := uc.service.UpdateItemVariant(ctx, validateResult.OrgID, variant)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(updatedVariant)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeItemVariant,
		TargetObjectID:     updatedVariant.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return updatedVariant, nil
}

func (uc *ItemUseCase) DeleteItemVariant(ctx context.Context, id uuid.UUID, variantId uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	// Get the current state for audit
	currentVariant, err := uc.service.GetItemVariantById(ctx, validateResult.OrgID, id, variantId)
	if err != nil {
		return err
	}

	prechangeState, err := json.Marshal(currentVariant)
	if err != nil {
		return err
	}

	err = uc.service.DeleteItemVariant(ctx, validateResult.OrgID, id, variantId)
	if err != nil {
		return err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionDelete,
		TargetObjectTypeId: models.ObjectTypeItemVariant,
		TargetObjectID:     variantId,
		PrechangeState:     prechangeState,
		PostchangeState:    nil,
	})
	if err != nil {
		return err
	}

	return nil
}

func (uc *ItemUseCase) CreateItemInstance(ctx context.Context, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	itemInstance.OrgID = validateResult.OrgID

	createdInstance, err := uc.service.CreateItemInstance(ctx, itemInstance)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(createdInstance)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeItemInstance,
		TargetObjectID:     createdInstance.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return createdInstance, nil
}

func (uc *ItemUseCase) GetItemInstances(ctx context.Context, id uuid.UUID) ([]*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemInstances(ctx, validateResult.OrgID, id)
}
