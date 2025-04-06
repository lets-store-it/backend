package services

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/repositories"
)

type UnitService struct {
	OrgRepository repositories.OrganizationRepository
}

const (
	DefaultLimit  = 10
	DefaultOffset = 0
)

// func toApiOrg(org *database.Org) *api.Org {
// 	return &api.Org{
// 		Id:        api.NewOptUUID(org.ID.UUIDValue()
// 		Name:      org.Name,
// 		Subdomain: org.Subdomain,
// 	}
// }

// GetOrgs implements api.Handler.
func (s *UnitService) GetOrgs(ctx context.Context, params api.GetOrgsParams) (*api.GetOrgsOK, error) {
	orgs, err := s.OrgRepository.GetOrgs(ctx, params.Limit.Or(DefaultLimit), params.Offset.Or(DefaultOffset))
	if err != nil {
		return nil, err
	}

	items := make([]api.GetOrgsOKItemsItem, 0, len(orgs)) // Initialize with capacity
	for _, org := range orgs {
		items = append(items, api.GetOrgsOKItemsItem{
			ID:        api.NewOptUUID(org.ID),
			Name:      org.Name,
			Subdomain: org.Subdomain,
			CreatedAt: api.NewNilDateTime(org.AuditFields.CreatedAt),
			UpdatedAt: api.NewNilDateTime(org.AuditFields.UpdatedAt),
			CreatedBy: api.NewNilUUID(org.AuditFields.CreatedBy),
			UpdatedBy: api.NewNilUUID(org.AuditFields.UpdatedBy),
		})
	}

	return &api.GetOrgsOK{
		Items: items,
	}, nil
}

// CreateOrg implements api.Handler.
func (s *UnitService) CreateOrg(ctx context.Context, req *api.Org) (*api.Org, error) {
	org, err := s.OrgRepository.CreateOrg(ctx, req.Name, req.Subdomain)
	if err != nil {
		return nil, err
	}

	return &api.Org{
		ID:        api.NewOptUUID(org.ID),
		Name:      org.Name,
		Subdomain: org.Subdomain,
	}, nil
}

// DeleteOrg implements api.Handler.
func (s *UnitService) DeleteOrg(ctx context.Context, params api.DeleteOrgParams) error {
	panic("unimplemented")
}

// GetOrgById implements api.Handler.
func (s *UnitService) GetOrgById(ctx context.Context, params api.GetOrgByIdParams) (*api.Org, error) {
	panic("unimplemented")
}

// UpdateOrg implements api.Handler.
func (s *UnitService) UpdateOrg(ctx context.Context, req *api.Org, params api.UpdateOrgParams) (*api.Org, error) {
	panic("unimplemented")
}

// NewError implements api.Handler.
func (s *UnitService) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			Message: err.Error(),
		},
	}
}

func (s *UnitService) CreateUnit(ctx context.Context, request *api.Unit) (*api.Unit, error) {
	return nil, nil
}

func (s *UnitService) DeleteUnit(ctx context.Context, params api.DeleteUnitParams) error {
	return nil
}

func (s *UnitService) GetUnitById(ctx context.Context, params api.GetUnitByIdParams) (*api.GetUnitByIdOK, error) {
	return nil, nil
}

func (s *UnitService) GetUnits(ctx context.Context, params api.GetUnitsParams) (*api.GetUnitsOK, error) {
	return nil, nil
}

func (s *UnitService) UpdateUnit(ctx context.Context, request *api.Unit, params api.UpdateUnitParams) (*api.Unit, error) {
	return nil, nil
}
