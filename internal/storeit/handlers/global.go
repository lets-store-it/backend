package handlers

import (
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

type RestApiImplementation struct {
	orgUseCase          *usecases.OrganizationUseCase
	orgUnitUseCase      *usecases.OrganizationUnitUseCase
	storageGroupUseCase *usecases.StorageGroupUseCase
	itemUseCase         *usecases.ItemUseCase
}

func NewRestApiImplementation(
	orgUseCase *usecases.OrganizationUseCase,
	orgUnitUseCase *usecases.OrganizationUnitUseCase,
	storageGroupUseCase *usecases.StorageGroupUseCase,
	itemUseCase *usecases.ItemUseCase,
) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:          orgUseCase,
		orgUnitUseCase:      orgUnitUseCase,
		storageGroupUseCase: storageGroupUseCase,
		itemUseCase:         itemUseCase,
	}
}
