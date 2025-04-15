package handlers

import (
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

type RestApiImplementation struct {
	orgUseCase          *usecases.OrganizationUseCase
	orgUnitUseCase      *usecases.OrganizationUnitUseCase
	storageGroupUseCase *usecases.StorageGroupUseCase
}

func NewRestApiImplementation(
	orgUseCase *usecases.OrganizationUseCase,
	orgUnitUseCase *usecases.OrganizationUnitUseCase,
	storageGroupUseCase *usecases.StorageGroupUseCase,
) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:          orgUseCase,
		orgUnitUseCase:      orgUnitUseCase,
		storageGroupUseCase: storageGroupUseCase,
	}
}
