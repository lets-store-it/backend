package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/repositories"
)

var (
	ErrOrganizationUnitNotFound = errors.New("organization unit not found")
)

type OrganizationUnitService struct {
	repo *repositories.OrganizationUnitRepository
}

func NewOrganizationUnitService(repo *repositories.OrganizationUnitRepository) *OrganizationUnitService {
	return &OrganizationUnitService{
		repo: repo,
	}
}

func validateOrganizationUnitData(name string, alias string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("organization unit name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("organization unit name is too long (max 100 characters)")
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("organization unit alias cannot be empty")
	}
	if len(alias) > 100 {
		return fmt.Errorf("organization unit alias is too long (max 100 characters)")
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("organization unit alias can only contain letters, numbers, and hyphens (no spaces)")
	}
	return nil
}

func (s *OrganizationUnitService) Create(ctx context.Context, orgID uuid.UUID, name string, alias string, address string) (*models.OrganizationUnit, error) {
	return s.repo.CreateOrganizationUnit(ctx, orgID, name, alias, address)
}

func (s *OrganizationUnitService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.OrganizationUnit, error) {
	return s.repo.GetOrganizationUnits(ctx, orgID)
}

func (s *OrganizationUnitService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.OrganizationUnit, error) {
	unit, err := s.repo.GetOrganizationUnit(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if unit == nil {
		return nil, ErrOrganizationUnitNotFound
	}
	return unit, nil
}

func (s *OrganizationUnitService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.repo.DeleteOrganizationUnit(ctx, orgID, id)
}

func (s *OrganizationUnitService) Update(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	return s.repo.UpdateOrganizationUnit(ctx, unit)
}

func (s *OrganizationUnitService) IsOrganizationUnitExists(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) (bool, error) {
	return s.repo.IsOrganizationUnitExists(ctx, orgID, unitID)
}
