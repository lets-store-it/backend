package handlers

import (
	"context"

	"github.com/let-store-it/backend/generated/api"
	"github.com/let-store-it/backend/internal/models"
)

func convertItemVariantToDTO(variant *models.ItemVariant) api.ItemVariant {
	var article api.NilString
	PtrToApiNil(variant.Article, &article)

	var ean13 api.NilInt32
	PtrToApiNil(variant.EAN13, &ean13)

	return api.ItemVariant{
		ID:      variant.ID,
		Name:    variant.Name,
		Article: article,
		Ean13:   ean13,
	}
}

func convertItemVariantsToDTO(variants []*models.ItemVariant) []api.ItemVariant {
	dtoVariants := make([]api.ItemVariant, 0, len(variants))
	for _, variant := range variants {
		dtoVariants = append(dtoVariants, convertItemVariantToDTO(variant))
	}
	return dtoVariants
}

func convertItemInstancesForItemToDTO(itemInstances *[]models.ItemInstance) []api.InstanceForItem {
	if itemInstances == nil {
		return nil
	}

	dtoInstances := make([]api.InstanceForItem, 0, len(*itemInstances))
	for _, instance := range *itemInstances {

		var cellPath []api.CellForInstanceCellPathItem
		for _, pathSegment := range *instance.Cell.Path {
			cellPath = append(cellPath, api.CellForInstanceCellPathItem{
				ID:         pathSegment.ID,
				Alias:      pathSegment.Alias,
				Name:       pathSegment.Name,
				ObjectType: api.CellForInstanceCellPathItemObjectType(pathSegment.ObjectType),
			})
		}

		var article api.NilString
		PtrToApiNil(instance.Variant.Article, &article)

		var ean13 api.NilInt32
		PtrToApiNil(instance.Variant.EAN13, &ean13)

		dtoInstances = append(dtoInstances, api.InstanceForItem{
			ID:     instance.ID,
			Status: api.InstanceForItemStatus(instance.Status),
			Variant: api.ItemVariant{
				ID:      instance.Variant.ID,
				Name:    instance.Variant.Name,
				Article: article,
				Ean13:   ean13,
			},
			Cell: api.CellForInstance{
				ID:       instance.Cell.ID,
				Alias:    instance.Cell.Alias,
				Row:      instance.Cell.Row,
				Level:    instance.Cell.Level,
				Position: instance.Cell.Position,
				CellPath: cellPath,
			},
		})
	}
	return dtoInstances
}
func convertItemToFullDTO(item *models.Item, itemInstances *[]models.ItemInstance) api.ItemFull {
	var variants []api.ItemVariant
	if item.Variants != nil {
		for _, variant := range *item.Variants {
			variants = append(variants, convertItemVariantToDTO(&variant))
		}
	}

	var description api.NilString
	PtrToApiNil(item.Description, &description)

	return api.ItemFull{
		ID:          item.ID,
		Name:        item.Name,
		Description: description,
		Variants:    variants,
		Instances:   convertItemInstancesForItemToDTO(itemInstances),
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

	createdItem, err := h.itemUseCase.Create(ctx, item)
	if err != nil {
		return nil, err
	}

	return &api.CreateItemResponse{
		Data: convertItemToFullDTO(createdItem, nil),
	}, nil
}

func (h *RestApiImplementation) GetItemById(ctx context.Context, params api.GetItemByIdParams) (api.GetItemByIdRes, error) {
	item, err := h.itemUseCase.GetByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &api.GetItemByIdResponse{
		Data: convertItemToFullDTO(item, item.Instances),
	}, nil
}

func (h *RestApiImplementation) GetItems(ctx context.Context) (api.GetItemsRes, error) {
	items, err := h.itemUseCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtoItems := make([]api.ItemForList, 0, len(items))
	for _, item := range items {
		var variants []api.ItemVariant
		if item.Variants != nil {
			for _, variant := range *item.Variants {
				variants = append(variants, convertItemVariantToDTO(&variant))
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
	err := h.itemUseCase.Delete(ctx, params.ID)
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

	updatedItem, err := h.itemUseCase.Update(ctx, newItem)
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
