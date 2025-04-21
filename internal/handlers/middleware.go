package handlers

import (
	"context"
	"fmt"
	"net/http"

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
