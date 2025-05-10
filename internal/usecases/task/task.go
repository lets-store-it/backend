package task

import (
	"context"

	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/tasks"
	"github.com/let-store-it/backend/internal/usecases"
)

type TaskUseCase struct {
	taskService *tasks.TaskService
	authService *auth.AuthService
}

type TaskUseCaseConfig struct {
	TaskService *tasks.TaskService
	AuthService *auth.AuthService
}

func New(config TaskUseCaseConfig) *TaskUseCase {
	return &TaskUseCase{
		taskService: config.TaskService,
		authService: config.AuthService,
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

	return uc.taskService.CreateTask(ctx, validateResult.OrgID, task)
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
