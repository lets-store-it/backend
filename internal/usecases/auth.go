package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/yandex"
)

type AuthUseCase struct {
	authService        *auth.AuthService
	yandexOAuthService *yandex.YandexOAuthService
}

func NewAuthUseCase(authService *auth.AuthService, yandexOAuthService *yandex.YandexOAuthService) *AuthUseCase {
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
		if errors.Is(err, auth.ErrUserNotFound) {
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

	isSystemUser := ctx.Value(IsSystemUserKey).(bool)
	if isSystemUser {
		return orgID, nil
	}

	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	_, err = uc.authService.GetUserRole(ctx, userID, orgID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	// if roles != auth.RoleOwner {
	// 	return uuid.Nil, fmt.Errorf("user is not an owner of the organization")
	// }

	return orgID, nil
}

func (uc *AuthUseCase) GetApiTokens(ctx context.Context) ([]*models.ApiToken, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	err = uc.authService.SetEmployeeRole(ctx, orgID, id, roleID)
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
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return err
	}

	err = uc.authService.DeleteEmployee(ctx, orgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) InviteEmployee(ctx context.Context, email string, roleID int) (*models.Employee, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}
	user, err := uc.authService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	err = uc.authService.AssignRoleToUser(ctx, orgID, user.ID, roleID)
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
	roles, err := uc.authService.GetRoles(ctx)
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
