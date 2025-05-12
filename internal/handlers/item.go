package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

func convertCellPathToOptionalDTO(cellPath *[]models.CellPathSegment) []api.CellForInstanceOptionalCellPathItem {
	if cellPath == nil {
		return nil
	}

	dtoCellPath := make([]api.CellForInstanceOptionalCellPathItem, 0, len(*cellPath))
	for _, pathSegment := range *cellPath {
		dtoCellPath = append(dtoCellPath, api.CellForInstanceOptionalCellPathItem{
			ID:         pathSegment.ID,
			Alias:      pathSegment.Alias,
			Name:       pathSegment.Name,
			ObjectType: api.CellForInstanceOptionalCellPathItemObjectType(pathSegment.ObjectType),
		})
	}
	return dtoCellPath
}

func convertCellPathToDTO(cellPath *[]models.CellPathSegment) []api.CellForInstanceCellPathItem {
	if cellPath == nil {
		return nil
	}

	dtoCellPath := make([]api.CellForInstanceCellPathItem, 0, len(*cellPath))
	for _, pathSegment := range *cellPath {
		dtoCellPath = append(dtoCellPath, api.CellForInstanceCellPathItem{
			ID:         pathSegment.ID,
			Alias:      pathSegment.Alias,
			Name:       pathSegment.Name,
			ObjectType: api.CellForInstanceCellPathItemObjectType(pathSegment.ObjectType),
		})
	}
	return dtoCellPath
}

func convertCellToNilDTO(cell *models.Cell) api.CellForInstance {
	res := api.CellForInstance{}
	if cell == nil {
		return res
	}
	modelCell := convertCellToDTO(cell)
	res = modelCell
	return res
}

func convertCellOptionalToNilDTO(cell *models.Cell) api.NilCellForInstanceOptional {
	res := api.NilCellForInstanceOptional{}
	if cell == nil {
		res.SetToNull()
		return res
	}
	modelCell := api.CellForInstanceOptional{
		ID:       cell.ID,
		Alias:    cell.Alias,
		Row:      cell.Row,
		Level:    cell.Level,
		Position: cell.Position,
		CellPath: convertCellPathToOptionalDTO(cell.Path),
	}
	res.SetTo(modelCell)
	return res
}

func convertCellToDTO(cell *models.Cell) api.CellForInstance {
	return api.CellForInstance{
		ID:       cell.ID,
		Alias:    cell.Alias,
		Row:      cell.Row,
		Level:    cell.Level,
		Position: cell.Position,
		CellPath: convertCellPathToDTO(cell.Path),
	}
}

func convertItemInstanceToDTO(itemInstance *models.ItemInstance) api.InstanceForItem {
	return api.InstanceForItem{
		ID:      itemInstance.ID,
		Status:  api.InstanceForItemStatus(itemInstance.Status),
		Variant: convertItemVariantToDTO(itemInstance.Variant),
		Cell:    convertCellOptionalToNilDTO(itemInstance.Cell),
	}
}

func convertItemInstanceToTaskItemDTO(itemInstance *models.ItemInstance) api.InstanceFull {
	if itemInstance == nil {
		return api.InstanceFull{}
	}

	var item api.ItemForList
	if itemInstance.Item != nil {
		var description api.NilString
		PtrToApiNil(itemInstance.Item.Description, &description)
		item = api.ItemForList{
			ID:          itemInstance.Item.ID,
			Name:        itemInstance.Item.Name,
			Description: description,
			Variants:    convertItemVariantsToDTO(itemInstance.Item.Variants),
		}
	}
	return api.InstanceFull{
		ID:      itemInstance.ID,
		Status:  api.InstanceFullStatus(itemInstance.Status),
		Variant: convertItemVariantToDTO(itemInstance.Variant),
		Cell:    convertCellOptionalToNilDTO(itemInstance.Cell),
		Item:    item,
	}
}

