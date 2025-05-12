package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/employee"
	"github.com/let-store-it/backend/internal/services/yandex"
	"github.com/let-store-it/backend/internal/usecases"
)

type AuthUseCase struct {
	authService        *auth.AuthService
	yandexOAuthService *yandex.YandexOAuthService
	employeeService    *employee.EmployeeService
}

type AuthUseCaseConfig struct {
	AuditService       *audit.AuditService
	AuthService        *auth.AuthService
	YandexOAuthService *yandex.YandexOAuthService
	EmployeeService    *employee.EmployeeService
}

func New(config AuthUseCaseConfig) *AuthUseCase {
	return &AuthUseCase{
		authService:        config.AuthService,
		yandexOAuthService: config.YandexOAuthService,
		employeeService:    config.EmployeeService,
	}
}

func (u *AuthUseCase) GetCurrentUser(ctx context.Context) (*models.User, error) {
	userID, err := common.GetUserIDFromContext(ctx)
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
		if errors.Is(err, common.ErrNotFound) {
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
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, false)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	apiTokens, err := uc.authService.GetApiTokens(ctx, valRes.OrgID)
	if err != nil {
		return nil, err
	}

	return apiTokens, nil
}

func (uc *AuthUseCase) CreateApiToken(ctx context.Context, name string) (*models.ApiToken, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, false)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	apiToken, err := uc.authService.CreateApiToken(ctx, valRes.OrgID, name)
	if err != nil {
		return nil, err
	}

	return apiToken, nil
}

func (uc *AuthUseCase) RevokeApiToken(ctx context.Context, id uuid.UUID) error {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, false)
	if err != nil {
		return err
	}

	if !valRes.IsAllowed {
		return usecases.ErrForbidden
	}

	err = uc.authService.RevokeApiToken(ctx, valRes.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *AuthUseCase) GetEmployees(ctx context.Context) ([]*models.Employee, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	employees, err := uc.employeeService.GetEmployees(ctx, valRes.OrgID)
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (uc *AuthUseCase) GetEmployee(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	employee, err := uc.employeeService.GetEmployee(ctx, valRes.OrgID, id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *AuthUseCase) SetEmployeeRole(ctx context.Context, id uuid.UUID, roleID int) (*models.Employee, error) {
	neededAccessLevel := models.AccessLevelAdmin
	if models.RoleID(roleID) == models.RoleOwnerID {
		neededAccessLevel = models.AccessLevelOwner
	}

	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, neededAccessLevel, true)
	if err != nil {
		return nil, err
	}

	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	if id == *valRes.UserID {
		return nil, errors.New("cannot set role for yourself")
	}

	err = uc.authService.SetUserRole(ctx, valRes.OrgID, id, models.RoleID(roleID))
	if err != nil {
		return nil, err
	}

	employee, err := uc.employeeService.GetEmployee(ctx, valRes.OrgID, id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (uc *AuthUseCase) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	orgID, err := common.GetOrganizationIDFromContext(ctx)
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
	neededAccessLevel := models.AccessLevelAdmin
	if models.RoleID(roleID) == models.RoleOwnerID {
		neededAccessLevel = models.AccessLevelOwner
	}

	valRes, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, neededAccessLevel, false)
	if err != nil {
		return nil, err
	}
	if !valRes.IsAllowed {
		return nil, usecases.ErrForbidden
	}
	user, err := uc.authService.GetUserByEmail(ctx, email)
	if err != nil {
		if database.IsNotFound(err) {
			return nil, errors.Join(common.ErrDetailedValidationErrorWithMessage("user not found"))
		}
		return nil, err
	}

	if user.ID == *valRes.UserID {
		return nil, common.ErrDetailedValidationErrorWithMessage("cannot invite yourself")
	}

	err = uc.authService.SetUserRole(ctx, valRes.OrgID, user.ID, models.RoleID(roleID))
	if err != nil {
		return nil, err
	}

	employee, err := uc.employeeService.GetEmployee(ctx, valRes.OrgID, user.ID)
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

func (u *AuthUseCase) CreateUser(ctx context.Context, email, firstName, lastName string) (*models.User, error) {
	user := &models.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	return u.authService.CreateUser(ctx, user)
}
