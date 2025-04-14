package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/usecases"
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

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
