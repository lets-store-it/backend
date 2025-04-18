package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/storeit/models"
)

func convertGroupToDTO(group *models.StorageGroup) api.StorageGroup {
	var parentID api.NilUUID
	if group.ParentID == nil {
		parentID.SetToNull()
	} else {
		parentID.SetTo(*group.ParentID)
	}

	return api.StorageGroup{
		ID:       group.ID,
		ParentId: parentID,
		Name:     group.Name,
		Alias:    group.Alias,
		UnitId:   group.UnitID,
	}
}

// CreateStorageGroup implements api.Handler.
func (h *RestApiImplementation) CreateStorageGroup(ctx context.Context, req *api.CreateStorageGroupRequest) (*api.CreateStorageGroupResponse, error) {
	var parentID *uuid.UUID
	if val, ok := req.ParentId.Get(); ok {
		parentID = &val
	}

	group, err := h.storageGroupUseCase.Create(ctx, req.UnitId, parentID, req.Name, req.Alias)
	if err != nil {
		return nil, err
	}

	groupDTO := convertGroupToDTO(group)
	return &api.CreateStorageGroupResponse{
		Data: groupDTO,
	}, nil
}

// DeleteStorageGroup implements api.Handler.
func (h *RestApiImplementation) DeleteStorageGroup(ctx context.Context, params api.DeleteStorageGroupParams) error {
	return h.storageGroupUseCase.Delete(ctx, params.ID)
}

// GetStorageGroupById implements api.Handler.
func (h *RestApiImplementation) GetStorageGroupById(ctx context.Context, params api.GetStorageGroupByIdParams) (*api.GetStorageGroupByIdResponse, error) {
	group, err := h.storageGroupUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetStorageGroupByIdResponse{
		Data: convertGroupToDTO(group),
	}, nil
}

// GetStorageGroups implements api.Handler.
func (h *RestApiImplementation) GetStorageGroups(ctx context.Context) (*api.GetStorageGroupsResponse, error) {
	groups, err := h.storageGroupUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.StorageGroup, 0, len(groups))
	for _, group := range groups {
		items = append(items, convertGroupToDTO(group))
	}

	return &api.GetStorageGroupsResponse{
		Data: items,
	}, nil
}

// PatchStorageGroup implements api.Handler.
func (h *RestApiImplementation) PatchStorageGroup(ctx context.Context, req *api.PatchStorageGroupRequest, params api.PatchStorageGroupParams) (*api.PatchStorageGroupResponse, error) {
	updates := make(map[string]interface{})

	if req.Name.IsSet() {
		updates["name"] = req.Name.Value
	}
	if req.Alias.IsSet() {
		updates["alias"] = req.Alias.Value
	}
	if req.ParentId.IsSet() {
		updates["parent_id"] = req.ParentId.Value
	}

	group, err := h.storageGroupUseCase.Patch(ctx, params.ID, updates)
	if err != nil {
		return nil, err
	}

	return &api.PatchStorageGroupResponse{
		Data: []api.StorageGroup{convertGroupToDTO(group)},
	}, nil
}

// UpdateStorageGroup implements api.Handler.
func (h *RestApiImplementation) UpdateStorageGroup(ctx context.Context, req *api.UpdateStorageGroupRequest, params api.UpdateStorageGroupParams) (*api.UpdateStorageGroupResponse, error) {
	var parentID *uuid.UUID
	if req.ParentId.IsSet() {
		parentID = &req.ParentId.Value
	}
	group := &models.StorageGroup{
		ID:       params.ID,
		ParentID: parentID,
		Name:     req.Name,
		Alias:    req.Alias,
		UnitID:   req.UnitId,
	}

	updatedGroup, err := h.storageGroupUseCase.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	return &api.UpdateStorageGroupResponse{
		Data: []api.StorageGroup{convertGroupToDTO(updatedGroup)},
	}, nil
}

// CreateCell implements api.Handler.
func (h *RestApiImplementation) CreateCell(ctx context.Context, req *api.CreateCellRequest, params api.CreateCellParams) (*api.CreateCellResponse, error) {
	panic("unimplemented")
}

// CreateCellsGroup implements api.Handler.
func (h *RestApiImplementation) CreateCellsGroup(ctx context.Context, req *api.CreateCellsGroupRequest) (*api.CreateCellsGroupResponse, error) {
	panic("unimplemented")
}

// DeleteCell implements api.Handler.
func (h *RestApiImplementation) DeleteCell(ctx context.Context, params api.DeleteCellParams) error {
	panic("unimplemented")
}

// DeleteCellsGroup implements api.Handler.
func (h *RestApiImplementation) DeleteCellsGroup(ctx context.Context, params api.DeleteCellsGroupParams) error {
	panic("unimplemented")
}

// GetCellById implements api.Handler.
func (h *RestApiImplementation) GetCellById(ctx context.Context, params api.GetCellByIdParams) (*api.GetCellByIdResponse, error) {
	panic("unimplemented")
}

// GetCells implements api.Handler.
func (h *RestApiImplementation) GetCells(ctx context.Context, params api.GetCellsParams) (*api.GetCellsResponse, error) {
	panic("unimplemented")
}

// GetCellsGroupById implements api.Handler.
func (h *RestApiImplementation) GetCellsGroupById(ctx context.Context, params api.GetCellsGroupByIdParams) (*api.GetCellsGroupByIdResponse, error) {
	panic("unimplemented")
}

// GetCellsGroups implements api.Handler.
func (h *RestApiImplementation) GetCellsGroups(ctx context.Context) (*api.GetCellsGroupsResponse, error) {
	panic("unimplemented")
}

// PatchCell implements api.Handler.
func (h *RestApiImplementation) PatchCell(ctx context.Context, req *api.PatchCellRequest, params api.PatchCellParams) (*api.PatchCellResponse, error) {
	panic("unimplemented")
}

// PatchCellsGroup implements api.Handler.
func (h *RestApiImplementation) PatchCellsGroup(ctx context.Context, req *api.PatchCellsGroupRequest, params api.PatchCellsGroupParams) (*api.PatchCellsGroupResponse, error) {
	panic("unimplemented")
}

// UpdateCell implements api.Handler.
func (h *RestApiImplementation) UpdateCell(ctx context.Context, req *api.UpdateCellRequest, params api.UpdateCellParams) (*api.UpdateCellResponse, error) {
	panic("unimplemented")
}

// UpdateCellsGroup implements api.Handler.
func (h *RestApiImplementation) UpdateCellsGroup(ctx context.Context, req *api.UpdateCellsGroupRequest, params api.UpdateCellsGroupParams) (*api.UpdateCellsGroupResponse, error) {
	panic("unimplemented")
}
