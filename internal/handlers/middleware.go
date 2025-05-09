package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services"
)

const (
	// Cookie configuration
	cookieName       = "storeit_session"
	refreshThreshold = 6 * 24 * time.Hour   // 6 days
	defaultExpiresIn = 360 * 24 * time.Hour // 360 days
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionRevoked  = errors.New("session revoked")
	ErrSessionExpired  = errors.New("session expired")
	ErrSessionInvalid  = errors.New("session invalid")
)

func WithOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var orgID uuid.UUID = uuid.Nil

		orgIDStr := r.Header.Get("x-organization-id")
		if orgIDStr != "" {
			var err error
			orgID, err = uuid.Parse(orgIDStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid organization ID format: %v", err), http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), models.OrganizationIDContextKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *RestApiImplementation) HandleApiToken(ctx context.Context, operationName api.OperationName, t api.ApiToken) (context.Context, error) {
	orgID, err := h.authUseCase.GetOrgIdByApiToken(ctx, t.GetAPIKey())
	if err != nil {
		if errors.Is(err, services.ErrNotFoundError) {
			return nil, h.NewUnauthorizedError(ctx)
		}
		return nil, fmt.Errorf("failed to process API token: %w", err)
	}

	ctx = context.WithValue(ctx, models.OrganizationIDContextKey, orgID)
	ctx = context.WithValue(ctx, models.IsSystemUserContextKey, true)
	return ctx, nil
}

func generateAuthCookie(token string, timeToLive time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(timeToLive.Seconds()),
		HttpOnly: true,
		Secure:   false, // TODO: Set to true in production
		SameSite: http.SameSiteLaxMode,
	}
}

func (h *RestApiImplementation) refreshSessionIfNeeded(ctx context.Context, session *models.UserSession) (*models.UserSession, error) {
	if session.ExpiresAt == nil || time.Until(*session.ExpiresAt) > refreshThreshold {
		return session, nil
	}

	newSession, err := h.authUseCase.CreateSession(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("error creating new session: %w", err)
	}

	// Set new session cookie
	expiresIn := defaultExpiresIn
	if newSession.ExpiresAt != nil {
		expiresIn = time.Until(*newSession.ExpiresAt)
	}

	cookie := generateAuthCookie(newSession.Token, expiresIn)
	if cookieSetFn, ok := ctx.Value(models.SetCookieContextKey).(func(*http.Cookie)); ok {
		cookieSetFn(cookie)
	}

	err = h.authUseCase.InvalidateSession(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("error invalidating session: %w", err)
	}

	return newSession, nil
}

func (h *RestApiImplementation) HandleCookie(ctx context.Context, operationName api.OperationName, t api.Cookie) (context.Context, error) {
	session, err := h.authUseCase.GetSessionBySecret(ctx, t.GetAPIKey())
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.RevokedAt != nil {
		return nil, ErrSessionRevoked
	}

	if session.ExpiresAt != nil && session.ExpiresAt.Before(time.Now()) {
		return nil, ErrSessionExpired
	}

	session, err = h.refreshSessionIfNeeded(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	return context.WithValue(ctx, models.UserIDContextKey, session.UserID), nil
}

func WithSetCookieFromContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), models.SetCookieContextKey, func(cookie *http.Cookie) {
			http.SetCookie(w, cookie)
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
