package handlers

import (
	"context"
	"net/http"

	"github.com/evevseev/storeit/backend/repositories"
	"github.com/evevseev/storeit/backend/generated/api"
)

type RestApiImplementation struct {
	repo repositories.OrganizationRepository
}

func NewGlobalHandler(repo repositories.OrganizationRepository) *RestApiImplementation {
	return &RestApiImplementation{repo: repo}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
