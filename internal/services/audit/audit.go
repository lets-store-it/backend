package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/telemetry"
	"github.com/let-store-it/backend/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type AuditService struct {
	pgxPool *pgxpool.Pool
	queries *sqlc.Queries

	auth *auth.AuthService

	kafka      *KafkaConfig
	kafkaTopic string
	tracer     trace.Tracer
}

type AuditServiceConfig struct {
	PGXPool *pgxpool.Pool
	Queries *sqlc.Queries

	Auth *auth.AuthService

	KafkaEnabled bool
	KafkaBrokers []string
	KafkaTopic   string
}

func New(cfg AuditServiceConfig) (*AuditService, error) {
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
		return fmt.Errorf("%w: object change is nil", services.ErrValidationError)
	}
	if objectChange.OrgID == uuid.Nil {
		return fmt.Errorf("%w: organization ID is nil", services.ErrValidationError)
	}
	if objectChange.TargetObjectID == uuid.Nil {
		return fmt.Errorf("%w: target object ID is nil", services.ErrValidationError)
	}
	return nil
}

func (s *AuditService) getObjectTypeInfo(ctx context.Context, typeID int32) (*models.ObjectType, error) {
	return telemetry.WithTrace(ctx, s.tracer, "getObjectTypeInfo", func(ctx context.Context, span trace.Span) (*models.ObjectType, error) {
		objectType, err := s.queries.GetObjectTypeById(ctx, typeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get object type: %w", err)
		}

		result := &models.ObjectType{
			ID:    models.ObjectTypeId(objectType.ID),
			Group: objectType.ObjectGroup,
			Name:  objectType.ObjectName,
		}

		span.SetAttributes(
			attribute.String("object.group", result.Group),
			attribute.String("object.name", result.Name),
		)
		span.SetStatus(codes.Ok, "successfully retrieved object type")

		return result, nil
	})
}

func (s *AuditService) publishToKafka(ctx context.Context, objectChange *models.ObjectChange) error {
	if s.kafka == nil {
		return nil
	}

	return telemetry.WithVoidTrace(ctx, s.tracer, "publishToKafka", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("change.id", objectChange.ID.String()),
			attribute.String("org.id", objectChange.OrgID.String()),
			attribute.String("target.object.id", objectChange.TargetObjectID.String()),
			attribute.String("kafka.topic", s.kafkaTopic),
		)

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
	})
}

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
	return telemetry.WithVoidTrace(ctx, s.tracer, "CreateObjectChange", func(ctx context.Context, span trace.Span) error {
		span.SetAttributes(
			attribute.String("org.id", objectChange.OrgID.String()),
			attribute.String("user.id", utils.SafeUUIDString(objectChange.UserID)),
			attribute.String("action", string(objectChange.Action)),
			attribute.String("target.object.id", objectChange.TargetObjectID.String()),
		)

		if err := s.validateObjectChange(objectChange); err != nil {
			return err
		}

		employee, err := s.auth.GetUserAsEmployeeInOrg(ctx, objectChange.OrgID, *objectChange.UserID)
		if err != nil {
			return err
		}

		err = database.WithVoidTransaction(ctx, s.pgxPool, s.tracer, func(ctx context.Context, tx pgx.Tx) error {
			qtx := s.queries.WithTx(tx)
			change, err := qtx.CreateObjectChange(ctx, sqlc.CreateObjectChangeParams{
				OrgID:            database.PgUUID(objectChange.OrgID),
				UserID:           database.PgUUIDPtr(objectChange.UserID),
				Action:           string(objectChange.Action),
				TargetObjectType: int32(objectChange.TargetObjectTypeId),
				TargetObjectID:   database.PgUUID(objectChange.TargetObjectID),
				PrechangeState:   objectChange.PrechangeState,
				PostchangeState:  objectChange.PostchangeState,
			})
			if err != nil {
				return fmt.Errorf("failed to create object change: %w", err)
			}

			objectChange.ID = change.ID.Bytes
			objectChange.Employee = employee

			if err := s.publishToKafka(ctx, objectChange); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}

		span.SetAttributes(
			attribute.String("change.id", objectChange.ID.String()),
		)

		return nil
	})
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	return telemetry.WithTrace(ctx, s.tracer, "GetObjectChanges", func(ctx context.Context, span trace.Span) ([]*models.ObjectChange, error) {
		span.SetAttributes(
			attribute.String("org.id", orgID.String()),
			attribute.String("target.object.type.id", fmt.Sprintf("%d", targetObjectTypeId)),
			attribute.String("target.object.id", targetObjectID.String()),
		)

		objectChanges, err := s.queries.GetObjectChanges(ctx, sqlc.GetObjectChangesParams{
			OrgID:            database.PgUUID(orgID),
			TargetObjectType: int32(targetObjectTypeId),
			TargetObjectID:   database.PgUUID(targetObjectID),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get object changes: %w", err)
		}

		objectType, err := s.getObjectTypeInfo(ctx, int32(targetObjectTypeId))
		if err != nil {
			return nil, err
		}

		objectChangesModels := make([]*models.ObjectChange, len(objectChanges))
		for i, change := range objectChanges {
			employee, err := s.auth.GetUserAsEmployeeInOrg(ctx, change.OrgID.Bytes, change.UserID.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to get employee info for change %s: %w", change.ID.Bytes, err)
			}

			objectChangesModels[i] = &models.ObjectChange{
				ID:                 change.ID.Bytes,
				OrgID:              change.OrgID.Bytes,
				UserID:             database.UUIDPtrFromPgx(change.UserID),
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
			attribute.Int("changes.count", len(objectChangesModels)),
			attribute.String("object.type", objectType.Name),
		)
		return objectChangesModels, nil
	})
}
