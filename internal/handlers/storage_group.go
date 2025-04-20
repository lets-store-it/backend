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
	group, err := h.storageGroupUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	if req.Name.IsSet() {
		group.Name = req.Name.Value
	}
	if req.Alias.IsSet() {
		group.Alias = req.Alias.Value
	}
	if req.ParentId.IsSet() {
		group.ParentID = &req.ParentId.Value
	}

	updatedGroup, err := h.storageGroupUseCase.Update(ctx, group)
	if err != nil {
		return nil, err
	}

	return &api.PatchStorageGroupResponse{
		Data: []api.StorageGroup{convertGroupToDTO(updatedGroup)},
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

func toCellsGroupDTO(group *models.CellsGroup) *api.CellGroupBase {
	return &api.CellGroupBase{
		ID:             group.ID,
		Name:           group.Name,
		Alias:          group.Alias,
		StorageGroupID: group.StorageGroupID,
	}
}

// GetCellsGroups implements api.Handler.
func (h *RestApiImplementation) GetCellsGroups(ctx context.Context) (*api.GetCellsGroupsResponse, error) {
	cellsGroups, err := h.storageGroupUseCase.GetCellsGroups(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]api.CellGroupBase, 0, len(cellsGroups))
	for _, group := range cellsGroups {
		items = append(items, *toCellsGroupDTO(group))
	}

	return &api.GetCellsGroupsResponse{
		Data: items,
	}, nil
}

// CreateCellsGroup implements api.Handler.
func (h *RestApiImplementation) CreateCellsGroup(ctx context.Context, req *api.CreateCellsGroupRequest) (*api.CreateCellsGroupResponse, error) {
	cellGroup, err := h.storageGroupUseCase.CreateCellsGroup(ctx, req.StorageGroupID, req.Name, req.Alias)
	if err != nil {
		return nil, err
	}

	return &api.CreateCellsGroupResponse{
		Data: *toCellsGroupDTO(cellGroup),
	}, nil
}

