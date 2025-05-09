package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/yandex"
	"github.com/let-store-it/backend/internal/usecases"
)

type AuthUseCase struct {
	authService        *auth.AuthService
	yandexOAuthService *yandex.YandexOAuthService
}

type AuthUseCaseConfig struct {
	AuditService       *audit.AuditService
	AuthService        *auth.AuthService
	YandexOAuthService *yandex.YandexOAuthService
}

func New(config AuthUseCaseConfig) *AuthUseCase {
	return &AuthUseCase{
		authService:        config.AuthService,
		yandexOAuthService: config.YandexOAuthService,
	}
}

func (u *AuthUseCase) GetCurrentUser(ctx context.Context) (*models.User, error) {
	userID, err := usecases.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUseCase) CreateSessionByEmail(ctx context.Context, email string) (*models.UserSession, error) {
	user, err := u.authService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	session, err := u.authService.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *AuthUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (*models.UserSession, error) {
	return u.authService.CreateSession(ctx, userID)
}

func (u *AuthUseCase) InvalidateSession(ctx context.Context, sessionID uuid.UUID) error {
	return u.authService.InvalidateSession(ctx, sessionID)
}

func (u *AuthUseCase) GetSessionBySecret(ctx context.Context, sessionSecret string) (*models.UserSession, error) {
	return u.authService.GetSessionBySecret(ctx, sessionSecret)
}

func (u *AuthUseCase) ExchangeYandexAccessToken(ctx context.Context, accessToken string) (*models.UserSession, error) {
	userInfo, err := u.yandexOAuthService.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	user, err := u.authService.GetUserByEmail(ctx, userInfo.DefaultEmail)
	if err != nil {
		if errors.Is(err, yandex.ErrInvalidOrExpiredToken) {
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

	session, err := u.authService.CreateSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (uc *AuthUseCase) GetApiTokens(ctx context.Context) ([]*models.ApiToken, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	apiTokens, err := uc.authService.GetApiTokens(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return apiTokens, nil
}

func (uc *AuthUseCase) CreateApiToken(ctx context.Context, name string) (*models.ApiToken, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	apiToken, err := uc.authService.CreateApiToken(ctx, orgID, name)
	if err != nil {
		return nil, err
	}

	return apiToken, nil
}

func (uc *AuthUseCase) RevokeApiToken(ctx context.Context, id uuid.UUID) error {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = uc.authService.RevokeApiToken(ctx, orgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) GetEmployees(ctx context.Context) ([]*models.Employee, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	employees, err := uc.authService.GetEmployees(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (uc *AuthUseCase) GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	employee, err := uc.authService.GetEmployee(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *AuthUseCase) SetEmployeeRole(ctx context.Context, id uuid.UUID, roleID int) (*models.Employee, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = uc.authService.SetUserRole(ctx, orgID, id, models.RoleID(roleID))
	if err != nil {
		return nil, err
	}

	employee, err := uc.authService.GetEmployee(ctx, orgID, id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *AuthUseCase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return err
	}

	err = uc.authService.RemoveUserRole(ctx, orgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) InviteEmployee(ctx context.Context, email string, roleID int) (*models.Employee, error) {
	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user, err := uc.authService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = uc.authService.SetUserRole(ctx, orgID, user.ID, models.RoleID(roleID))
	if err != nil {
		return nil, err
	}

	employee, err := uc.authService.GetEmployee(ctx, orgID, user.ID)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *AuthUseCase) GetRoles(ctx context.Context) ([]*models.Role, error) {
	roles, err := uc.authService.GetAvailableRoles(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (uc *AuthUseCase) GetOrgIdByApiToken(ctx context.Context, token string) (uuid.UUID, error) {
	orgID, err := uc.authService.GetOrgIdByApiToken(ctx, token)
	if err != nil {
		return uuid.Nil, err
	}
	return orgID, nil
}
