package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/usecases"
)

func WithOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgIDStr := r.Header.Get("x-organization-id")
		orgID := uuid.Nil
		if orgIDStr != "" {
			parsedOrgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid organization ID: %v", err), http.StatusBadRequest)
				return
			}
			orgID = parsedOrgID
		}

		ctx := context.WithValue(r.Context(), usecases.OrganizationIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type AuthMiddleware struct {
	skipPathsPrefix []string
	authUseCase     *usecases.AuthUseCase
	cookieName      string
}

func NewAuthMiddleware(authUseCase *usecases.AuthUseCase, cookieName string, skipPaths []string) *AuthMiddleware {
	return &AuthMiddleware{authUseCase: authUseCase, cookieName: cookieName, skipPathsPrefix: skipPaths}
}

func (m *AuthMiddleware) Process(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				for _, path := range m.skipPathsPrefix {
					if strings.HasPrefix(r.URL.Path, path) {
						next.ServeHTTP(w, r)
						return
					}
				}
				slog.Info("cookie not present, user not logged in")
				http.Error(w, fmt.Sprintf("cookie %s not present, authorize first", m.cookieName), http.StatusUnauthorized)
				return
			}
			slog.Error("error getting cookie", "error", err)
			http.Error(w, "error getting cookie", http.StatusBadRequest)
			return
		}

		userID, err := m.authUseCase.GetUserIdFromSession(r.Context(), cookie.Value)
		if err != nil {
			http.SetCookie(w, &http.Cookie{Name: m.cookieName, Value: "", Expires: time.Now().Add(-1 * time.Hour)})
			http.Error(w, "invalid session, cookie was reset", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), usecases.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
