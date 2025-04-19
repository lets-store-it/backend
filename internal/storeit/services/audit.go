package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

type AuditService struct {
	queries *database.Queries
}

func NewAuditService(queries *database.Queries) *AuditService {
	return &AuditService{queries: queries}
}

func (s *AuditService) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
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
		return err
	}
	return nil
}

func (s *AuditService) GetObjectChanges(ctx context.Context, orgID uuid.UUID, targetObjectType models.ObjectType, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	objectChanges, err := s.queries.GetObjectChanges(ctx, database.GetObjectChangesParams{
		OrgID:            pgtype.UUID{Bytes: orgID, Valid: true},
		TargetObjectType: int32(targetObjectType),
		TargetObjectID:   pgtype.UUID{Bytes: targetObjectID, Valid: true},
	})
	if err != nil {
		return nil, err
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
