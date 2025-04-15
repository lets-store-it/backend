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
	ErrOrganizationNotFound = errors.New("organization not found")
)

type OrganizationService struct {
	repo *repositories.OrganizationRepository
}

func NewOrganizationService(repo *repositories.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repo: repo,
	}
}

func validateOrganizationData(name, subdomain string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("organization name is too long (max 100 characters)")
	}
	if strings.TrimSpace(subdomain) == "" {
		return fmt.Errorf("subdomain cannot be empty")
	}
	if len(subdomain) > 63 {
		return fmt.Errorf("subdomain is too long (max 63 characters)")
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", subdomain)
	if !matched {
		return fmt.Errorf("subdomain can only contain lowercase letters, numbers, and hyphens")
	}
	return nil
}

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	if err := validateOrganizationData(name, subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return s.repo.CreateOrganization(ctx, name, subdomain)
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]*models.Organization, error) {
	return s.repo.GetOrganizations(ctx)
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.repo.GetOrganization(ctx, id)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteOrganization(ctx, id)
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	exists, err := s.repo.IsOrganizationExists(ctx, org.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrOrganizationNotFound
	}
	if err := validateOrganizationData(org.Name, org.Subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return s.repo.UpdateOrganization(ctx, org)
}

func (s *OrganizationService) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.IsOrganizationExists(ctx, id)
}
