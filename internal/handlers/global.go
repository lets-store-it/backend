package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/common"
	auditUC "github.com/let-store-it/backend/internal/usecases/audit"
	authUC "github.com/let-store-it/backend/internal/usecases/auth"
	itemUC "github.com/let-store-it/backend/internal/usecases/item"
	orgUC "github.com/let-store-it/backend/internal/usecases/organization"
	storageUC "github.com/let-store-it/backend/internal/usecases/storage"
	"github.com/ogen-go/ogen/ogenerrors"
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

func (h *RestApiImplementation) NewConflictError(ctx context.Context, message string) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusConflict,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "conflict",
				Message: message,
			},
		},
	}
}

func (h *RestApiImplementation) NewUnauthorizedError(ctx context.Context) *api.DefaultErrorStatusCode {
	return h.NewUnauthorizedErrorWithMessage(ctx, "Unauthorized")
}

func (h *RestApiImplementation) NewUnauthorizedErrorWithMessage(ctx context.Context, message string) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusUnauthorized,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "unauthorized",
				Message: message,
			},
		},
	}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.DefaultErrorStatusCode {
	if errors.Is(err, ErrSessionNotFound) {
		return h.NewUnauthorizedError(ctx)
	}
	if errors.Is(err, ErrSessionRevoked) {
		return h.NewUnauthorizedErrorWithMessage(ctx, "Session was revoked")
	}
	if errors.Is(err, ErrSessionExpired) {
		return h.NewUnauthorizedErrorWithMessage(ctx, "Session expired")
	}

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		return h.NewUnauthorizedError(ctx)
	}

	if errors.Is(err, common.ErrNotAuthorized) {
		return h.NewUnauthorizedError(ctx)
	}

	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "internal_server_error",
				Message: err.Error(),
			},
		},
	}
}
