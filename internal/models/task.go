package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypePickmentItem TaskType = "pickment"
	TaskTypeMovement     TaskType = "movement"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReady      TaskStatus = "ready"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusCancelled  TaskStatus = "cancelled"
)

type TaskItemStatus string

const (
	TaskItemStatusPending  TaskItemStatus = "pending"
	TaskItemStatusPicked   TaskItemStatus = "picked"
	TaskItemStatusDone     TaskItemStatus = "done"
	TaskItemStatusReturned TaskItemStatus = "returned"
)

type TaskItem struct {
	OrgID      uuid.UUID `json:"org_id"`
	TaskID     uuid.UUID `json:"task_id"`
	InstanceID uuid.UUID `json:"instance_id"`

	SourceCellID *uuid.UUID `json:"source_cell_id"`
	TargetCellID *uuid.UUID `json:"target_cell_id"`

	Status TaskItemStatus `json:"status"`

	Instance   *ItemInstance `json:"instance"`
	SourceCell *Cell         `json:"source_cell"`
	TargetCell *Cell         `json:"target_cell"`
}

type Task struct {
	ID uuid.UUID `json:"id"`

	OrgID  uuid.UUID `json:"org_id"`
	UnitID uuid.UUID `json:"unit_id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`

	Status           TaskStatus `json:"status"`
	AssignedToUserID *uuid.UUID `json:"assigned_to_user_id"`
	Type             TaskType   `json:"type"`

	Items      []*TaskItem       `json:"items"`
	AssignedTo *Employee         `json:"assigned_to"`
	Unit       *OrganizationUnit `json:"unit"`

	CreatedAt   time.Time  `json:"created_at"`
	AssignedAt  *time.Time `json:"assigned_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
