package usecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/services"
)

type OrganizationUseCase struct {
	service *services.OrganizationService
}

func NewOrganizationUseCase(service *services.OrganizationService) *OrganizationUseCase {
	return &OrganizationUseCase{
		service: service,
	}
}

func (uc *OrganizationUseCase) validateOrganizationData(name, subdomain string) error {
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

func (uc *OrganizationUseCase) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	if err := uc.validateOrganizationData(name, subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Create(ctx, name, subdomain)
}

func (uc *OrganizationUseCase) GetAll(ctx context.Context) ([]*models.Organization, error) {
	return uc.service.GetAll(ctx)
}

func (uc *OrganizationUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	org, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}

func (uc *OrganizationUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.service.Delete(ctx, id)
}

func (uc *OrganizationUseCase) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	if err := uc.validateOrganizationData(org.Name, org.Subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, org)
}

func (uc *OrganizationUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	// Get current organization
	org, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		org.Name = name
	}
	if subdomain, ok := updates["subdomain"].(string); ok {
		org.Subdomain = subdomain
	}

	// Validate updated data
	if err := uc.validateOrganizationData(org.Name, org.Subdomain); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, org)
}

func (uc *OrganizationUseCase) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return uc.service.IsOrganizationExistsByID(ctx, id)
}
