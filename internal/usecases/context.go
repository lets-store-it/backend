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
		return uuid.Nil, ErrOrganizationIDMissing
	}
	return orgID, nil
}

func GetTvBoardIDFromContext(ctx context.Context) (uuid.UUID, error) {
	tvBoardID, ok := ctx.Value(models.TvBoardIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, ErrTvBoardIDMissing
	}
	return tvBoardID, nil
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(models.UserIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context or invalid type")
	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

type AuthService interface {
	CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel models.AccessLevel) (bool, error)
}

type ValidateAccessResult struct {
	HasAccess  bool
	OrgID      uuid.UUID
	IsApiToken bool
	UserID     *uuid.UUID
}

func ValidateAccess(ctx context.Context, service AuthService, accessLevel models.AccessLevel) (ValidateAccessResult, error) {
	isSystemUser, ok := ctx.Value(models.IsSystemUserContextKey).(bool)
	if ok && isSystemUser {
		return ValidateAccessResult{
			HasAccess:  true,
			OrgID:      uuid.Nil,
			IsApiToken: true,
			UserID:     nil,
		}, nil
	}

	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, err
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, err
	}

	ok, err = service.CheckUserAccess(ctx, orgID, userID, accessLevel)
	if err != nil {
		return ValidateAccessResult{}, err
	}

	return ValidateAccessResult{
		HasAccess: ok,
		OrgID:     orgID,
		UserID:    &userID,
	}, nil
}
func ValidateAccessWithOptionalApiToken(ctx context.Context, service AuthService, accessLevel models.AccessLevel, allowApiToken bool) (ValidateAccessResult, error) {
	isSystemUser, ok := ctx.Value(models.IsSystemUserContextKey).(bool)
	if ok && isSystemUser {
		orgID, err := GetOrganizationIDFromContext(ctx)
		if err != nil {
			return ValidateAccessResult{}, err
		}

		return ValidateAccessResult{
			HasAccess:  true,
			OrgID:      orgID,
			IsApiToken: true,
			UserID:     nil,
		}, nil
	}

	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, err
	}
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, err
	}

	ok, err = service.CheckUserAccess(ctx, orgID, userID, accessLevel)
	if err != nil {
		return ValidateAccessResult{}, err
	}

	return ValidateAccessResult{
		HasAccess: ok,
		OrgID:     orgID,
		UserID:    &userID,
	}, nil
}
