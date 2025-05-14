package item

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/item"
	"github.com/let-store-it/backend/internal/usecases"
)

type ItemUseCase struct {
	service     *item.ItemService
	authService *auth.AuthService
}

type ItemUseCaseConfig struct {
	Service     *item.ItemService
	AuthService *auth.AuthService
}

func New(config ItemUseCaseConfig) *ItemUseCase {
	if config.Service == nil || config.AuthService == nil {
		panic("Service and AuthService are required")
	}
	return &ItemUseCase{
		service:     config.Service,
		authService: config.AuthService,
	}
}

func (uc *ItemUseCase) CreateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	createdItem, err := uc.service.CreateItem(ctx, validateResult.OrgID, item)
	if err != nil {
		return nil, err
	}

	fullItem, err := uc.service.GetItemByID(ctx, validateResult.OrgID, createdItem.ID)
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

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemsAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) GetItemById(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemByID(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	updatedItem, err := uc.service.UpdateItem(ctx, validateResult.OrgID, item)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (uc *ItemUseCase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteItem(ctx, validateResult.OrgID, id)
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

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	createdVariant, err := uc.service.CreateItemVariant(ctx, validateResult.OrgID, variant)
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

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemVariantById(ctx, validateResult.OrgID, id, variantId)
}

func (uc *ItemUseCase) GetItemVariants(ctx context.Context, id uuid.UUID) ([]*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemVariantsAll(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItemVariant(ctx context.Context, variant *models.ItemVariant) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	updatedVariant, err := uc.service.UpdateItemVariant(ctx, validateResult.OrgID, variant)
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

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteItemVariant(ctx, validateResult.OrgID, id, variantId)
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

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	itemInstance.OrgID = validateResult.OrgID

	createdInstance, err := uc.service.CreateItemInstance(ctx, itemInstance)
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

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemInstances(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) GetItemInstanceById(ctx context.Context, id uuid.UUID) (*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemInstanceById(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) GetItemInstancesAll(ctx context.Context) ([]*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetItemInstancesAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) UpdateItemInstance(ctx context.Context, instanceId uuid.UUID, variantId uuid.UUID, cellId *uuid.UUID) (*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	itemInstance, err := uc.service.GetItemInstanceById(ctx, validateResult.OrgID, instanceId)
	if err != nil {
		return nil, err
	}	

	itemInstance.VariantID = variantId
	itemInstance.CellID = cellId

	updatedInstance, err := uc.service.UpdateItemInstance(ctx, validateResult.OrgID, itemInstance)
	if err != nil {
		return nil, err
	}

	return updatedInstance, nil
}

func (uc *ItemUseCase) DeleteItemInstance(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteItemInstance(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}
