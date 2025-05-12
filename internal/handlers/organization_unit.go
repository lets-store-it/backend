package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
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
		Alias:   api.StorageAlias(unit.Alias),
		Address: address,
	}
}

func (h *RestApiImplementation) GetOrganizationUnits(ctx context.Context) (api.GetOrganizationUnitsRes, error) {
	units, err := h.orgUnitUseCase.GetAllUnits(ctx)
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

func (h *RestApiImplementation) CreateUnit(ctx context.Context, req *api.UnitBase) (api.CreateUnitRes, error) {
	unit, err := h.orgUnitUseCase.CreateUnit(ctx, req.Name, string(req.Alias), req.Address.Value)
	if err != nil {
		return nil, err
	}

	unitDTO := convertUnitToDTO(unit)
	return &api.CreateOrganizationUnitResponse{
		Data: unitDTO,
	}, nil
}

func (h *RestApiImplementation) DeleteOrganizationUnit(ctx context.Context, params api.DeleteOrganizationUnitParams) (api.DeleteOrganizationUnitRes, error) {
	err := h.orgUnitUseCase.DeleteUnit(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.DefaultNoContent{}, nil
}

func (h *RestApiImplementation) GetOrganizationUnitById(ctx context.Context, params api.GetOrganizationUnitByIdParams) (api.GetOrganizationUnitByIdRes, error) {
	unit, err := h.orgUnitUseCase.GetUnitByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationUnitByIdResponse{
		Data: convertUnitToDTO(unit),
	}, nil
}

func (h *RestApiImplementation) UpdateOrganizationUnit(ctx context.Context, req *api.UnitBase, params api.UpdateOrganizationUnitParams) (api.UpdateOrganizationUnitRes, error) {
	var address *string
	if req.Address.IsSet() {
		address = &req.Address.Value
	}
	unit := &models.OrganizationUnit{
		ID:      params.ID,
		Name:    req.Name,
		Alias:   string(req.Alias),
		Address: address,
	}

	updatedUnit, err := h.orgUnitUseCase.UpdateUnit(ctx, unit)
	if err != nil {
		return nil, err
	}

	return &api.UpdateOrganizationUnitResponse{
		Data: convertUnitToDTO(updatedUnit),
	}, nil
}
