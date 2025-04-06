package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/repositories"
)

type APIImplementation struct {
	OrganizationRepository repositories.OrganizationRepository
}

const (
	DefaultLimit  = 10
	DefaultOffset = 0
)

func (s *APIImplementation) GetOrgs(ctx context.Context, params api.GetOrgsParams) (*api.OrganizationsPagedResponse, error) {
	orgs, err := s.OrganizationRepository.GetOrgs(ctx, params.Limit.Or(DefaultLimit), params.Offset.Or(DefaultOffset))
	if err != nil {
		return nil, err
	}

	items := make([]api.OrganizationsPagedResponseItemsItem, 0, len(orgs)) // Initialize with capacity
	for _, org := range orgs {
		items = append(items, api.OrganizationsPagedResponseItemsItem{
			ID:        api.NewOptUUID(org.ID),
			Name:      org.Name,
			Subdomain: org.Subdomain,
		})
	}

	return &api.OrganizationsPagedResponse{
		Items: items,
		Metadata: api.PaginationMetadata{
			Total:  int32(len(orgs)),
			Limit:  params.Limit.Or(DefaultLimit),
			Offset: params.Offset.Or(DefaultOffset),
		},
	}, nil
}

func (s *APIImplementation) CreateOrg(ctx context.Context, req *api.Organization) (*api.Organization, error) {
	exists, err := s.OrganizationRepository.IsOrgExists(ctx, req.Name, req.Subdomain)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &api.ErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response: api.Error{
				Message: "Organization with such name or subdomain already exists",
			},
		}
	}

	org, err := s.OrganizationRepository.CreateOrg(ctx, req.Name, req.Subdomain)
	if err != nil {
		return nil, err
	}

	return &api.Organization{
		ID:        api.NewOptUUID(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

func (s *APIImplementation) DeleteOrg(ctx context.Context, params api.DeleteOrgParams) error {
	exists, err := s.OrganizationRepository.IsOrgExistsById(ctx, params.ID)
	if err != nil {
		return err
	}
	if !exists {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response: api.Error{
				Message: "Organization not found",
			},
		}
	}

	err = s.OrganizationRepository.DeleteOrg(ctx, params.ID)
	if err != nil {
		return err
	}

	return nil
}
func (s *APIImplementation) GetOrgById(ctx context.Context, params api.GetOrgByIdParams) (*api.Organization, error) {
	org, err := s.OrganizationRepository.GetOrgById(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	// if org == nil {
	// 	return &api.ErrorStatusCode{
	// 		StatusCode: http.StatusNotFound,
	// 		Response: api.Error{
	// 			Message: "Organization not found",
	// 		},
	// 	}
	// }

	return &api.Organization{
		ID:        api.NewOptUUID(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

func (s *APIImplementation) UpdateOrg(ctx context.Context, req *api.Organization, params api.UpdateOrgParams) (*api.Organization, error) {
	panic("unimplemented")
}

func (s *APIImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Message: err.Error(),
		},
	}
}

func (s *APIImplementation) CreateUnit(ctx context.Context, request *api.Unit) (*api.Unit, error) {
	return nil, nil
}

func (s *APIImplementation) DeleteUnit(ctx context.Context, params api.DeleteUnitParams) error {
	return nil
}

func (s *APIImplementation) GetUnitById(ctx context.Context, params api.GetUnitByIdParams) (*api.GetUnitByIdOK, error) {
	return nil, nil
}

func (s *APIImplementation) GetUnits(ctx context.Context, params api.GetUnitsParams) (*api.GetUnitsOK, error) {
	return nil, nil
}

func (s *APIImplementation) UpdateUnit(ctx context.Context, request *api.Unit, params api.UpdateUnitParams) (*api.Unit, error) {
	return nil, nil
}