func convertItemVariantToDTO(variant *models.ItemVariant) api.ItemVariant {
	if variant == nil {
		return api.ItemVariant{}
	}

	var article api.NilString
	PtrToApiNil(variant.Article, &article)

	var ean13 api.NilInt64
	PtrToApiNil(variant.EAN13, &ean13)

	return api.ItemVariant{
		ID:      variant.ID,
		Name:    variant.Name,
		Article: article,
		Ean13:   ean13,
	}
}

func convertItemVariantsToDTO(variants []*models.ItemVariant) []api.ItemVariant {
	if variants == nil {
		return []api.ItemVariant{}
	}
	dtoVariants := make([]api.ItemVariant, 0, len(variants))
	for _, variant := range variants {
		if variant == nil {
			continue
		}
		dtoVariants = append(dtoVariants, convertItemVariantToDTO(variant))
	}
	return dtoVariants
}

func convertItemInstancesForItemToDTO(itemInstances []*models.ItemInstance) []api.InstanceForItem {
	if itemInstances == nil {
		return []api.InstanceForItem{}
	}

	dtoInstances := make([]api.InstanceForItem, 0, len(itemInstances))
	for _, instance := range itemInstances {
		if instance == nil {
			continue
		}

		var cellPath []api.CellForInstanceCellPathItem
		if instance.Cell != nil && instance.Cell.Path != nil {
			cellPath = make([]api.CellForInstanceCellPathItem, 0, len(*instance.Cell.Path))
			for _, pathSegment := range *instance.Cell.Path {
				cellPath = append(cellPath, api.CellForInstanceCellPathItem{
					ID:         pathSegment.ID,
					Alias:      pathSegment.Alias,
					Name:       pathSegment.Name,
					ObjectType: api.CellForInstanceCellPathItemObjectType(pathSegment.ObjectType),
				})
			}
		}

		var article api.NilString
		if instance.Variant != nil {
			PtrToApiNil(instance.Variant.Article, &article)
		}

		var ean13 api.NilInt64
		if instance.Variant != nil {
			PtrToApiNil(instance.Variant.EAN13, &ean13)
		}

		var variant api.ItemVariant
		if instance.Variant != nil {
			variant = api.ItemVariant{
				ID:      instance.Variant.ID,
				Name:    instance.Variant.Name,
				Article: article,
				Ean13:   ean13,
			}
		}

		dtoInstances = append(dtoInstances, api.InstanceForItem{
			ID:      instance.ID,
			Status:  api.InstanceForItemStatus(instance.Status),
			Variant: variant,
			Cell:    convertCellOptionalToNilDTO(instance.Cell),
		})
	}
	return dtoInstances
}

func convertItemToFullDTO(item *models.Item, itemInstances []*models.ItemInstance) api.ItemFull {
	variants := make([]api.ItemVariant, 0, len(item.Variants))
	if item.Variants != nil {
		for _, variant := range item.Variants {
			variants = append(variants, convertItemVariantToDTO(variant))
		}
	}

	var description api.NilString
	PtrToApiNil(item.Description, &description)

	return api.ItemFull{
		ID:          item.ID,
		Name:        item.Name,
		Description: description,
		Variants:    variants,
		Items:       convertItemInstancesForItemToDTO(itemInstances),
	}
}

func (h *RestApiImplementation) CreateItem(ctx context.Context, req *api.CreateItemRequest) (api.CreateItemRes, error) {
	var description *string
	if val, ok := req.Description.Get(); ok {
		description = &val
	}

	item := &models.Item{
		Name:        req.Name,
		Description: description,
	}

	createdItem, err := h.itemUseCase.CreateItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemResponse{
		Data: convertItemToFullDTO(createdItem, nil),
	}, nil
}

