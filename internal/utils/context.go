package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
)

func GetOrganizationIDFromContext(ctx context.Context) (uuid.UUID, error) {
	orgID, ok := ctx.Value(models.OrganizationIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("organization ID not found in context or invalid type")
	}
	if orgID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("organization ID not found in context")
	}
	return orgID, nil
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(models.UserIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context or invalid type")
	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
