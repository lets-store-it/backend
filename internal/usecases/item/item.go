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

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.CreateItem(ctx, validateResult.OrgID, item)
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

	return uc.service.UpdateItem(ctx, validateResult.OrgID, item)
}

func (uc *ItemUseCase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelWorker)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteItem(ctx, validateResult.OrgID, id)
}

func (uc *ItemUseCase) CreateItemVariant(ctx context.Context, variant *models.ItemVariant) (*models.ItemVariant, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.CreateItemVariant(ctx, validateResult.OrgID, variant)
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

	return uc.service.UpdateItemVariant(ctx, validateResult.OrgID, variant)
}

func (uc *ItemUseCase) DeleteItemVariant(ctx context.Context, id uuid.UUID, variantId uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteItemVariant(ctx, validateResult.OrgID, id, variantId)
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

	return uc.service.CreateItemInstance(ctx, itemInstance)
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
