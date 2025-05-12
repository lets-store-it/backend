package audit

import (
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toObjectChange(change sqlc.AppObjectChange) *models.ObjectChange {
	res := models.ObjectChange{
		ID:               change.ID.Bytes,
		OrgID:            change.OrgID.Bytes,
		UserID:           database.UUIDPtrFromPgx(change.UserID),
		Action:           models.ObjectChangeAction(change.Action),
		TargetObjectType: models.ObjectTypeId(change.TargetObjectType),
		TargetObjectID:   change.TargetObjectID.Bytes,
		PrechangeState:   change.PrechangeState,
		PostchangeState:  change.PostchangeState,
		Timestamp:        change.Time.Time,
	}
	return &res
}