// GetCellsGroupById implements api.Handler.
func (h *RestApiImplementation) GetCellsGroupById(ctx context.Context, params api.GetCellsGroupByIdParams) (*api.GetCellsGroupByIdResponse, error) {
	cellGroup, err := h.storageGroupUseCase.GetCellsGroupByID(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	return &api.GetCellsGroupByIdResponse{
		Data: *toCellsGroupDTO(cellGroup),
	}, nil
}

// DeleteCellsGroup implements api.Handler.
func (h *RestApiImplementation) DeleteCellsGroup(ctx context.Context, params api.DeleteCellsGroupParams) error {
	return h.storageGroupUseCase.DeleteCellsGroup(ctx, params.GroupId)
}

// UpdateCell implements api.Handler.
func (h *RestApiImplementation) UpdateCellsGroup(ctx context.Context, req *api.UpdateCellsGroupRequest, params api.UpdateCellsGroupParams) (*api.UpdateCellsGroupResponse, error) {
	model := &models.CellsGroup{
		ID:    params.GroupId,
		Name:  req.Name,
		Alias: req.Alias,
	}
	cellGroup, err := h.storageGroupUseCase.UpdateCellsGroup(ctx, model)
	if err != nil {
		return nil, err
	}

	return &api.UpdateCellsGroupResponse{
		Data: *toCellsGroupDTO(cellGroup),
	}, nil
}

// PatchCellsGroup implements api.Handler.
func (h *RestApiImplementation) PatchCellsGroup(ctx context.Context, req *api.PatchCellsGroupRequest, params api.PatchCellsGroupParams) (*api.PatchCellsGroupResponse, error) {
	group, err := h.storageGroupUseCase.GetCellsGroupByID(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	if req.Name.IsSet() {
		group.Name = req.Name.Value
	}
	if req.Alias.IsSet() {
		group.Alias = req.Alias.Value
	}

	updatedGroup, err := h.storageGroupUseCase.UpdateCellsGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return &api.PatchCellsGroupResponse{
		Data: *toCellsGroupDTO(updatedGroup),
	}, nil
}

// CreateCell implements api.Handler.
func (h *RestApiImplementation) CreateCell(ctx context.Context, req *api.CreateCellRequest, params api.CreateCellParams) (*api.CreateCellResponse, error) {
	cell, err := h.storageGroupUseCase.CreateCell(ctx, params.GroupId, req.Alias, req.Row, req.Level, req.Position)
	if err != nil {
		return nil, err
	}

	return &api.CreateCellResponse{
		Data: api.CellBase{
			ID:       cell.ID,
			Alias:    cell.Alias,
			Row:      cell.Row,
			Level:    cell.Level,
			Position: cell.Position,
		},
	}, nil
}

// DeleteCell implements api.Handler.
func (h *RestApiImplementation) DeleteCell(ctx context.Context, params api.DeleteCellParams) error {
	return h.storageGroupUseCase.DeleteCell(ctx, params.GroupId, params.CellId)
}

// GetCellById implements api.Handler.
func (h *RestApiImplementation) GetCellById(ctx context.Context, params api.GetCellByIdParams) (*api.GetCellByIdResponse, error) {
	cell, err := h.storageGroupUseCase.GetCellByID(ctx, params.GroupId, params.CellId)
	if err != nil {
		return nil, err
	}

	return &api.GetCellByIdResponse{
		ID:       cell.ID,
		Alias:    cell.Alias,
		Row:      cell.Row,
		Level:    cell.Level,
		Position: cell.Position,
	}, nil
}

// GetCells implements api.Handler.
func (h *RestApiImplementation) GetCells(ctx context.Context, params api.GetCellsParams) (*api.GetCellsResponse, error) {
	cells, err := h.storageGroupUseCase.GetCells(ctx, params.GroupId)
	if err != nil {
		return nil, err
	}

	items := make([]api.CellBase, 0, len(cells))
	for _, cell := range cells {
		items = append(items, api.CellBase{
			ID:       cell.ID,
			Alias:    cell.Alias,
			Row:      cell.Row,
			Level:    cell.Level,
			Position: cell.Position,
		})
	}

	return &api.GetCellsResponse{
		Data: items,
	}, nil
}

// PatchCell implements api.Handler.
func (h *RestApiImplementation) PatchCell(ctx context.Context, req *api.PatchCellRequest, params api.PatchCellParams) (*api.PatchCellResponse, error) {
	cell, err := h.storageGroupUseCase.GetCellByID(ctx, params.GroupId, params.CellId)
	if err != nil {
		return nil, err
	}

	if req.Alias.IsSet() {
		cell.Alias = req.Alias.Value
	}
	if req.Row.IsSet() {
		cell.Row = req.Row.Value
	}
	if req.Level.IsSet() {
		cell.Level = req.Level.Value
	}
	if req.Position.IsSet() {
		cell.Position = req.Position.Value
	}

	updatedCell, err := h.storageGroupUseCase.UpdateCell(ctx, params.GroupId, cell)
	if err != nil {
		return nil, err
	}

	return &api.PatchCellResponse{
		Data: api.CellBase{
			ID:       updatedCell.ID,
			Alias:    updatedCell.Alias,
			Row:      updatedCell.Row,
			Level:    updatedCell.Level,
			Position: updatedCell.Position,
		},
	}, nil
}

// UpdateCell implements api.Handler.
func (h *RestApiImplementation) UpdateCell(ctx context.Context, req *api.UpdateCellRequest, params api.UpdateCellParams) (*api.UpdateCellResponse, error) {
	cell := &models.Cell{
		ID:       params.CellId,
		Alias:    req.Alias,
		Row:      req.Row,
		Level:    req.Level,
		Position: req.Position,
	}

	updatedCell, err := h.storageGroupUseCase.UpdateCell(ctx, params.GroupId, cell)
	if err != nil {
		return nil, err
	}

	return &api.UpdateCellResponse{
		Data: api.CellBase{
			ID:       updatedCell.ID,
			Alias:    updatedCell.Alias,
			Row:      updatedCell.Row,
			Level:    updatedCell.Level,
			Position: updatedCell.Position,
		},
	}, nil
}
