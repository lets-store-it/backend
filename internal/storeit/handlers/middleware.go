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
	orgUseCase *usecases.OrganizationUseCase
}

func NewOrganizationIDMiddleware(orgUseCase *usecases.OrganizationUseCase) *OrganizationIDMiddleware {
	return &OrganizationIDMiddleware{
		orgUseCase: orgUseCase,
	}
}

func (m *OrganizationIDMiddleware) WithOrganizationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip organization header check for /orgs paths
		if strings.HasPrefix(r.URL.Path, "/orgs") {
			next.ServeHTTP(w, r)
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

		ctx := context.WithValue(r.Context(), usecases.OrganizationIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
