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

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
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

	_, err := s.queries.CreateObjectChange(ctx, database.CreateObjectChangeParams{
		OrgID:            pgtype.UUID{Bytes: objectChange.OrgID, Valid: true},
		UserID:           pgtype.UUID{Bytes: objectChange.UserID, Valid: true},
		Action:           string(objectChange.Action),
		TargetObjectType: int32(objectChange.TargetObjectType),
		TargetObjectID:   pgtype.UUID{Bytes: objectChange.TargetObjectID, Valid: true},
		PrechangeState:   objectChange.PrechangeState,
		PostchangeState:  objectChange.PostchangeState,
	})
	if err != nil {
		return fmt.Errorf("failed to create object change: %w", err)
	}

	// Send to Kafka if enabled
	if s.kafka != nil {
		message, err := json.Marshal(objectChange)
		if err != nil {
			return fmt.Errorf("failed to marshal object change: %w", err)
		}

		// Use random number as key for even distribution
		key := []byte(fmt.Sprintf("%d", rand.Int()))
		if err := s.kafka.SendMessage(ctx, key, message); err != nil {
			return fmt.Errorf("failed to send message to kafka: %w", err)
		}
	}

	return nil
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectType models.ObjectType, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if targetObjectID == uuid.Nil {
		return nil, ErrInvalidTargetObject
	}

	objectChanges, err := s.queries.GetObjectChanges(ctx, database.GetObjectChangesParams{
		OrgID:            pgtype.UUID{Bytes: orgID, Valid: true},
		TargetObjectType: int32(targetObjectType),
		TargetObjectID:   pgtype.UUID{Bytes: targetObjectID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object changes: %w", err)
	}

	objectChangesModels := make([]*models.ObjectChange, len(objectChanges))
	for i, objectChange := range objectChanges {
		objectChangesModels[i] = &models.ObjectChange{
			ID:               objectChange.ID.Bytes,
			OrgID:            objectChange.OrgID.Bytes,
			UserID:           objectChange.UserID.Bytes,
			Action:           models.ObjectChangeAction(objectChange.Action),
			TargetObjectType: models.ObjectType(objectChange.TargetObjectType),
			TargetObjectID:   objectChange.TargetObjectID.Bytes,
			PrechangeState:   objectChange.PrechangeState,
			PostchangeState:  objectChange.PostchangeState,
		}
	}
	return objectChangesModels, nil
}
