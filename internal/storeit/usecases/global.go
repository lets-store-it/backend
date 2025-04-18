package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type contextKey string

const OrganizationIDKey contextKey = "organization_id"
const UserIDKey contextKey = "user_id"

func GetOrganizationIDFromContext(ctx context.Context) (uuid.UUID, error) {
	orgID := ctx.Value(OrganizationIDKey).(uuid.UUID)
	if orgID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("organization ID not found in context")
	}
	return orgID, nil
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userID := ctx.Value(UserIDKey).(uuid.UUID)
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
