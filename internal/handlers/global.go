package handlers

import (
	auditUC "github.com/let-store-it/backend/internal/usecases/audit"
	authUC "github.com/let-store-it/backend/internal/usecases/auth"
	itemUC "github.com/let-store-it/backend/internal/usecases/item"
	orgUC "github.com/let-store-it/backend/internal/usecases/organization"
	storageUC "github.com/let-store-it/backend/internal/usecases/storage"
	taskUC "github.com/let-store-it/backend/internal/usecases/task"
	tvboardUC "github.com/let-store-it/backend/internal/usecases/tv_board"
)

type RestApiImplementation struct {
	orgUseCase          *orgUC.OrganizationUseCase
	orgUnitUseCase      *orgUC.OrganizationUseCase
	storageGroupUseCase *storageUC.StorageUseCase
	itemUseCase         *itemUC.ItemUseCase
	authUseCase         *authUC.AuthUseCase
	auditUseCase        *auditUC.AuditUseCase
	taskUseCase         *taskUC.TaskUseCase
	tvBoardUseCase      *tvboardUC.TvBoardUseCase
}

func NewRestApiImplementation(
	orgUseCase *orgUC.OrganizationUseCase,
	orgUnitUseCase *orgUC.OrganizationUseCase,
	storageGroupUseCase *storageUC.StorageUseCase,
	itemUseCase *itemUC.ItemUseCase,
	authUseCase *authUC.AuthUseCase,
	auditUseCase *auditUC.AuditUseCase,
	taskUseCase *taskUC.TaskUseCase,
	tvBoardUseCase *tvboardUC.TvBoardUseCase,
) *RestApiImplementation {
	return &RestApiImplementation{
		orgUseCase:          orgUseCase,
		orgUnitUseCase:      orgUnitUseCase,
		storageGroupUseCase: storageGroupUseCase,
		itemUseCase:         itemUseCase,
		authUseCase:         authUseCase,
		auditUseCase:        auditUseCase,
		taskUseCase:         taskUseCase,
		tvBoardUseCase:      tvBoardUseCase,
	}
}
