package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/usecases"
	"github.com/ogen-go/ogen/ogenerrors"
)

func (h *RestApiImplementation) NewConflictError(ctx context.Context, message string) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusConflict,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "conflict",
				Message: message,
			},
		},
	}
}

func (h *RestApiImplementation) NewUnauthorizedError(ctx context.Context) *api.DefaultErrorStatusCode {
	return h.NewUnauthorizedErrorWithMessage(ctx, "you are not authorized to access this resource")
}

func (h *RestApiImplementation) NewUnauthorizedErrorWithMessage(ctx context.Context, message string) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusUnauthorized,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "unauthorized",
				Message: message,
			},
		},
	}
}

func (h *RestApiImplementation) NewValidationError(ctx context.Context, message string) *api.DefaultErrorStatusCode {
	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusBadRequest,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "validation_error",
				Message: message,
			},
		},
	}
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.DefaultErrorStatusCode {
	// var ogenErr ogenerrors.Error
	var detailedErr *usecases.ErrDetailedValidationError
	switch {
	case errors.As(err, &detailedErr):
		return h.NewValidationError(ctx, detailedErr.Message)
	case errors.Is(err, usecases.ErrNotAuthorized):
		return h.NewUnauthorizedError(ctx)
	case errors.Is(err, usecases.ErrOrganizationIDMissing):
		return h.NewValidationError(ctx, "x-organization-id header is missing")
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		return h.NewUnauthorizedError(ctx)
	case errors.Is(err, ErrSessionNotFound):
		return h.NewUnauthorizedError(ctx)
	case errors.Is(err, ErrSessionRevoked):
		return h.NewUnauthorizedErrorWithMessage(ctx, "Session was revoked")
	case errors.Is(err, ErrSessionExpired):
		return h.NewUnauthorizedErrorWithMessage(ctx, "Session expired")
		// case errors.As(err, &ogenErr):
		// code = ogenErr.
	}

	return &api.DefaultErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api.ErrorContent{
			Error: api.ErrorContentError{
				Code:    "internal_server_error",
				Message: err.Error(),
			},
		},
	}
}
