package handlers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
)

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

	return &api.AuthResponse{
		SetCookie: generateAuthCookie(session.Token, 60*60*24*30).String(), // 30 days
	}, nil
}

func (h *RestApiImplementation) Logout(ctx context.Context) (api.LogoutRes, error) {
	return &api.LogoutResponse{
		SetCookie: generateAuthCookie("", -1).String(), // Expire cookie
	}, nil
}

func toApiToken(token *models.ApiToken) api.Token {
	return api.Token{
		ID:    token.ID,
		Token: token.Token,
		Name:  token.Name,
	}
}

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

func (h *RestApiImplementation) CreateApiToken(ctx context.Context, req *api.CreateApiTokenRequest) (api.CreateApiTokenRes, error) {
	apiToken, err := h.authUseCase.CreateApiToken(ctx, req.Name)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.CreateApiTokenResponse{
		Data: toApiToken(apiToken),
	}, nil
}

func (h *RestApiImplementation) RevokeApiToken(ctx context.Context, params api.RevokeApiTokenParams) (api.RevokeApiTokenRes, error) {
	err := h.authUseCase.RevokeApiToken(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.RevokeApiTokenNoContent{}, nil
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

func (h *RestApiImplementation) DeleteEmployeeById(ctx context.Context, params api.DeleteEmployeeByIdParams) (api.DeleteEmployeeByIdRes, error) {
	err := h.authUseCase.DeleteEmployee(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.DeleteEmployeeByIdOK{}, nil
}

func (h *RestApiImplementation) GetEmployeeById(ctx context.Context, params api.GetEmployeeByIdParams) (api.GetEmployeeByIdRes, error) {
	employee, err := h.authUseCase.GetEmployee(ctx, params.ID)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.GetEmployeeResponse{
		Data: toEmployeeDTO(employee),
	}, nil
}

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

func (h *RestApiImplementation) InviteEmployee(ctx context.Context, req *api.InviteEmployeeRequest) (api.InviteEmployeeRes, error) {
	employee, err := h.authUseCase.InviteEmployee(ctx, req.Email, req.RoleId)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.GetEmployeeResponse{
		Data: toEmployeeDTO(employee),
	}, nil
}

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

func (h *RestApiImplementation) TestAuth(c echo.Context) error {
	var body struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	if err := c.Bind(&body); err != nil {
		return h.NewError(c.Request().Context(), err)
	}
	if body.Email == "" || body.FirstName == "" || body.LastName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "email, firstName and lastName are required",
		})
	}

	session, err := h.authUseCase.CreateSessionByEmail(c.Request().Context(), body.Email)
	if err != nil {
		if err == common.ErrNotFound {
			user, err := h.authUseCase.CreateUser(c.Request().Context(), body.Email, body.FirstName, body.LastName)
			if err != nil {
				return h.NewError(c.Request().Context(), err)
			}

			session, err = h.authUseCase.CreateSession(c.Request().Context(), user.ID)
			if err != nil {
				return h.NewError(c.Request().Context(), err)
			}
		} else {
			return h.NewError(c.Request().Context(), err)
		}
	}

	cookie := generateAuthCookie(session.Token, defaultExpiresIn)
	c.SetCookie(cookie)

	return c.NoContent(http.StatusOK)
}
