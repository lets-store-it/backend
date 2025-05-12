package common

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
)

func GetOrganizationIDFromContext(ctx context.Context) (uuid.UUID, error) {
	orgID, ok := ctx.Value(models.OrganizationIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrOrganizationIDMissing
	}
	if orgID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: organization ID is null", ErrOrganizationIDMissing)
	}
	return orgID, nil
}

func GetTvBoardIDFromContext(ctx context.Context) (uuid.UUID, error) {
	tvBoardID, ok := ctx.Value(models.TvBoardIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrTvBoardIDMissing
	}
	if tvBoardID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: TV board ID is null", ErrTvBoardIDMissing)
	}
	return tvBoardID, nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(models.UserIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrUserIDMissing
	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: user ID is null", ErrUserIDMissing)
	}
	return userID, nil
}

func GetUserIDFromContextIfExists(ctx context.Context) (*uuid.UUID, error) {
	userID, ok := ctx.Value(models.UserIDContextKey).(uuid.UUID)
	if !ok {
		return nil, nil
	}
	return &userID, nil
}
