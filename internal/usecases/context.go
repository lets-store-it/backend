package usecases

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

type AuthService interface {
	CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel models.AccessLevel) (bool, error)
}

type ValidateAccessResult struct {
	IsAuthorized bool
	OrgID        uuid.UUID
	IsApiToken   bool
	UserID       *uuid.UUID
}

func ValidateAccessWithOptionalApiToken(ctx context.Context, service AuthService, accessLevel models.AccessLevel, allowApiToken bool) (ValidateAccessResult, error) {
	isSystemUser, ok := ctx.Value(models.IsSystemUserContextKey).(bool)
	if ok && isSystemUser {
		if !allowApiToken {
			return ValidateAccessResult{}, fmt.Errorf("action can not be performed by api token")
		}
		orgID, err := GetOrganizationIDFromContext(ctx)
		if err != nil {
			return ValidateAccessResult{}, fmt.Errorf("%w: %v", ErrOrganizationIDMissing, err)
		}

		return ValidateAccessResult{
			IsAuthorized: true,
			OrgID:        orgID,
			IsApiToken:   true,
			UserID:       nil,
		}, nil
	}

	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, fmt.Errorf("%w: %v", ErrOrganizationIDMissing, err)
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, fmt.Errorf("%w: %v", ErrUserIDMissing, err)
	}

	ok, err = service.CheckUserAccess(ctx, orgID, userID, accessLevel)
	if err != nil {
		return ValidateAccessResult{}, fmt.Errorf("%w: %v", ErrNotAuthorized, err)
	}

	return ValidateAccessResult{
		IsAuthorized: ok,
		OrgID:        orgID,
		UserID:       &userID,
	}, nil
}
