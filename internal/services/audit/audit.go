package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
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
	pgxPool *pgxpool.Pool
	queries *database.Queries
	auth    *auth.AuthService
	kafka   *KafkaConfig
	logger  *slog.Logger
}

type AuditServiceConfig struct {
	PGXPool      *pgxpool.Pool
	Queries      *database.Queries
	Auth         *auth.AuthService
	KafkaEnabled bool
	KafkaBrokers []string
	Logger       *slog.Logger
}

func New(cfg *AuditServiceConfig) (*AuditService, error) {
	if cfg == nil || cfg.Queries == nil || cfg.PGXPool == nil {
		return nil, fmt.Errorf("invalid configuration")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	// Add service name prefix to all log messages
	logger = logger.With("service", "audit")

	service := &AuditService{
		pgxPool: cfg.PGXPool,
		queries: cfg.Queries,
		auth:    cfg.Auth,
		logger:  logger,
	}

	if cfg.KafkaEnabled {
		kafka := NewKafkaConfig(cfg.KafkaBrokers)
		if err := kafka.Connect(context.Background(), auditTopic); err != nil {
			return nil, fmt.Errorf("failed to connect to kafka: %w", err)
		}
		service.kafka = kafka
		service.logger.Info("kafka connection established", "brokers", cfg.KafkaBrokers)
	}

	return service, nil
}

func (s *AuditService) Close() error {
	if s.kafka != nil {
		if err := s.kafka.Close(); err != nil {
			s.logger.Error("failed to close kafka connection",
				"method", "Close",
				"error", err)
			return fmt.Errorf("failed to close kafka connection: %w", err)
		}
		s.logger.Info("kafka connection closed",
			"method", "Close")
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

func (s *AuditService) publishToKafka(ctx context.Context, objectChange *models.ObjectChange) error {
	if s.kafka == nil {
		return nil
	}

	message, err := json.Marshal(objectChange)
	if err != nil {
		s.logger.Error("failed to marshal object change",
			"method", "publishToKafka",
			"error", err)
		return fmt.Errorf("failed to marshal object change: %w", err)
	}

	key := []byte(fmt.Sprintf("%d", rand.Int()))
	if err := s.kafka.SendMessage(ctx, key, message); err != nil {
		s.logger.Error("failed to send message to kafka",
			"method", "publishToKafka",
			"error", err,
			"object_change_id", objectChange.ID,
			"org_id", objectChange.OrgID,
			"target_object_id", objectChange.TargetObjectID)
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	s.logger.Info("kafka message sent successfully",
		"method", "publishToKafka",
		"object_change_id", objectChange.ID,
		"org_id", objectChange.OrgID,
		"target_object_id", objectChange.TargetObjectID)

	return nil
}

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.logger.Error("failed to begin transaction",
			"method", "CreateObjectChange",
			"error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.validateObjectChange(objectChange); err != nil {
		s.logger.Error("invalid object change",
			"method", "CreateObjectChange",
			"error", err,
			"org_id", objectChange.OrgID,
			"user_id", objectChange.UserID,
			"target_object_id", objectChange.TargetObjectID)
		return err
	}

	// Get related information
	objectType, err := s.getObjectTypeInfo(ctx, int32(objectChange.TargetObjectTypeId))
	if err != nil {
		s.logger.Error("failed to get object type info",
			"method", "CreateObjectChange",
			"error", err,
			"type_id", objectChange.TargetObjectTypeId)
		return err
	}

	employee, err := s.auth.GetEmployeeWithRole(ctx, objectChange.OrgID, objectChange.UserID)
	if err != nil {
		s.logger.Error("failed to get employee with role",
			"method", "CreateObjectChange",
			"error", err,
			"org_id", objectChange.OrgID,
			"user_id", objectChange.UserID)
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
		s.logger.Error("failed to create object change",
			"method", "CreateObjectChange",
			"error", err,
			"org_id", objectChange.OrgID,
			"user_id", objectChange.UserID,
			"action", objectChange.Action,
			"target_object_id", objectChange.TargetObjectID)
		return fmt.Errorf("failed to create object change: %w", err)
	}

	// Update the object change with additional information
	objectChange.ID = change.ID.Bytes
	objectChange.ObjectType = objectType
	objectChange.Employee = employee

	s.logger.Info("object change created successfully",
		"method", "CreateObjectChange",
		"change_id", objectChange.ID,
		"org_id", objectChange.OrgID,
		"user_id", objectChange.UserID,
		"action", objectChange.Action,
		"target_object_id", objectChange.TargetObjectID)

	return s.publishToKafka(ctx, objectChange)
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	if orgID == uuid.Nil {
		s.logger.Error("invalid organization ID",
			"method", "GetObjectChanges",
			"org_id", orgID)
		return nil, ErrInvalidOrganization
	}
	if targetObjectID == uuid.Nil {
		s.logger.Error("invalid target object ID",
			"method", "GetObjectChanges",
			"target_object_id", targetObjectID)
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
		employee, err := s.auth.GetEmployeeWithRole(ctx, change.OrgID.Bytes, change.UserID.Bytes)
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
