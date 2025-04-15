package handlers

import (
	"context"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/internal/storeit/models"
)

func convertToModel(org *api.Organization) *models.Organization {
	return &models.Organization{
		ID:        org.ID.Value,
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}
}

func convertToDTO(org *models.Organization) *api.Organization {
	return &api.Organization{
		ID:        api.NewOptUUID(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}
}

func createOrganizationResponse(org *models.Organization) *api.Organization {
	return convertToDTO(org)
}

func createErrorResponse(statusCode int, message string) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: statusCode,
		Response: api.Error{
			Message: message,
		},
	}
}

// CreateOrganization implements api.Handler.
func (h *RestApiImplementation) CreateOrganization(ctx context.Context, req *api.CreateOrganizationRequest) (*api.CreateOrganizationResponse, error) {
	org, err := h.orgUseCase.Create(ctx, req.Name, req.Subdomain)
	if err != nil {
		return nil, err
	}

	return &api.CreateOrganizationResponse{
		Data: *createOrganizationResponse(org),
	}, nil
}

// GetOrganizations implements api.Handler.
func (h *RestApiImplementation) GetOrganizations(ctx context.Context) (*api.GetOrganizationsResponse, error) {
	orgs, err := h.orgUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.Organization, 0, len(orgs))
	for _, org := range orgs {
		items = append(items, *convertToDTO(org))
	}

	return &api.GetOrganizationsResponse{
		Data: items,
	}, nil
}

// DeleteOrganization implements api.Handler.
func (h *RestApiImplementation) DeleteOrganization(ctx context.Context, params api.DeleteOrganizationParams) error {
	return h.orgUseCase.Delete(ctx, params.ID)
}

// GetOrganizationById implements api.Handler.
func (h *RestApiImplementation) GetOrganizationById(ctx context.Context, params api.GetOrganizationByIdParams) (*api.GetOrganizationByIdResponse, error) {
	org, err := h.orgUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationByIdResponse{
		Data: *convertToDTO(org),
	}, nil
}

// PatchOrganization implements api.Handler.
func (h *RestApiImplementation) PatchOrganization(ctx context.Context, req *api.PatchOrganizationRequest, params api.PatchOrganizationParams) (*api.PatchOrganizationResponse, error) {
	updates := make(map[string]interface{})

	if req.Name.IsSet() {
		updates["name"] = req.Name.Value
	}
	if req.Subdomain.IsSet() {
		updates["subdomain"] = req.Subdomain.Value
	}

	org, err := h.orgUseCase.Patch(ctx, params.ID, updates)
	if err != nil {
		return nil, err
	}

	return &api.PatchOrganizationResponse{
		Data: []api.Organization{*createOrganizationResponse(org)},
	}, nil
}

// UpdateOrganization implements api.Handler.
func (h *RestApiImplementation) UpdateOrganization(ctx context.Context, req *api.UpdateOrganizationRequest, params api.UpdateOrganizationParams) (*api.UpdateOrganizationResponse, error) {
	org := &models.Organization{
		ID:        params.ID,
		Name:      req.Name,
		Subdomain: req.Subdomain,
	}

	updatedOrg, err := h.orgUseCase.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	return &api.UpdateOrganizationResponse{
		Data: []api.Organization{*createOrganizationResponse(updatedOrg)},
	}, nil
}
