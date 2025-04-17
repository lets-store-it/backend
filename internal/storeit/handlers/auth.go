package handlers

import (
	"context"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
)

func (h *RestApiImplementation) GetAuthCookieByEmail(ctx context.Context, req *api.GetAuthCookieByEmailRequest) (*api.GetAuthCookieByEmailOK, error) {
	session, err := h.authUseCase.CreateSessionByEmail(ctx, req.Email)
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

	return &api.GetAuthCookieByEmailOK{
		SetCookie: cookie.String(),
	}, nil
}

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
