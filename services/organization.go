package services

import (
	"context"

	"github.com/evevseev/storeit/backend/models"
	"github.com/google/uuid"
)

type OrganizationRepository interface {
	CreateOrganization(ctx context.Context, name string, subdomain string) (*models.Organization, error)
	GetOrganizations(ctx context.Context) ([]*models.Organization, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
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
