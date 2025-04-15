package services

import (
	"context"
	"errors"

	"github.com/evevseev/storeit/backend/internal/storeit/models"
	"github.com/google/uuid"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, subdomain string) (*models.Organization, error)
	GetOrganizations(ctx context.Context) ([]*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	UpdateOrganization(ctx context.Context, org *models.Organization) (*models.Organization, error)
	IsOrganizationExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type OrganizationService struct {
	repo OrganizationRepository
}

func NewOrganizationService(repo OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repo: repo,
	}
}

func (s *OrganizationService) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	return s.repo.CreateOrganization(ctx, name, subdomain)
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]*models.Organization, error) {
	return s.repo.GetOrganizations(ctx)
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.repo.GetOrganizationByID(ctx, id)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteOrganization(ctx, id)
}

func (s *OrganizationService) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	exists, err := s.repo.IsOrganizationExistsByID(ctx, org.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrOrganizationNotFound
	}
	return s.repo.UpdateOrganization(ctx, org)
}

func (s *OrganizationService) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	org, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply partial updates
	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				org.Name = name
			}
		case "subdomain":
			if subdomain, ok := value.(string); ok {
				org.Subdomain = subdomain
			}
		}
	}

	return s.repo.UpdateOrganization(ctx, org)
}

func (s *OrganizationService) IsOrganizationExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	return s.repo.IsOrganizationExistsByID(ctx, id)
}