func (h *RestApiImplementation) GetItemById(ctx context.Context, params api.GetItemByIdParams) (api.GetItemByIdRes, error) {
	item, err := h.itemUseCase.GetItemById(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetItemByIdResponse{
		Data: convertItemToFullDTO(item, item.Instances),
	}, nil
}

func (h *RestApiImplementation) GetItems(ctx context.Context) (api.GetItemsRes, error) {
	items, err := h.itemUseCase.GetItemsAll(ctx)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]api.ItemForList, 0, len(items))
	for _, item := range items {
		variants := make([]api.ItemVariant, 0, len(item.Variants))
		if item.Variants != nil {
			for _, variant := range item.Variants {
				variants = append(variants, convertItemVariantToDTO(variant))
			}
		}

		var description api.NilString
		PtrToApiNil(item.Description, &description)

		dtoItems = append(dtoItems, api.ItemForList{
			ID:          item.ID,
			Name:        item.Name,
			Description: description,
			Variants:    variants,
		})
	}

	return &api.GetItemsResponse{
		Data: dtoItems,
	}, nil
}

func (h *RestApiImplementation) DeleteItem(ctx context.Context, params api.DeleteItemParams) (api.DeleteItemRes, error) {
	err := h.itemUseCase.DeleteItem(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.DeleteItemNoContent{}, nil
}

func (h *RestApiImplementation) UpdateItem(ctx context.Context, req *api.UpdateItemRequest, params api.UpdateItemParams) (api.UpdateItemRes, error) {

	newItem := &models.Item{
		ID:          params.ID,
		Name:        req.Name,
		Description: ApiValueToPtr(req.Description),
	}

	updatedItem, err := h.itemUseCase.UpdateItem(ctx, newItem)
	if err != nil {
		return nil, err
	}

	return &api.UpdateItemResponse{
		Data: convertItemToFullDTO(updatedItem, nil),
	}, nil
}

// // PatchItem implements api.Handler.
// func (h *RestApiImplementation) PatchItem(ctx context.Context, req *api.PatchItemRequest, params api.PatchItemParams) (*api.PatchItemResponse, error) {
// 	updates := make(map[string]interface{})

// 	if val, ok := req.Name.Get(); ok {
// 		updates["name"] = val
// 	}
// 	if val, ok := req.Description.Get(); ok {
// 		updates["description"] = &val
// 	}

// 	// Handle variants updates
// 	if req.Variants != nil {
// 		variants := make([]interface{}, len(req.Variants))
// 		for i, v := range req.Variants {
// 			variant := make(map[string]interface{})
// 			variant["id"] = v.ID
// 			variant["name"] = v.Name
// 			if v.Article.Set {
// 				variant["article"] = &v.Article.Value
// 			}
// 			if v.Ean13.Set {
// 				variant["ean13"] = float64(v.Ean13.Value)
// 			}
// 			variants[i] = variant
// 		}
// 		updates["variants"] = variants
// 	}

// 	orgID, err := usecases.GetOrganizationIDFromContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	item, err := h.itemUseCase.Patch(ctx, orgID, params.ID, updates)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &api.PatchItemResponse{
// 		Data: *convertItemToFullDTO(item, nil),
// 	}, nil
// }

func (h *RestApiImplementation) CreateItemVariant(ctx context.Context, req *api.CreateItemVariantRequest, params api.CreateItemVariantParams) (api.CreateItemVariantRes, error) {
	variant := &models.ItemVariant{
		Name:    req.Name,
		ItemID:  params.ID,
		Article: ApiValueToPtr(req.Article),
		EAN13:   ApiValueToPtr(req.Ean13),
	}

	res, err := h.itemUseCase.CreateItemVariant(ctx, variant)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemVariantResponse{
		Data: convertItemVariantToDTO(res),
	}, nil
}

func (h *RestApiImplementation) GetItemVariants(ctx context.Context, params api.GetItemVariantsParams) (api.GetItemVariantsRes, error) {
	variants, err := h.itemUseCase.GetItemVariants(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetItemVariantsResponse{
		Data: convertItemVariantsToDTO(variants),
	}, nil
}

func (h *RestApiImplementation) GetItemVariantById(ctx context.Context, params api.GetItemVariantByIdParams) (api.GetItemVariantByIdRes, error) {
	variant, err := h.itemUseCase.GetItemVariantById(ctx, params.ID, params.VariantId)
	if err != nil {
		return nil, err
	}

	return &api.GetItemVariantByIdResponse{
		Data: convertItemVariantToDTO(variant),
	}, nil
}

func (h *RestApiImplementation) UpdateItemVariant(ctx context.Context, req *api.UpdateItemVariantRequest, params api.UpdateItemVariantParams) (api.UpdateItemVariantRes, error) {
	variant := &models.ItemVariant{
		ID:      params.VariantId,
		Name:    req.Name,
		ItemID:  params.ID,
		Article: ApiValueToPtr(req.Article),
		EAN13:   ApiValueToPtr(req.Ean13),
	}

	updatedVariant, err := h.itemUseCase.UpdateItemVariant(ctx, variant)
	if err != nil {
		return nil, err
	}

	return &api.UpdateItemVariantResponse{
		Data: convertItemVariantToDTO(updatedVariant),
	}, nil
}
func (h *RestApiImplementation) DeleteItemVariant(ctx context.Context, params api.DeleteItemVariantParams) (api.DeleteItemVariantRes, error) {
	err := h.itemUseCase.DeleteItemVariant(ctx, params.ID, params.VariantId)
	if err != nil {
		return nil, err
	}
	return &api.DeleteItemVariantNoContent{}, nil
}

// instances

// CreateInstanceForItem implements api.Handler.
func (h *RestApiImplementation) CreateInstanceForItem(ctx context.Context, req *api.CreateInstanceForItemRequest, params api.CreateInstanceForItemParams) (api.CreateInstanceForItemRes, error) {
	itemInstance := &models.ItemInstance{
		ItemID:    params.ItemId,
		VariantID: req.VariantId,
		CellID:    &req.CellId,
	}

	itemInstance, err := h.itemUseCase.CreateItemInstance(ctx, itemInstance)
	if err != nil {
		return nil, err
	}

	return &api.CreateInstanceForItemResponse{
		Data: convertItemInstanceToDTO(itemInstance),
	}, nil

}

func (h *RestApiImplementation) DeleteInstanceById(ctx context.Context, params api.DeleteInstanceByIdParams) (api.DeleteInstanceByIdRes, error) {
	err := h.itemUseCase.DeleteItemInstance(ctx, params.InstanceId)
	if err != nil {
		return nil, err
	}
	return &api.DeleteInstanceByIdOK{}, nil
}

// GetInstanceById implements api.Handler.
func (h *RestApiImplementation) GetInstanceById(ctx context.Context, params api.GetInstanceByIdParams) (api.GetInstanceByIdRes, error) {
	instance, err := h.itemUseCase.GetItemInstanceById(ctx, params.InstanceId)
	if err != nil {
		return nil, err
	}
	return &api.GetInstanceByIdResponse{
		Data: convertItemInstanceToTaskItemDTO(instance),
	}, nil
}

func (h *RestApiImplementation) UpdateInstanceById(ctx context.Context, req *api.UpdateInstanceRequest, params api.UpdateInstanceByIdParams) (api.UpdateInstanceByIdRes, error) {
	instance := &models.ItemInstance{
		ID:        params.InstanceId,
		VariantID: req.VariantId,
		CellID:    &req.CellId,
	}

	updatedInstance, err := h.itemUseCase.UpdateItemInstance(ctx, instance)
	if err != nil {
		return nil, err
	}

	return &api.UpdateInstanceResponse{
		Data: convertItemInstanceToTaskItemDTO(updatedInstance),
	}, nil
}

func (h *RestApiImplementation) GetInstances(ctx context.Context) (api.GetInstancesRes, error) {
	instances, err := h.itemUseCase.GetItemInstancesAll(ctx)
	if err != nil {
		return nil, err
	}

	dtoInstances := make([]api.InstanceFull, 0, len(instances))
	for _, instance := range instances {
		dtoInstances = append(dtoInstances, convertItemInstanceToTaskItemDTO(instance))
	}

	return &api.GetInstancesResponse{
		Data: dtoInstances,
	}, nil
}

// GetInstancesByItemId implements api.Handler.
func (h *RestApiImplementation) GetInstancesByItemId(ctx context.Context, params api.GetInstancesByItemIdParams) (api.GetInstancesByItemIdRes, error) {
	panic("unimplemented")
}
