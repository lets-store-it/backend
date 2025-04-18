package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/storeit/models"
)

func convertUnitToDTO(unit *models.OrganizationUnit) api.Unit {
	var address api.NilString
	if unit.Address == nil {
		address.SetToNull()
	} else {
		address.SetTo(*unit.Address)
	}

	return api.Unit{
		ID:      unit.ID,
		Name:    unit.Name,
		Alias:   unit.Alias,
		Address: address,
	}
}

// GetOrganizationUnits implements api.Handler.
func (h *RestApiImplementation) GetOrganizationUnits(ctx context.Context) (*api.GetOrganizationUnitsResponse, error) {
	units, err := h.orgUnitUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.Unit, 0, len(units))
	for _, unit := range units {
		items = append(items, convertUnitToDTO(unit))
	}

	return &api.GetOrganizationUnitsResponse{
		Data: items,
	}, nil
}

// CreateUnit implements api.Handler.
func (h *RestApiImplementation) CreateUnit(ctx context.Context, req *api.CreateOrganizationUnitRequest) (*api.CreateOrganizationUnitResponse, error) {
	unit, err := h.orgUnitUseCase.Create(ctx, req.Name, req.Alias, req.Address.Value)
	if err != nil {
		return nil, err
	}

	unitDTO := convertUnitToDTO(unit)
	return &api.CreateOrganizationUnitResponse{
		Data: unitDTO,
	}, nil
}

// DeleteOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) DeleteOrganizationUnit(ctx context.Context, params api.DeleteOrganizationUnitParams) error {
	return h.orgUnitUseCase.Delete(ctx, params.ID)
}

// GetOrganizationUnitById implements api.Handler.
func (h *RestApiImplementation) GetOrganizationUnitById(ctx context.Context, params api.GetOrganizationUnitByIdParams) (*api.GetOrganizationUnitByIdResponse, error) {
	unit, err := h.orgUnitUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationUnitByIdResponse{
		Data: convertUnitToDTO(unit),
	}, nil
}

// PatchOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) PatchOrganizationUnit(ctx context.Context, req *api.PatchOrganizationUnitRequest, params api.PatchOrganizationUnitParams) (*api.PatchOrganizationUnitResponse, error) {
	updates := make(map[string]interface{})

	if req.Name.IsSet() {
		updates["name"] = req.Name.Value
	}
	if req.Address.IsSet() {
		updates["address"] = req.Address.Value
	}
	if req.Alias.IsSet() {
		updates["alias"] = req.Alias.Value
	}

	unit, err := h.orgUnitUseCase.Patch(ctx, params.ID, updates)
	if err != nil {
		return nil, err
	}

	return &api.PatchOrganizationUnitResponse{
		Data: []api.Unit{convertUnitToDTO(unit)},
	}, nil
}

// UpdateOrganizationUnit implements api.Handler.
func (h *RestApiImplementation) UpdateOrganizationUnit(ctx context.Context, req *api.UpdateOrganizationUnitRequest, params api.UpdateOrganizationUnitParams) (*api.UpdateOrganizationUnitResponse, error) {
	var address *string
	if req.Address.IsSet() {
		address = &req.Address.Value
	}
	unit := &models.OrganizationUnit{
		ID:      params.ID,
		Name:    req.Name,
		Alias:   req.Alias,
		Address: address,
	}

	updatedUnit, err := h.orgUnitUseCase.Update(ctx, unit)
	if err != nil {
		return nil, err
	}

	return &api.UpdateOrganizationUnitResponse{
		Data: []api.Unit{convertUnitToDTO(updatedUnit)},
	}, nil
}
