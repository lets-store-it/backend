package storeit

import (
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/handlers"
	"github.com/let-store-it/backend/internal/storeit/repositories"
	"github.com/let-store-it/backend/internal/storeit/services"
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

func NewRestApiImplementation(queries *database.Queries) *handlers.RestApiImplementation {
	// Repositories
	orgRepo := &repositories.OrganizationRepository{Queries: queries}
	orgUnitRepo := &repositories.OrganizationUnitRepository{Queries: queries}
	storageGroupRepo := &repositories.StorageGroupRepository{Queries: queries}

	// Services
	orgService := services.NewOrganizationService(orgRepo)
	orgUnitService := services.NewOrganizationUnitService(orgUnitRepo)
	storageGroupService := services.NewStorageGroupService(storageGroupRepo)

	// Use cases
	orgUseCase := usecases.NewOrganizationUseCase(orgService)
	orgUnitUseCase := usecases.NewOrganizationUnitUseCase(orgUnitService, orgService)
	storageGroupUseCase := usecases.NewStorageGroupUseCase(storageGroupService, orgService)

	// Handlers
	return handlers.NewRestApiImplementation(orgUseCase, orgUnitUseCase, storageGroupUseCase)
}
