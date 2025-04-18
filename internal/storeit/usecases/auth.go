package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/services"
	"github.com/let-store-it/backend/internal/storeit/services/yandex"
)

type AuthUseCase struct {
	authService        *services.AuthService
	yandexOAuthService *yandex.YandexOAuthService
}

func NewAuthUseCase(authService *services.AuthService, yandexOAuthService *yandex.YandexOAuthService) *AuthUseCase {
	return &AuthUseCase{authService: authService, yandexOAuthService: yandexOAuthService}
}

func (u *AuthUseCase) GetCurrentUser(ctx context.Context) (*models.User, error) {
	userID := ctx.Value(UserIDKey).(uuid.UUID)
	user, err := u.authService.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUseCase) CreateSessionByEmail(ctx context.Context, email string) (*models.Session, error) {
	user, err := u.authService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	session, err := u.authService.CreateUserSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *AuthUseCase) GetUserIdFromSession(ctx context.Context, sessionSecret string) (uuid.UUID, error) {
	user, err := u.authService.GetUserBySessionSecret(ctx, sessionSecret)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (u *AuthUseCase) ExchangeYandexAccessToken(ctx context.Context, accessToken string) (*models.Session, error) {
	userInfo, err := u.yandexOAuthService.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.authService.GetUserByEmail(ctx, userInfo.DefaultEmail)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			user = &models.User{
				Email:     userInfo.DefaultEmail,
				FirstName: userInfo.FirstName,
				LastName:  userInfo.LastName,
				YandexID:  &userInfo.ID,
			}
			user, err = u.authService.CreateUser(ctx, user)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	session, err := u.authService.CreateUserSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *AuthUseCase) validateOrganizationAccess(ctx context.Context) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	roles, err := uc.authService.GetUserRoles(ctx, userID, orgID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	if _, ok := roles[services.RoleOwner]; !ok {
		return uuid.Nil, fmt.Errorf("user is not an owner of the organization")
	}

	return orgID, nil
}
