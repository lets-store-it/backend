package handlers

import (
	"context"
	"net/http"

	"github.com/let-store-it/backend/generated/api"
)

// CreateStorageGroup implements api.Handler.
func (h *RestApiImplementation) CreateStorageGroup(ctx context.Context, req *api.CreateStorageGroupRequest) (*api.CreateStorageGroupResponse, error) {
	panic("unimplemented")
}

// DeleteStorageGroup implements api.Handler.
func (h *RestApiImplementation) DeleteStorageGroup(ctx context.Context, params api.DeleteStorageGroupParams) error {
	panic("unimplemented")
}

// GetStorageGroupById implements api.Handler.
func (h *RestApiImplementation) GetStorageGroupById(ctx context.Context, params api.GetStorageGroupByIdParams) (*api.GetStorageGroupByIdResponse, error) {
	panic("unimplemented")
}

// GetStorageGroups implements api.Handler.
func (h *RestApiImplementation) GetStorageGroups(ctx context.Context) (*api.GetStorageGroupsResponse, error) {
	panic("unimplemented")
}

// PatchStorageGroup implements api.Handler.
func (h *RestApiImplementation) PatchStorageGroup(ctx context.Context, req *api.PatchStorageGroupRequest, params api.PatchStorageGroupParams) (*api.PatchStorageGroupResponse, error) {
	panic("unimplemented")
}

// UpdateStorageGroup implements api.Handler.
func (h *RestApiImplementation) UpdateStorageGroup(ctx context.Context, req *api.UpdateStorageGroupRequest, params api.UpdateStorageGroupParams) (*api.UpdateStorageGroupResponse, error) {
	panic("unimplemented")
}

func (h *RestApiImplementation) NewError(ctx context.Context, err error) *api.ErrorStatusCode {
	return createErrorResponse(http.StatusInternalServerError, err.Error())
}
