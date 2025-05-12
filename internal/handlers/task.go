package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

func toAssignedToDTO(assignedTo *models.Employee) api.EmployeeOptional {
	var middleName api.NilString
	PtrToApiNil(assignedTo.MiddleName, &middleName)
	return api.EmployeeOptional{
		UserId:     assignedTo.UserID,
		FirstName:  assignedTo.FirstName,
		LastName:   assignedTo.LastName,
		MiddleName: middleName,
		Email:      assignedTo.Email,
		Role:       toRoleDTO(assignedTo.Role),
	}
}

func tasksToDto(tasks []*models.Task) []api.TaskBase {
	res := make([]api.TaskBase, len(tasks))
	for i, task := range tasks {
		res[i] = taskToDto(task)
	}
	return res
}

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

	var assignedTo api.NilEmployeeOptional
	if task.AssignedTo != nil {
		assignedTo.SetTo(toAssignedToDTO(task.AssignedTo))
	} else {
		assignedTo.SetToNull()
	}
	res.AssignedTo = assignedTo
	if task.Unit != nil {
		res.Unit = convertUnitToDTO(task.Unit)
	}
	return res
}

func taskItemToDto(item *models.TaskItem) api.TaskItem {
	var instance api.InstanceFull
	if item.Instance != nil {
		instance = convertItemInstanceToTaskItemDTO(item.Instance)
	}

	var sourceCell api.CellForInstance
	if item.SourceCell != nil {
		sourceCell = convertCellToDTO(item.SourceCell)
	}

	var targetCell api.NilCellForInstanceOptional
	if item.TargetCell != nil {
		targetCell.SetTo(api.CellForInstanceOptional{
			ID:       item.TargetCell.ID,
			Alias:    item.TargetCell.Alias,
			Row:      item.TargetCell.Row,
			Level:    item.TargetCell.Level,
			Position: item.TargetCell.Position,
			CellPath: convertCellPathToOptionalDTO(item.TargetCell.Path),
		})
	} else {
		targetCell.SetToNull()
	}

	return api.TaskItem{
		Instance:   instance,
		SourceCell: sourceCell,
		TargetCell: targetCell,
		Status:     api.TaskItemStatus(item.Status),
	}
}

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
	}

	if task.Unit != nil {
		res.Unit = convertUnitToDTO(task.Unit)
	}

	if task.AssignedTo != nil {
		var assignedTo api.NilEmployeeOptional
		assignedTo.SetTo(toAssignedToDTO(task.AssignedTo))
		res.AssignedTo = assignedTo
	} else {
		res.AssignedTo.SetToNull()
	}

	taskItems := make([]api.TaskItem, len(task.Items))
	for i, item := range task.Items {
		taskItems[i] = taskItemToDto(item)
	}
	res.Items = taskItems
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

func (h *RestApiImplementation) GetTaskById(ctx context.Context, params api.GetTaskByIdParams) (api.GetTaskByIdRes, error) {
	task, err := h.taskUseCase.GetTaskById(ctx, params.ID)
	if err != nil {
		return nil, err
	}
	return &api.GetTaskResponse{
		Data: taskToFullDto(task),
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

// PickInstanceFromCell implements api.Handler.
func (h *RestApiImplementation) PickInstanceFromCell(ctx context.Context, req *api.PickInstanceFromCellReq, params api.PickInstanceFromCellParams) (api.PickInstanceFromCellRes, error) {
	err := h.taskUseCase.PickInstanceFromCell(ctx, params.ID, req.InstanceId)
	if err != nil {
		return nil, err
	}
	return &api.PickInstanceFromCellNoContent{}, nil
}
