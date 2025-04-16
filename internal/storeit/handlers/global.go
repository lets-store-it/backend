package handlers

import (
	"context"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
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

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.Error{
			ErrorID: "internal_server_error",
			Message: err.Error(),
		},
	}
}
