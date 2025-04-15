package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/internal/storeit/usecases"
)

type RestApiImplementation struct {
	orgUseCase     *usecases.OrganizationUseCase
	orgUnitUseCase *usecases.OrganizationUnitUseCase
}

// CreateStorageSpace implements api.Handler.
func (h *RestApiImplementation) CreateStorageSpace(ctx context.Context, req *api.CreateStorageSpaceRequest) (*api.CreateStorageSpaceResponse, error) {
	panic("unimplemented")
}

// DeleteStorageSpace implements api.Handler.
func (h *RestApiImplementation) DeleteStorageSpace(ctx context.Context, params api.DeleteStorageSpaceParams) error {
	panic("unimplemented")
}

// GetStorageSpaceById implements api.Handler.
func (h *RestApiImplementation) GetStorageSpaceById(ctx context.Context, params api.GetStorageSpaceByIdParams) (*api.GetStorageSpaceByIdResponse, error) {
	panic("unimplemented")
}

// GetStorageSpaces implements api.Handler.
func (h *RestApiImplementation) GetStorageSpaces(ctx context.Context) (*api.GetStorageSpacesResponse, error) {
	panic("unimplemented")
}

// PatchStorageSpace implements api.Handler.
func (h *RestApiImplementation) PatchStorageSpace(ctx context.Context, req *api.PatchStorageSpaceRequest, params api.PatchStorageSpaceParams) (*api.PatchStorageSpaceResponse, error) {
	panic("unimplemented")
}

// UpdateStorageSpace implements api.Handler.
func (h *RestApiImplementation) UpdateStorageSpace(ctx context.Context, req *api.UpdateStorageSpaceRequest, params api.UpdateStorageSpaceParams) (*api.UpdateStorageSpaceResponse, error) {
	panic("unimplemented")
}

func NewRestApiImplementation(orgUseCase *usecases.OrganizationUseCase, orgUnitUseCase *usecases.OrganizationUnitUseCase) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:     orgUseCase,
		orgUnitUseCase: orgUnitUseCase,
	}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
