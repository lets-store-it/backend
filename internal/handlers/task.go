package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

func taskToDto(task *models.Task) api.TaskBase {
	var description api.NilString
	PtrToApiNil(task.Description, &description)

	var assignedAt api.NilDateTime
	PtrToApiNil(task.AssignedAt, &assignedAt)

	var completedAt api.NilDateTime
	PtrToApiNil(task.CompletedAt, &completedAt)

	res := api.TaskBase{
		ID:          task.ID,
		Name:        task.Name,
		Description: description,
		Status:      api.TaskBaseStatus(task.Status),
		CreatedAt:   task.CreatedAt,
		AssignedAt:  assignedAt,
		CompletedAt: completedAt,
		Type:        api.TaskBaseType(task.Type),
	}
	if task.AssignedTo != nil {
		res.AssignedTo = toEmployeeDTO(task.AssignedTo)
	}
	if task.Unit != nil {
		res.Unit = convertUnitToDTO(task.Unit)
	}
	return res
}

// func taskItemToDto(item *models.TaskItem) api.TaskItem {
// 	return api.TaskItem{
// 		Instance:   convertItemInstanceToTaskItemDTO(item.Instance),
// 		SourceCell: convertCellToDTO(item.SourceCell),
// 		TargetCell: convertCellToDTO(item.TargetCell),
// 	}
// }

func taskToFullDto(task *models.Task) api.TaskFull {
	var description api.NilString
	PtrToApiNil(task.Description, &description)

	var assignedAt api.NilDateTime
	PtrToApiNil(task.AssignedAt, &assignedAt)

	var completedAt api.NilDateTime
	PtrToApiNil(task.CompletedAt, &completedAt)

	res := api.TaskFull{
		ID:          task.ID,
		Name:        task.Name,
		Description: description,
		Status:      api.TaskFullStatus(task.Status),
		CreatedAt:   task.CreatedAt,
		AssignedAt:  assignedAt,
		CompletedAt: completedAt,
		Type:        api.TaskFullType(task.Type),
		Unit:        convertUnitToDTO(task.Unit),
		AssignedTo:  toEmployeeDTO(task.AssignedTo),
	}
	// for _, item := range task.Items {
	// 	res.Items = append(res.Items, taskItemToDto(item))
	// }
	return res
}

func (h *RestApiImplementation) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (api.CreateTaskRes, error) {
	task := &models.Task{
		Name:             req.Name,
		Description:      ApiValueToPtr(req.Description),
		Status:           models.TaskStatusPending,
		Type:             models.TaskType(req.Type),
		UnitID:           req.UnitId,
		AssignedToUserID: ApiValueToPtr(req.AssignedTo),
	}

	items := make([]*models.TaskItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = &models.TaskItem{
			InstanceID:   item.InstanceId,
			TargetCellID: ApiValueToPtr(item.TargetCellId),
		}
	}
	task.Items = items

	createdTask, err := h.taskUseCase.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return &api.CreateTaskResponse{
		Data: taskToFullDto(createdTask),
	}, nil
}

func (h *RestApiImplementation) GetTasks(ctx context.Context) (api.GetTasksRes, error) {
	tasks, err := h.taskUseCase.GetTasks(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]api.TaskBase, len(tasks))
	for i, task := range tasks {
		dtos[i] = taskToDto(task)
	}
	return &api.GetTasksResponse{
		Data: dtos,
	}, nil
}
