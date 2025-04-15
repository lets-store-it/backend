package handlers

import (
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

type RestApiImplementation struct {
	orgUseCase     *usecases.OrganizationUseCase
	orgUnitUseCase *usecases.OrganizationUnitUseCase
}

func NewRestApiImplementation(orgUseCase *usecases.OrganizationUseCase, orgUnitUseCase *usecases.OrganizationUnitUseCase) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:     orgUseCase,
		orgUnitUseCase: orgUnitUseCase,
	}
}
