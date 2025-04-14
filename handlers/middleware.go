package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func WithOrganizationID(next http.Handler) http.Handler {
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

		ctx := context.WithValue(r.Context(), "organization_id", orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
