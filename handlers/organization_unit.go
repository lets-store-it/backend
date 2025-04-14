package handlers

import (
	"context"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/models"
	"github.com/google/uuid"
)

func convertUnitToDTO(unit *models.OrganizationUnit) api.Unit {
	return api.Unit{
		ID:      api.NewOptUUID(unit.ID),
		Name:    unit.Name,
		Alias:   api.OptString{Value: unit.Alias},
		Address: api.OptNilString{Value: unit.Address},
	}
}

// CreateUnit implements api.Handler.
func (h *RestApiImplementation) CreateUnit(ctx context.Context, req *api.CreateOrganizationUnitRequest) (*api.CreateOrganizationUnitResponse, error) {
	orgID := ctx.Value("organization_id").(uuid.UUID)

	unit, err := h.orgUnitUseCase.Create(ctx, orgID, req.Name, req.Alias.Value, req.Address.Value)
	if err != nil {
		return nil, err
	}

	unitDTO := convertUnitToDTO(unit)
	return &api.CreateOrganizationUnitResponse{
		Data: api.NewOptUnit(unitDTO),
	}, nil
}

// GetOrganizationUnits implements api.Handler.
func (h *RestApiImplementation) GetOrganizationUnits(ctx context.Context) (*api.GetOrganizationUnitsResponse, error) {
	orgID := ctx.Value("organization_id").(uuid.UUID)

	units, err := h.orgUnitUseCase.GetAll(ctx, orgID)
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
	unit := &models.OrganizationUnit{
		ID:      params.ID,
		Name:    req.Name,
		Address: req.Address.Value,
	}

	updatedUnit, err := h.orgUnitUseCase.Update(ctx, unit)
	if err != nil {
		return nil, err
	}

	return &api.UpdateOrganizationUnitResponse{
		Data: []api.Unit{convertUnitToDTO(updatedUnit)},
	}, nil
}
