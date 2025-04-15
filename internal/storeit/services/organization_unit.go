package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrOrganizationUnitNotFound = errors.New("organization unit not found")
)

type OrganizationUnitRepository interface {
	GetOrganizationUnitByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error)
	GetOrganizationUnits(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error)
	CreateOrganizationUnit(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error)
	DeleteOrganizationUnit(ctx context.Context, id uuid.UUID) error
	UpdateOrganizationUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error)
	IsOrganizationUnitExistsForOrganization(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) (bool, error)
}

type OrganizationUnitService struct {
	repo OrganizationUnitRepository
}

func NewOrganizationUnitService(repo OrganizationUnitRepository) *OrganizationUnitService {
	return &OrganizationUnitService{
		repo: repo,
	}
}

func (s *OrganizationUnitService) Create(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	return s.repo.CreateOrganizationUnit(ctx, orgID, name, alias, address)
}

func (s *OrganizationUnitService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	return s.repo.GetOrganizationUnits(ctx, orgID)
}

func (s *OrganizationUnitService) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	unit, err := s.repo.GetOrganizationUnitByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrOrganizationUnitNotFound
	}
	return unit, nil
}

func (s *OrganizationUnitService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteOrganizationUnit(ctx, id)
}

func (s *OrganizationUnitService) Update(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	return s.repo.UpdateOrganizationUnit(ctx, unit)
}

func (s *OrganizationUnitService) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationUnit, error) {
	unit, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		unit.Name = name
	}
	if alias, ok := updates["alias"].(string); ok {
		unit.Alias = alias
	}
	if address, ok := updates["address"].(string); ok {
		unit.Address = &address
	}

	return s.repo.UpdateOrganizationUnit(ctx, unit)
}

func (s *OrganizationUnitService) IsOrganizationUnitExistsForOrganization(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) (bool, error) {
	return s.repo.IsOrganizationUnitExistsForOrganization(ctx, orgID, unitID)
}
