package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	auditUC "github.com/let-store-it/backend/internal/usecases/audit"
	authUC "github.com/let-store-it/backend/internal/usecases/auth"
	itemUC "github.com/let-store-it/backend/internal/usecases/item"
	orgUC "github.com/let-store-it/backend/internal/usecases/organization"
	storageUC "github.com/let-store-it/backend/internal/usecases/storage"
	taskUC "github.com/let-store-it/backend/internal/usecases/task"
)

type RestApiImplementation struct {
	orgUseCase          *orgUC.OrganizationUseCase
	orgUnitUseCase      *orgUC.OrganizationUseCase
	storageGroupUseCase *storageUC.StorageUseCase
	itemUseCase         *itemUC.ItemUseCase
	authUseCase         *authUC.AuthUseCase
	auditUseCase        *auditUC.AuditUseCase
	taskUseCase         *taskUC.TaskUseCase
}


// PatchCurrentUser implements api.Handler.
func (h *RestApiImplementation) PatchCurrentUser(ctx context.Context, req *api.PatchCurrentUserRequest) (api.PatchCurrentUserRes, error) {
	panic("unimplemented")
}

// PutCurrentUser implements api.Handler.
func (h *RestApiImplementation) PutCurrentUser(ctx context.Context, req *api.UpdateCurrentUserRequest) (api.PutCurrentUserRes, error) {
	panic("unimplemented")
}

func NewRestApiImplementation(
	orgUseCase *orgUC.OrganizationUseCase,
	orgUnitUseCase *orgUC.OrganizationUseCase,
	storageGroupUseCase *storageUC.StorageUseCase,
	itemUseCase *itemUC.ItemUseCase,
	authUseCase *authUC.AuthUseCase,
	auditUseCase *auditUC.AuditUseCase,
	taskUseCase *taskUC.TaskUseCase,
) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:          orgUseCase,
		orgUnitUseCase:      orgUnitUseCase,
		storageGroupUseCase: storageGroupUseCase,
		itemUseCase:         itemUseCase,
		authUseCase:         authUseCase,
		auditUseCase:        auditUseCase,
		taskUseCase:         taskUseCase,
	}
}
