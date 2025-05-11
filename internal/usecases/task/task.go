package task

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/tasks"
	"github.com/let-store-it/backend/internal/usecases"
)

type TaskUseCase struct {
	taskService *tasks.TaskService
	authService *auth.AuthService
	orgService  *organization.OrganizationService
}

type TaskUseCaseConfig struct {
	TaskService *tasks.TaskService
	AuthService *auth.AuthService
	OrgService  *organization.OrganizationService
}

func New(config TaskUseCaseConfig) *TaskUseCase {
	return &TaskUseCase{
		taskService: config.TaskService,
		authService: config.AuthService,
		orgService:  config.OrgService,
	}
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, services.ErrNotAuthorized
	}

	task.OrgID = validateResult.OrgID

	createdTask, err := uc.taskService.CreateTask(ctx, validateResult.OrgID, task)
	if err != nil {
		return nil, err
	}

	createdTask.Unit, err = uc.orgService.GetUnitByID(ctx, validateResult.OrgID, createdTask.UnitID)
	if err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (uc *TaskUseCase) GetTaskById(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, services.ErrNotAuthorized
	}

	task, err := uc.taskService.GetTaskById(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (uc *TaskUseCase) GetTasks(ctx context.Context) ([]*models.Task, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, services.ErrNotAuthorized
	}

	tasks, err := uc.taskService.GetTasks(ctx, validateResult.OrgID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (uc *TaskUseCase) PickInstanceFromCell(ctx context.Context, taskID uuid.UUID, instanceID uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return services.ErrNotAuthorized
	}

	return uc.taskService.PickInstance(ctx, validateResult.OrgID, taskID, instanceID)
}
