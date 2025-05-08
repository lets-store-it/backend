package audit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	database "github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	auth *auth.AuthService

	kafka      *KafkaConfig
	kafkaTopic string
	tracer     trace.Tracer
}

type AuditServiceConfig struct {
	PGXPool *pgxpool.Pool
	Queries *database.Queries

	Auth *auth.AuthService

	KafkaEnabled bool
	KafkaBrokers []string
	KafkaTopic   string
}

func New(cfg AuditServiceConfig) (*AuditService, error) {
	if cfg.Queries == nil || cfg.PGXPool == nil {
		return nil, fmt.Errorf("invalid configuration")
	}

	service := &AuditService{
		pgxPool:    cfg.PGXPool,
		queries:    cfg.Queries,
		auth:       cfg.Auth,
		tracer:     otel.GetTracerProvider().Tracer("audit-service"),
		kafkaTopic: cfg.KafkaTopic,
	}

	if cfg.KafkaEnabled {
		if cfg.KafkaTopic == "" {
			return nil, fmt.Errorf("kafka topic is required when kafka is enabled")
		}
		kafka := NewKafkaConfig(cfg.KafkaBrokers)
		if err := kafka.Connect(context.Background(), cfg.KafkaTopic); err != nil {
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
	ctx, span := s.tracer.Start(ctx, "getObjectTypeInfo",
		trace.WithAttributes(
			attribute.Int("type_id", int(typeID)),
		),
	)
	defer span.End()

	objectType, err := s.queries.GetObjectTypeById(ctx, typeID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get object type")
		return nil, fmt.Errorf("failed to get object type: %w", err)
	}

	result := &models.ObjectType{
		ID:    models.ObjectTypeId(objectType.ID),
		Group: objectType.ObjectGroup,
		Name:  objectType.ObjectName,
	}

	span.SetAttributes(
		attribute.String("object_group", result.Group),
		attribute.String("object_name", result.Name),
	)
	span.SetStatus(codes.Ok, "successfully retrieved object type")

	return result, nil
}

func (s *AuditService) publishToKafka(ctx context.Context, objectChange *models.ObjectChange) error {
	if s.kafka == nil {
		return nil
	}

	ctx, span := s.tracer.Start(ctx, "publishToKafka",
		trace.WithAttributes(
			attribute.String("change_id", objectChange.ID.String()),
			attribute.String("org_id", objectChange.OrgID.String()),
			attribute.String("target_object_id", objectChange.TargetObjectID.String()),
			attribute.String("kafka_topic", s.kafkaTopic),
		),
	)
	defer span.End()

	message, err := json.Marshal(objectChange)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to marshal object change")
		return fmt.Errorf("failed to marshal object change: %w", err)
	}

	key := []byte(fmt.Sprintf("%d", rand.Int()))
	if err := s.kafka.SendMessage(ctx, key, message); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to send message to kafka")
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	span.SetStatus(codes.Ok, "successfully published to kafka")
	return nil
}

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
	ctx, span := s.tracer.Start(ctx, "CreateObjectChange",
		trace.WithAttributes(
			attribute.String("org_id", objectChange.OrgID.String()),
			attribute.String("user_id", objectChange.UserID.String()),
			attribute.String("action", string(objectChange.Action)),
			attribute.String("target_object_id", objectChange.TargetObjectID.String()),
		),
	)
	defer span.End()

	tx, err := s.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to begin transaction")
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.validateObjectChange(objectChange); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid object change")
		return err
	}

	// Get related information
	objectType, err := s.getObjectTypeInfo(ctx, int32(objectChange.TargetObjectTypeId))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get object type info")
		return err
	}

	employee, err := s.auth.GetEmployeeWithRole(ctx, objectChange.OrgID, objectChange.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get employee with role")
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
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create object change")
		return fmt.Errorf("failed to create object change: %w", err)
	}

	// Update the object change with additional information
	objectChange.ID = change.ID.Bytes
	objectChange.ObjectType = objectType
	objectChange.Employee = employee

	span.SetAttributes(
		attribute.String("change_id", objectChange.ID.String()),
		attribute.String("object_type", objectType.Name),
	)

	if err := s.publishToKafka(ctx, objectChange); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to publish to kafka")
		return err
	}

	span.SetStatus(codes.Ok, "object change created successfully")
	return nil
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	ctx, span := s.tracer.Start(ctx, "GetObjectChanges",
		trace.WithAttributes(
			attribute.String("org_id", orgID.String()),
			attribute.String("target_object_type_id", fmt.Sprintf("%d", targetObjectTypeId)),
			attribute.String("target_object_id", targetObjectID.String()),
		),
	)
	defer span.End()

	if orgID == uuid.Nil {
		span.RecordError(ErrInvalidOrganization)
		span.SetStatus(codes.Error, "invalid organization ID")
		return nil, ErrInvalidOrganization
	}
	if targetObjectID == uuid.Nil {
		span.RecordError(ErrInvalidTargetObject)
		span.SetStatus(codes.Error, "invalid target object ID")
		return nil, ErrInvalidTargetObject
	}

	// Get the object changes
	objectChanges, err := s.queries.GetObjectChanges(ctx, database.GetObjectChangesParams{
		OrgID:            pgtype.UUID{Bytes: orgID, Valid: true},
		TargetObjectType: int32(targetObjectTypeId),
		TargetObjectID:   pgtype.UUID{Bytes: targetObjectID, Valid: true},
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get object changes")
		return nil, fmt.Errorf("failed to get object changes: %w", err)
	}

	// Get the object type information
	objectType, err := s.getObjectTypeInfo(ctx, int32(targetObjectTypeId))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get object type info")
		return nil, err
	}

	objectChangesModels := make([]*models.ObjectChange, len(objectChanges))
	for i, change := range objectChanges {
		employee, err := s.auth.GetEmployeeWithRole(ctx, change.OrgID.Bytes, change.UserID.Bytes)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "failed to get employee info")
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

	span.SetAttributes(
		attribute.Int("changes_count", len(objectChangesModels)),
		attribute.String("object_type", objectType.Name),
	)
	span.SetStatus(codes.Ok, "successfully retrieved object changes")

	return objectChangesModels, nil
}
