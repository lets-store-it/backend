package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

func convertToDTO(org *models.Organization) api.Organization {
	return api.Organization{
		ID:        org.ID,
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}
}

func convertToDTOs(orgs []*models.Organization) []api.Organization {
	items := make([]api.Organization, 0, len(orgs))
	for _, org := range orgs {
		items = append(items, convertToDTO(org))
	}
	return items
}

func (h *RestApiImplementation) CreateOrganization(ctx context.Context, req *api.CreateOrganizationRequest) (api.CreateOrganizationRes, error) {
	org, err := h.orgUseCase.Create(ctx, req.Name, req.Subdomain)
	if err != nil {
		return nil, err
	}

	return &api.CreateOrganizationResponse{
		Data: convertToDTO(org),
	}, nil
}

func (h *RestApiImplementation) GetOrganizations(ctx context.Context) (api.GetOrganizationsRes, error) {
	orgs, err := h.orgUseCase.GetUsersOrgs(ctx)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationsResponse{
		Data: convertToDTOs(orgs),
	}, nil
}

func (h *RestApiImplementation) GetOrganizationById(ctx context.Context, params api.GetOrganizationByIdParams) (api.GetOrganizationByIdRes, error) {
	org, err := h.orgUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetOrganizationByIdResponse{
		Data: convertToDTO(org),
	}, nil
}

func (h *RestApiImplementation) UpdateOrganization(ctx context.Context, req *api.OrganizationUpdate, params api.UpdateOrganizationParams) (api.UpdateOrganizationRes, error) {
	ctx = context.WithValue(ctx, models.OrganizationIDContextKey, params.ID)

	org := &models.Organization{
		Name: req.Name,
	}

	updatedOrg, err := h.orgUseCase.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	return &api.UpdateOrganizationResponse{
		Data: convertToDTO(updatedOrg),
	}, nil
}

func (h *RestApiImplementation) DeleteOrganization(ctx context.Context, params api.DeleteOrganizationParams) (api.DeleteOrganizationRes, error) {
	ctx = context.WithValue(ctx, models.OrganizationIDContextKey, params.ID)

	err := h.orgUseCase.Delete(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.DefaultNoContent{}, nil
}
