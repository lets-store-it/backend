package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/generated/api"
	"github.com/evevseev/storeit/backend/usecases"
)

type RestApiImplementation struct {
	orgUseCase *usecases.OrganizationUseCase
}

func NewRestApiImplementation(orgUseCase *usecases.OrganizationUseCase) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase: orgUseCase,
	}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
