package handlers

import (
	"context"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

// func (h *RestApiImplementation) GetCurrentUserBySessionSecret(ctx context.Context) (*api.GetCurrentUserResponse, error) {

// 	userID := h.authUseCase.GetCurrentUser(ctx)
// 	user, err := h.authUseCase.GetCurrentUser(ctx, userID)
// 	if err != nil {
// 		return nil, h.NewError(ctx, err)
// 	}

// 	return user, nil
// }

// / GetCurrentUser implements api.Handler.
func (h *RestApiImplementation) GetCurrentUser(ctx context.Context) (*api.GetCurrentUserResponse, error) {
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

func (h *RestApiImplementation) ExchangeYandexAccessToken(ctx context.Context, req *api.ExchangeYandexAccessTokenReq) (*api.AuthResponse, error) {
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

func (h *RestApiImplementation) Logout(ctx context.Context) (*api.LogoutResponse, error) {
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
func (h *RestApiImplementation) GetApiTokens(ctx context.Context) (*api.GetApiTokensResponse, error) {
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
func (h *RestApiImplementation) CreateApiToken(ctx context.Context, req *api.CreateApiTokenRequest) (*api.CreateApiTokenResponse, error) {
	apiToken, err := h.authUseCase.CreateApiToken(ctx, req.Name)
	if err != nil {
		return nil, h.NewError(ctx, err)
	}
	return &api.CreateApiTokenResponse{
		Data: toApiToken(apiToken),
	}, nil
}

// RevokeApiToken implements api.Handler.
func (h *RestApiImplementation) RevokeApiToken(ctx context.Context, params api.RevokeApiTokenParams) error {
	err := h.authUseCase.RevokeApiToken(ctx, params.ID)
	if err != nil {
		return h.NewError(ctx, err)
	}
	return nil
}
