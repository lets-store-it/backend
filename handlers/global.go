package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/usecases"
)

type RestApiImplementation struct {
	orgUseCase *usecases.OrganizationUseCase
}

// CreateUnit implements api.Handler.
func (h *RestApiImplementation) CreateUnit(ctx context.Context, req *api.CreateOrganizationUnitRequest) (*api.CreateOrganizationUnitResponse, error) {
	panic("unimplemented")
}

// DeleteOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) DeleteOrganizationUnit(ctx context.Context, params api.DeleteOrganizationUnitParams) error {
	panic("unimplemented")
}

// GetOrganizationUnitById implements api.Handler.
func (h *RestApiImplementation) GetOrganizationUnitById(ctx context.Context, params api.GetOrganizationUnitByIdParams) (*api.GetOrganizationUnitByIdResponse, error) {
	panic("unimplemented")
}

// GetOrganizationUnits implements api.Handler.
func (h *RestApiImplementation) GetOrganizationUnits(ctx context.Context) (*api.GetOrganizationUnitsResponse, error) {
	panic("unimplemented")
}

// PatchOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) PatchOrganizationUnit(ctx context.Context, req *api.PatchOrganizationUnitRequest, params api.PatchOrganizationUnitParams) (*api.PatchOrganizationUnitResponse, error) {
	panic("unimplemented")
}

// UpdateOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) UpdateOrganizationUnit(ctx context.Context, req *api.UpdateOrganizationUnitRequest, params api.UpdateOrganizationUnitParams) (*api.UpdateOrganizationUnitResponse, error) {
	panic("unimplemented")
}

func NewRestApiImplementation(orgUseCase *usecases.OrganizationUseCase) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase: orgUseCase,
	}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
