package models

import (
	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypePickItem TaskType = "pick"
	TaskTypeMovement TaskType = "movement"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type TaskItem struct {
	ID           uuid.UUID `json:"id"`
	InstanceID   uuid.UUID `json:"instance_id"`
	TargetCellID uuid.UUID `json:"target_cell_id"`

	Instance   *ItemInstance `json:"instance"`
	TargetCell *Cell         `json:"target_cell"`
}

type Task struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	UnitID      uuid.UUID  `json:"unit_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      TaskStatus `json:"status"`
	Type        TaskType   `json:"type"`
	Items       []TaskItem `json:"items"`
}
