package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	auditUC "github.com/let-store-it/backend/internal/usecases/audit"
	authUC "github.com/let-store-it/backend/internal/usecases/auth"
	itemUC "github.com/let-store-it/backend/internal/usecases/item"
	orgUC "github.com/let-store-it/backend/internal/usecases/organization"
	storageUC "github.com/let-store-it/backend/internal/usecases/storage"
)

type RestApiImplementation struct {
	orgUseCase          *orgUC.OrganizationUseCase
	orgUnitUseCase      *orgUC.OrganizationUseCase
	storageGroupUseCase *storageUC.StorageUseCase
	itemUseCase         *itemUC.ItemUseCase
	authUseCase         *authUC.AuthUseCase
	auditUseCase        *auditUC.AuditUseCase
}

// CreateInstanceForItem implements api.Handler.
func (h *RestApiImplementation) CreateInstanceForItem(ctx context.Context, req *api.CreateInstanceForItemRequest, params api.CreateInstanceForItemParams) (api.CreateInstanceForItemRes, error) {
	panic("unimplemented")
}

// DeleteInstanceById implements api.Handler.
func (h *RestApiImplementation) DeleteInstanceById(ctx context.Context, params api.DeleteInstanceByIdParams) (api.DeleteInstanceByIdRes, error) {
	panic("unimplemented")
	return &api.DeleteInstanceByIdOK{}, nil
}

// GetInstances implements api.Handler.
func (h *RestApiImplementation) GetInstances(ctx context.Context) (api.GetInstancesRes, error) {
	panic("unimplemented")
	// return api.GetInstancesByItemIdResponse, nil
}

// GetInstancesByItemId implements api.Handler.
func (h *RestApiImplementation) GetInstancesByItemId(ctx context.Context, params api.GetInstancesByItemIdParams) (api.GetInstancesByItemIdRes, error) {
	panic("unimplemented")

}

// PatchCurrentUser implements api.Handler.
func (h *RestApiImplementation) PatchCurrentUser(ctx context.Context, req *api.PatchCurrentUserRequest) (api.PatchCurrentUserRes, error) {
	panic("unimplemented")
}

// PutCurrentUser implements api.Handler.
func (h *RestApiImplementation) PutCurrentUser(ctx context.Context, req *api.UpdateCurrentUserRequest) (api.PutCurrentUserRes, error) {
	panic("unimplemented")
}

func NewRestApiImplementation(
	orgUseCase *orgUC.OrganizationUseCase,
	orgUnitUseCase *orgUC.OrganizationUseCase,
	storageGroupUseCase *storageUC.StorageUseCase,
	itemUseCase *itemUC.ItemUseCase,
	authUseCase *authUC.AuthUseCase,
	auditUseCase *auditUC.AuditUseCase,
) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:          orgUseCase,
		orgUnitUseCase:      orgUnitUseCase,
		storageGroupUseCase: storageGroupUseCase,
		itemUseCase:         itemUseCase,
		authUseCase:         authUseCase,
		auditUseCase:        auditUseCase,
	}
}
