package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditFields struct {
	CreatedAt time.Time
	UpdatedAt *time.Time
	CreatedBy *uuid.UUID
	UpdatedBy *uuid.UUID
}

type Organization struct {
	ID          uuid.UUID
	Name        string
	Subdomain   string
	AuditFields *AuditFields
}
