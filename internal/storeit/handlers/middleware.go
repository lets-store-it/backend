package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/usecases"
)

type OrganizationIDMiddleware struct {
	orgUseCase  *usecases.OrganizationUseCase
	authUseCase *usecases.AuthUseCase
}

func NewOrganizationIDMiddleware(orgUseCase *usecases.OrganizationUseCase, authUseCase *usecases.AuthUseCase) *OrganizationIDMiddleware {
	return &OrganizationIDMiddleware{
		orgUseCase:  orgUseCase,
		authUseCase: authUseCase,
	}
}

func (m *OrganizationIDMiddleware) WithOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/auth") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("storeit_session")
		if err != nil {
			http.Error(w, "storeit_session cookie is required", http.StatusBadRequest)
			return
		}
		userID, err := m.authUseCase.GetUserIdFromSession(r.Context(), cookie.Value)
		if err != nil {
			http.Error(w, "invalid session", http.StatusBadRequest)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/me") {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), usecases.UserIDKey, userID)))
			return
		}

		if strings.HasPrefix(r.URL.Path, "/orgs") {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), usecases.UserIDKey, userID)))
			return
		}

		orgIDStr := r.Header.Get("x-organization-id")
		if orgIDStr == "" {
			http.Error(w, "x-organization-id header is required", http.StatusBadRequest)
			return
		}

		orgID, err := uuid.Parse(orgIDStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid organization ID: %v", err), http.StatusBadRequest)
			return
		}

		exists, err := m.orgUseCase.IsOrganizationExists(r.Context(), orgID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to check organization existence: %v", err), http.StatusInternalServerError)
			return
		}

		if !exists {
			http.Error(w, "organization not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), usecases.UserIDKey, userID)
		ctx = context.WithValue(ctx, usecases.OrganizationIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
