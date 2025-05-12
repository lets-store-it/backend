package item

import (
	"context"

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

	if !validateResult.IsAllowed {
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

	return fullItem, nil
}

func (uc *ItemUseCase) GetItemsAll(ctx context.Context) ([]*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemsAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) GetItemById(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemByID(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItem(ctx context.Context, item *models.Item) (*models.Item, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
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
		return nil, usecases.ErrNotAuthorized
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
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemVariantById(ctx, validateResult.OrgID, id, variantId)
}

func (uc *ItemUseCase) GetItemVariants(ctx context.Context, id uuid.UUID) ([]*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemVariantsAll(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) UpdateItemVariant(ctx context.Context, variant *models.ItemVariant) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
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
		return nil, usecases.ErrNotAuthorized
	}

	itemInstance.OrgID = validateResult.OrgID

	createdInstance, err := uc.service.CreateItemInstance(ctx, itemInstance)
	if err != nil {
		return nil, err
	}

	// cellPath, err := uc.service.GetCellPath(ctx, validateResult.OrgID, createdInstance.CellID)
	// if err != nil {
	// 	return nil, err
	// }

	return createdInstance, nil
}

func (uc *ItemUseCase) GetItemInstances(ctx context.Context, id uuid.UUID) ([]*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemInstances(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) GetItemInstanceById(ctx context.Context, id uuid.UUID) (*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemInstanceById(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) GetItemInstancesAll(ctx context.Context) ([]*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetItemInstancesAll(ctx, validateResult.OrgID)
}

func (uc *ItemUseCase) UpdateItemInstance(ctx context.Context, itemInstance *models.ItemInstance) (*models.ItemInstance, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	itemInstance.OrgID = validateResult.OrgID

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
