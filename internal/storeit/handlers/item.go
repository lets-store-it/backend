package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
)

// CreateItem implements api.Handler.
func (h *RestApiImplementation) CreateItem(ctx context.Context, req *api.CreateItemRequest) (*api.CreateItemResponse, error) {
	panic("unimplemented")
}

// DeleteItem implements api.Handler.
func (h *RestApiImplementation) DeleteItem(ctx context.Context, params api.DeleteItemParams) error {
	panic("unimplemented")
}

// GetItemById implements api.Handler.
func (h *RestApiImplementation) GetItemById(ctx context.Context, params api.GetItemByIdParams) (*api.GetItemByIdResponse, error) {
	panic("unimplemented")
}

// GetItems implements api.Handler.
func (h *RestApiImplementation) GetItems(ctx context.Context) (*api.GetItemsResponse, error) {
	panic("unimplemented")
}

// PatchItem implements api.Handler.
func (h *RestApiImplementation) PatchItem(ctx context.Context, req *api.PatchItemRequest, params api.PatchItemParams) (*api.PatchItemResponse, error) {
	panic("unimplemented")
}

// UpdateItem implements api.Handler.
func (h *RestApiImplementation) UpdateItem(ctx context.Context, req *api.UpdateItemRequest, params api.UpdateItemParams) (*api.UpdateItemResponse, error) {
	panic("unimplemented")
}
