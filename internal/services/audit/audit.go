package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

const (
	auditTopic = "audit.object-changes"
)

var (
	ErrInvalidObjectChange = errors.New("invalid object change")
	ErrInvalidOrganization = errors.New("invalid organization")
	ErrInvalidUser         = errors.New("invalid user")
	ErrInvalidTargetObject = errors.New("invalid target object")
)

type AuditService struct {
	queries *database.Queries
	kafka   *KafkaConfig
}

type AuditServiceConfig struct {
	Queries      *database.Queries
	KafkaEnabled bool
	KafkaBrokers []string
}

func New(cfg *AuditServiceConfig) (*AuditService, error) {
	if cfg == nil || cfg.Queries == nil {
		return nil, fmt.Errorf("invalid configuration")
	}

	service := &AuditService{
		queries: cfg.Queries,
	}

	if cfg.KafkaEnabled {
		kafka := NewKafkaConfig(cfg.KafkaBrokers)
		if err := kafka.Connect(context.Background(), auditTopic); err != nil {
			return nil, fmt.Errorf("failed to connect to kafka: %w", err)
		}
		service.kafka = kafka
	}

	return service, nil
}

func (s *AuditService) Close() error {
	if s.kafka != nil {
		if err := s.kafka.Close(); err != nil {
			return fmt.Errorf("failed to close kafka connection: %w", err)
		}
	}
	return nil
}

func (s *AuditService) validateObjectChange(objectChange *models.ObjectChange) error {
	if objectChange == nil {
		return ErrInvalidObjectChange
	}
	if objectChange.OrgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if objectChange.UserID == uuid.Nil {
		return ErrInvalidUser
	}
	if objectChange.TargetObjectID == uuid.Nil {
		return ErrInvalidTargetObject
	}
	return nil
}

func (s *AuditService) getObjectTypeInfo(ctx context.Context, typeID int32) (*models.ObjectType, error) {
	objectType, err := s.queries.GetObjectTypeById(ctx, typeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object type: %w", err)
	}

	return &models.ObjectType{
		ID:    models.ObjectTypeId(objectType.ID),
		Group: objectType.ObjectGroup,
		Name:  objectType.ObjectName,
	}, nil
}

func (s *AuditService) getEmployeeInfo(ctx context.Context, orgID, userID uuid.UUID) (*models.Employee, error) {
	employee, err := s.queries.GetEmployeeByUserId(ctx, database.GetEmployeeByUserIdParams{
		OrgID:  pgtype.UUID{Bytes: orgID, Valid: true},
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	role, err := s.queries.GetRoleById(ctx, employee.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	var middleName *string
	if employee.MiddleName.Valid {
		middleName = &employee.MiddleName.String
	}

	return &models.Employee{
		UserID:     employee.UserID.Bytes,
		Email:      employee.Email,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		MiddleName: middleName,
		RoleID:     int(employee.RoleID),
		Role: &models.Role{
			ID:          int(role.ID),
			Name:        role.Name,
			DisplayName: role.DisplayName,
			Description: role.Description,
		},
	}, nil
}

func (s *AuditService) publishToKafka(ctx context.Context, objectChange *models.ObjectChange) error {
	if s.kafka == nil {
		return nil
	}

	message, err := json.Marshal(objectChange)
	if err != nil {
		return fmt.Errorf("failed to marshal object change: %w", err)
	}

	key := []byte(fmt.Sprintf("%d", rand.Int()))
	if err := s.kafka.SendMessage(ctx, key, message); err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
	if err := s.validateObjectChange(objectChange); err != nil {
		return err
	}

	// Get related information
	objectType, err := s.getObjectTypeInfo(ctx, int32(objectChange.TargetObjectTypeId))
	if err != nil {
		return err
	}

	employee, err := s.getEmployeeInfo(ctx, objectChange.OrgID, objectChange.UserID)
	if err != nil {
		return err
	}

	// Create the object change record
	change, err := s.queries.CreateObjectChange(ctx, database.CreateObjectChangeParams{
		OrgID:            pgtype.UUID{Bytes: objectChange.OrgID, Valid: true},
		UserID:           pgtype.UUID{Bytes: objectChange.UserID, Valid: true},
		Action:           string(objectChange.Action),
		TargetObjectType: int32(objectChange.TargetObjectTypeId),
		TargetObjectID:   pgtype.UUID{Bytes: objectChange.TargetObjectID, Valid: true},
		PrechangeState:   objectChange.PrechangeState,
		PostchangeState:  objectChange.PostchangeState,
	})
	if err != nil {
		return fmt.Errorf("failed to create object change: %w", err)
	}

	// Update the object change with additional information
	objectChange.ID = change.ID.Bytes
	objectChange.ObjectType = objectType
	objectChange.Employee = employee

	return s.publishToKafka(ctx, objectChange)
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if targetObjectID == uuid.Nil {
		return nil, ErrInvalidTargetObject
	}

	// Get the object changes
	objectChanges, err := s.queries.GetObjectChanges(ctx, database.GetObjectChangesParams{
		OrgID:            pgtype.UUID{Bytes: orgID, Valid: true},
		TargetObjectType: int32(targetObjectTypeId),
		TargetObjectID:   pgtype.UUID{Bytes: targetObjectID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object changes: %w", err)
	}

	// Get the object type information
	objectType, err := s.getObjectTypeInfo(ctx, int32(targetObjectTypeId))
	if err != nil {
		return nil, err
	}

	objectChangesModels := make([]*models.ObjectChange, len(objectChanges))
	for i, change := range objectChanges {
		employee, err := s.getEmployeeInfo(ctx, change.OrgID.Bytes, change.UserID.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to get employee info for change %s: %w", change.ID.Bytes, err)
		}

		objectChangesModels[i] = &models.ObjectChange{
			ID:                 change.ID.Bytes,
			OrgID:              change.OrgID.Bytes,
			UserID:             change.UserID.Bytes,
			Action:             models.ObjectChangeAction(change.Action),
			TargetObjectTypeId: models.ObjectTypeId(objectType.ID),
			TargetObjectID:     change.TargetObjectID.Bytes,
			PrechangeState:     change.PrechangeState,
			PostchangeState:    change.PostchangeState,
			Timestamp:          change.Time.Time,
			ObjectType:         objectType,
			Employee:           employee,
		}
	}

	return objectChangesModels, nil
}
