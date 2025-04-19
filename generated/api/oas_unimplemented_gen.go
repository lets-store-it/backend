// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// CreateCell implements createCell operation.
//
// Create Cells.
//
// POST /cells-groups/{groupId}/cells
func (UnimplementedHandler) CreateCell(ctx context.Context, req *CreateCellRequest, params CreateCellParams) (r *CreateCellResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateCellsGroup implements createCellsGroup operation.
//
// Create Cells Group.
//
// POST /cells-groups
func (UnimplementedHandler) CreateCellsGroup(ctx context.Context, req *CreateCellsGroupRequest) (r *CreateCellsGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateItem implements createItem operation.
//
// Create Item.
//
// POST /items
func (UnimplementedHandler) CreateItem(ctx context.Context, req *CreateItemRequest) (r *CreateItemResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateOrganization implements createOrganization operation.
//
// Create Organization.
//
// POST /orgs
func (UnimplementedHandler) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest) (r *CreateOrganizationResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateStorageGroup implements createStorageGroup operation.
//
// Create Storage Group.
//
// POST /storage-groups
func (UnimplementedHandler) CreateStorageGroup(ctx context.Context, req *CreateStorageGroupRequest) (r *CreateStorageGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// CreateUnit implements createUnit operation.
//
// Create Organization Unit.
//
// POST /units
func (UnimplementedHandler) CreateUnit(ctx context.Context, req *CreateOrganizationUnitRequest) (r *CreateOrganizationUnitResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// DeleteCell implements deleteCell operation.
//
// Delete Cell.
//
// DELETE /cells-groups/{groupId}/cells/{cellId}
func (UnimplementedHandler) DeleteCell(ctx context.Context, params DeleteCellParams) error {
	return ht.ErrNotImplemented
}

// DeleteCellsGroup implements deleteCellsGroup operation.
//
// Delete Cells Group.
//
// DELETE /cells-groups/{groupId}
func (UnimplementedHandler) DeleteCellsGroup(ctx context.Context, params DeleteCellsGroupParams) error {
	return ht.ErrNotImplemented
}

// DeleteItem implements deleteItem operation.
//
// Delete Item.
//
// DELETE /items/{id}
func (UnimplementedHandler) DeleteItem(ctx context.Context, params DeleteItemParams) error {
	return ht.ErrNotImplemented
}

// DeleteOrganization implements deleteOrganization operation.
//
// Delete Organization.
//
// DELETE /orgs/{id}
func (UnimplementedHandler) DeleteOrganization(ctx context.Context, params DeleteOrganizationParams) error {
	return ht.ErrNotImplemented
}

// DeleteOrganizationUnit implements deleteOrganizationUnit operation.
//
// Delete Organization Unit.
//
// DELETE /units/{id}
func (UnimplementedHandler) DeleteOrganizationUnit(ctx context.Context, params DeleteOrganizationUnitParams) error {
	return ht.ErrNotImplemented
}

// DeleteStorageGroup implements deleteStorageGroup operation.
//
// Delete Storage Group.
//
// DELETE /storage-groups/{id}
func (UnimplementedHandler) DeleteStorageGroup(ctx context.Context, params DeleteStorageGroupParams) error {
	return ht.ErrNotImplemented
}

// ExchangeYandexAccessToken implements exchangeYandexAccessToken operation.
//
// Exchange Yandex Access token for Session token.
//
// POST /auth/oauth2/yandex
func (UnimplementedHandler) ExchangeYandexAccessToken(ctx context.Context, req *ExchangeYandexAccessTokenReq) (r *AuthResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCellById implements getCellById operation.
//
// Get Cell by ID.
//
// GET /cells-groups/{groupId}/cells/{cellId}
func (UnimplementedHandler) GetCellById(ctx context.Context, params GetCellByIdParams) (r *GetCellByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCells implements getCells operation.
//
// Get list of Cells.
//
// GET /cells-groups/{groupId}/cells
func (UnimplementedHandler) GetCells(ctx context.Context, params GetCellsParams) (r *GetCellsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCellsGroupById implements getCellsGroupById operation.
//
// Get Cells Group by ID.
//
// GET /cells-groups/{groupId}
func (UnimplementedHandler) GetCellsGroupById(ctx context.Context, params GetCellsGroupByIdParams) (r *GetCellsGroupByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCellsGroups implements getCellsGroups operation.
//
// Get list of Cells Groups.
//
// GET /cells-groups
func (UnimplementedHandler) GetCellsGroups(ctx context.Context) (r *GetCellsGroupsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetCurrentUser implements getCurrentUser operation.
//
// Get Current User.
//
// GET /me
func (UnimplementedHandler) GetCurrentUser(ctx context.Context) (r *GetCurrentUserResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetItemById implements getItemById operation.
//
// Get Item by ID.
//
// GET /items/{id}
func (UnimplementedHandler) GetItemById(ctx context.Context, params GetItemByIdParams) (r *GetItemByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetItems implements getItems operation.
//
// Get list of Items.
//
// GET /items
func (UnimplementedHandler) GetItems(ctx context.Context) (r *GetItemsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOrganizationById implements getOrganizationById operation.
//
// Get Organization by ID.
//
// GET /orgs/{id}
func (UnimplementedHandler) GetOrganizationById(ctx context.Context, params GetOrganizationByIdParams) (r *GetOrganizationByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOrganizationUnitById implements getOrganizationUnitById operation.
//
// Get Unit by ID.
//
// GET /units/{id}
func (UnimplementedHandler) GetOrganizationUnitById(ctx context.Context, params GetOrganizationUnitByIdParams) (r *GetOrganizationUnitByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOrganizationUnits implements getOrganizationUnits operation.
//
// Get list of Organization Units.
//
// GET /units
func (UnimplementedHandler) GetOrganizationUnits(ctx context.Context) (r *GetOrganizationUnitsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetOrganizations implements getOrganizations operation.
//
// Get list of Organizations.
//
// GET /orgs
func (UnimplementedHandler) GetOrganizations(ctx context.Context) (r *GetOrganizationsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetStorageGroupById implements getStorageGroupById operation.
//
// Get Storage Group by ID.
//
// GET /storage-groups/{id}
func (UnimplementedHandler) GetStorageGroupById(ctx context.Context, params GetStorageGroupByIdParams) (r *GetStorageGroupByIdResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// GetStorageGroups implements getStorageGroups operation.
//
// Get list of Storage Groups.
//
// GET /storage-groups
func (UnimplementedHandler) GetStorageGroups(ctx context.Context) (r *GetStorageGroupsResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// Logout implements logout operation.
//
// Logout user.
//
// GET /auth/logout
func (UnimplementedHandler) Logout(ctx context.Context) (r *LogoutResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchCell implements patchCell operation.
//
// Patch Cell.
//
// PATCH /cells-groups/{groupId}/cells/{cellId}
func (UnimplementedHandler) PatchCell(ctx context.Context, req *PatchCellRequest, params PatchCellParams) (r *PatchCellResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchCellsGroup implements patchCellsGroup operation.
//
// Patch Cells Group.
//
// PATCH /cells-groups/{groupId}
func (UnimplementedHandler) PatchCellsGroup(ctx context.Context, req *PatchCellsGroupRequest, params PatchCellsGroupParams) (r *PatchCellsGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchItem implements patchItem operation.
//
// Patch Item.
//
// PATCH /items/{id}
func (UnimplementedHandler) PatchItem(ctx context.Context, req *PatchItemRequest, params PatchItemParams) (r *PatchItemResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchOrganization implements patchOrganization operation.
//
// Update Organization.
//
// PATCH /orgs/{id}
func (UnimplementedHandler) PatchOrganization(ctx context.Context, req *PatchOrganizationRequest, params PatchOrganizationParams) (r *PatchOrganizationResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchOrganizationUnit implements patchOrganizationUnit operation.
//
// Patch Organization Unit.
//
// PATCH /units/{id}
func (UnimplementedHandler) PatchOrganizationUnit(ctx context.Context, req *PatchOrganizationUnitRequest, params PatchOrganizationUnitParams) (r *PatchOrganizationUnitResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// PatchStorageGroup implements patchStorageGroup operation.
//
// Patch Storage Group.
//
// PATCH /storage-groups/{id}
func (UnimplementedHandler) PatchStorageGroup(ctx context.Context, req *PatchStorageGroupRequest, params PatchStorageGroupParams) (r *PatchStorageGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateCell implements updateCell operation.
//
// Update Cell.
//
// PUT /cells-groups/{groupId}/cells/{cellId}
func (UnimplementedHandler) UpdateCell(ctx context.Context, req *UpdateCellRequest, params UpdateCellParams) (r *UpdateCellResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateCellsGroup implements updateCellsGroup operation.
//
// Update Cells Group.
//
// PUT /cells-groups/{groupId}
func (UnimplementedHandler) UpdateCellsGroup(ctx context.Context, req *UpdateCellsGroupRequest, params UpdateCellsGroupParams) (r *UpdateCellsGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateItem implements updateItem operation.
//
// Update Item.
//
// PUT /items/{id}
func (UnimplementedHandler) UpdateItem(ctx context.Context, req *UpdateItemRequest, params UpdateItemParams) (r *UpdateItemResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateOrganization implements updateOrganization operation.
//
// Update Organization.
//
// PUT /orgs/{id}
func (UnimplementedHandler) UpdateOrganization(ctx context.Context, req *UpdateOrganizationRequest, params UpdateOrganizationParams) (r *UpdateOrganizationResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateOrganizationUnit implements updateOrganizationUnit operation.
//
// Update Organization Unit.
//
// PUT /units/{id}
func (UnimplementedHandler) UpdateOrganizationUnit(ctx context.Context, req *UpdateOrganizationUnitRequest, params UpdateOrganizationUnitParams) (r *UpdateOrganizationUnitResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// UpdateStorageGroup implements updateStorageGroup operation.
//
// Update Storage Group.
//
// PUT /storage-groups/{id}
func (UnimplementedHandler) UpdateStorageGroup(ctx context.Context, req *UpdateStorageGroupRequest, params UpdateStorageGroupParams) (r *UpdateStorageGroupResponse, _ error) {
	return r, ht.ErrNotImplemented
}

// NewError creates *DefaultErrorStatusCode from error returned by handler.
//
// Used for common default response.
func (UnimplementedHandler) NewError(ctx context.Context, err error) (r *DefaultErrorStatusCode) {
	r = new(DefaultErrorStatusCode)
	return r
}
