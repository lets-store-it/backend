package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/usecases"
)

// HandleApiToken implements api.SecurityHandler.
func (h *RestApiImplementation) HandleApiToken(ctx context.Context, operationName api.OperationName, t api.ApiToken) (context.Context, error) {
	orgID, err := h.authUseCase.GetOrgIdByApiToken(ctx, t.GetAPIKey())
	if err != nil {
		if errors.Is(err, auth.ErrApiTokenNotFound) {
			return nil, h.NewUnauthorizedError(ctx)
		}
		return nil, fmt.Errorf("error processing api token")
	}
	ctx = context.WithValue(ctx, usecases.OrganizationIDKey, orgID)
	ctx = context.WithValue(ctx, usecases.IsSystemUserKey, true)
	return ctx, nil
}

// HandleCookie implements api.SecurityHandler.
func (h *RestApiImplementation) HandleCookie(ctx context.Context, operationName api.OperationName, t api.Cookie) (context.Context, error) {
	userID, err := h.authUseCase.GetUserIdFromSession(ctx, t.GetAPIKey())
	if err != nil {
		if errors.Is(err, auth.ErrSessionNotFound) {
			return nil, h.NewUnauthorizedError(ctx)
		}
		return nil, fmt.Errorf("error processing session")
	}

	return context.WithValue(ctx, usecases.UserIDKey, userID), nil
}

// / GetCurrentUser implements api.Handler.
func (h *RestApiImplementation) GetCurrentUser(ctx context.Context) (api.GetCurrentUserRes, error) {
	user, err := h.authUseCase.GetCurrentUser(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	var middleName api.NilString
	if user.MiddleName != nil {
		middleName.Value = *user.MiddleName
	}

	return &api.GetCurrentUserResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: middleName,
	}, nil
}

func (h *RestApiImplementation) ExchangeYandexAccessToken(ctx context.Context, req *api.ExchangeYandexAccessTokenReq) (api.ExchangeYandexAccessTokenRes, error) {
	session, err := h.authUseCase.ExchangeYandexAccessToken(ctx, req.AccessToken)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	cookie := &http.Cookie{
		Name:     "storeit_session",
		Value:    session.Secret,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &api.AuthResponse{
		SetCookie: cookie.String(),
	}, nil
}

func (h *RestApiImplementation) Logout(ctx context.Context) (api.LogoutRes, error) {
	cookie := &http.Cookie{
		Name:     "storeit_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	return &api.LogoutResponse{
		SetCookie: cookie.String(),
	}, nil
}

func toApiToken(token *models.ApiToken) api.Token {
	return api.Token{
		ID:    token.ID,
		Token: token.Token,
		Name:  token.Name,
	}
}

// GetApiTokens implements api.Handler.
func (h *RestApiImplementation) GetApiTokens(ctx context.Context) (api.GetApiTokensRes, error) {
	apiTokens, err := h.authUseCase.GetApiTokens(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}

	tokens := make([]api.Token, len(apiTokens))
	for i, token := range apiTokens {
		tokens[i] = toApiToken(token)
	}
	return &api.GetApiTokensResponse{
		Data: tokens,
	}, nil
}

// CreateApiToken implements api.Handler.
func (h *RestApiImplementation) CreateApiToken(ctx context.Context, req *api.CreateApiTokenRequest) (api.CreateApiTokenRes, error) {
	apiToken, err := h.authUseCase.CreateApiToken(ctx, req.Name)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.CreateApiTokenResponse{
		Data: toApiToken(apiToken),
	}, nil
}

// RevokeApiToken implements api.Handler.
func (h *RestApiImplementation) RevokeApiToken(ctx context.Context, params api.RevokeApiTokenParams) (api.RevokeApiTokenRes, error) {
	err := h.authUseCase.RevokeApiToken(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return nil, nil
}

func toRoleDTO(role *models.Role) api.Role {
	return api.Role{
		ID:          role.ID,
		Name:        string(role.Name),
		DisplayName: role.DisplayName,
		Description: role.Description,
	}
}

func toEmployeeDTO(employee *models.Employee) api.Employee {
	var middleName api.NilString
	if employee.MiddleName != nil {
		middleName.Value = *employee.MiddleName
	}
	return api.Employee{
		UserId:     employee.UserID,
		Email:      employee.Email,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		MiddleName: middleName,
		Role:       toRoleDTO(employee.Role),
	}
}

// GetEmployees implements api.Handler.
func (h *RestApiImplementation) GetEmployees(ctx context.Context) (api.GetEmployeesRes, error) {
	employees, err := h.authUseCase.GetEmployees(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	employeesDTO := make([]api.Employee, len(employees))
	for i, employee := range employees {
		employeesDTO[i] = toEmployeeDTO(employee)
	}
	return &api.GetEmployeesResponse{
		Data: employeesDTO,
	}, nil
}

// DeleteEmployeeById implements api.Handler.
func (h *RestApiImplementation) DeleteEmployeeById(ctx context.Context, params api.DeleteEmployeeByIdParams) (api.DeleteEmployeeByIdRes, error) {
	err := h.authUseCase.DeleteEmployee(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.DeleteEmployeeByIdOK{}, nil
}

// GetEmployeeById implements api.Handler.
func (h *RestApiImplementation) GetEmployeeById(ctx context.Context, params api.GetEmployeeByIdParams) (api.GetEmployeeByIdRes, error) {
	employee, err := h.authUseCase.GetEmployee(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.GetEmployeeResponse{
		Data: toEmployeeDTO(employee),
	}, nil
}

// GetRoles implements api.Handler.
func (h *RestApiImplementation) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	roles, err := h.authUseCase.GetRoles(ctx)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	rolesDTO := make([]api.Role, len(roles))
	for i, role := range roles {
		rolesDTO[i] = toRoleDTO(role)
	}
	return &api.GetRolesOK{
		Data: rolesDTO,
	}, nil
}

// InviteEmployee implements api.Handler.
func (h *RestApiImplementation) InviteEmployee(ctx context.Context, req *api.InviteEmployeeRequest) (api.InviteEmployeeRes, error) {
	employee, err := h.authUseCase.InviteEmployee(ctx, req.Email, req.RoleId)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.GetEmployeeResponse{
		Data: toEmployeeDTO(employee),
	}, nil
}

// PatchEmployeeById implements api.Handler.
func (h *RestApiImplementation) PatchEmployeeById(ctx context.Context, req *api.PatchEmployeeRequest, params api.PatchEmployeeByIdParams) (api.PatchEmployeeByIdRes, error) {
	roleID := req.RoleId.Value
	employee, err := h.authUseCase.SetEmployeeRole(ctx, params.ID, roleID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.GetEmployeeResponse{
		Data: toEmployeeDTO(employee),
	}, nil
}
