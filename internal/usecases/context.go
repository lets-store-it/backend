package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
)

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
			return ValidateAccessResult{}, fmt.Errorf("%w: action can not be performed by api token", ErrForbidden)
		}
		orgID, err := common.GetOrganizationIDFromContext(ctx)
		if err != nil {
			return ValidateAccessResult{}, fmt.Errorf("%w: %v", common.ErrOrganizationIDMissing, err)
		}

		return ValidateAccessResult{
			IsAuthorized: true,
			OrgID:        orgID,
			IsApiToken:   true,
			UserID:       nil,
		}, nil
	}

	orgID, err := common.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, fmt.Errorf("%w: %v", common.ErrOrganizationIDMissing, err)
	}
	userID, err := common.GetUserIDFromContext(ctx)
	if err != nil {
		return ValidateAccessResult{}, fmt.Errorf("%w: %v", common.ErrUserIDMissing, err)
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
