package models

import (
	"github.com/google/uuid"
)

type RoleName string

const (
	RoleOwner   RoleName = "org_owner"
	RoleAdmin   RoleName = "org_admin"
	RoleManager RoleName = "org_manager"
	RoleWorker  RoleName = "org_worker"
)

type RoleID int

const (
	RoleOwnerID   RoleID = 1
	RoleAdminID   RoleID = 2
	RoleManagerID RoleID = 3
	RoleWorkerID  RoleID = 4
)

type Role struct {
	ID          int
	Name        RoleName
	DisplayName string
	Description string
}

type Employee struct {
	UserID     uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName *string   `json:"middle_name"`
	RoleID     int       `json:"role_id"`

	Role *Role `json:"role"`
}
