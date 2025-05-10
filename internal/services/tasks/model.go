package tasks

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toTask(task sqlc.Task) *models.Task {
	model := &models.Task{
		ID:          database.UUIDFromPgx(task.ID),
		OrgID:       database.UUIDFromPgx(task.OrgID),
		UnitID:      database.UUIDFromPgx(task.UnitID),
		Name:        task.Name,
		Description: database.PgTextPtrFromPgx(task.Description),

		Status:           models.TaskStatus(task.Status),
		AssignedToUserID: database.UUIDPtrFromPgx(task.AssignedToUserID),
		Type:             models.TaskType(task.Type),

		// Items: []models.TaskItem{},

		CreatedAt:   task.CreatedAt.Time,
		AssignedAt:  database.PgTimePtrFromPgx(task.AssignedAt),
		CompletedAt: database.PgTimePtrFromPgx(task.CompletedAt),

		AssignedTo: nil,
	}

	return model
}

func toTaskItem(taskItem sqlc.TaskItem) *models.TaskItem {
	return &models.TaskItem{
		OrgID:        database.UUIDFromPgx(taskItem.OrgID),
		TaskID:       database.UUIDFromPgx(taskItem.TaskID),
		InstanceID:   database.UUIDFromPgx(taskItem.ItemInstanceID),
		SourceCellID: database.UUIDPtrFromPgx(taskItem.SourceCellID),
		TargetCellID: database.UUIDPtrFromPgx(taskItem.DestinationCellID),
		Status:       models.TaskItemStatus(taskItem.Status),
	}
}
