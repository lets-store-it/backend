package handlers

import (
	"context"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/usecases"
)

type RestApiImplementation struct {
	orgUseCase          *usecases.OrganizationUseCase
	orgUnitUseCase      *usecases.OrganizationUnitUseCase
	storageGroupUseCase *usecases.StorageUseCase
	itemUseCase         *usecases.ItemUseCase
	authUseCase         *usecases.AuthUseCase
	auditUseCase        *usecases.AuditUseCase
}



// CreateInstanceForItem implements api.Handler.
func (h *RestApiImplementation) CreateInstanceForItem(ctx context.Context, req *api.CreateInstanceForItemRequest, params api.CreateInstanceForItemParams) (*api.CreateInstanceForItemResponse, error) {
	panic("unimplemented")
}

// DeleteInstanceById implements api.Handler.
func (h *RestApiImplementation) DeleteInstanceById(ctx context.Context, params api.DeleteInstanceByIdParams) error {
	panic("unimplemented")
}

// GetInstances implements api.Handler.
func (h *RestApiImplementation) GetInstances(ctx context.Context) (*api.GetInstancesResponse, error) {
	panic("unimplemented")
}

// GetInstancesByItemId implements api.Handler.
func (h *RestApiImplementation) GetInstancesByItemId(ctx context.Context, params api.GetInstancesByItemIdParams) (*api.GetInstancesByItemIdResponse, error) {
	panic("unimplemented")
}

// PatchCurrentUser implements api.Handler.
func (h *RestApiImplementation) PatchCurrentUser(ctx context.Context, req *api.PatchCurrentUserRequest) (*api.GetCurrentUserResponse, error) {
	panic("unimplemented")
}

// PutCurrentUser implements api.Handler.
func (h *RestApiImplementation) PutCurrentUser(ctx context.Context, req *api.UpdateCurrentUserRequest) (*api.GetCurrentUserResponse, error) {
	panic("unimplemented")
}

func NewRestApiImplementation(
	orgUseCase *usecases.OrganizationUseCase,
	orgUnitUseCase *usecases.OrganizationUnitUseCase,
	storageGroupUseCase *usecases.StorageUseCase,
	itemUseCase *usecases.ItemUseCase,
	authUseCase *usecases.AuthUseCase,
	auditUseCase *usecases.AuditUseCase,
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

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			ErrorID: "internal_server_error",
			Message: err.Error(),
		},
	}
}
