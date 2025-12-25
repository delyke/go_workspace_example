package converter

import (
	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/order/internal/model"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func ProtoPartCategoryToModel(c inventoryV1.Category) model.PartCategory {
	switch c {
	case inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED:
		return model.PartCategoryUnknown
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.PartCategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.PartCategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.PartCategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.PartCategoryWing
	default:
		return model.PartCategoryUnknown
	}
}

func ProtoPartDimensionsToModel(d *inventoryV1.Dimensions) *model.PartDimensions {
	if d == nil {
		return nil
	}
	return &model.PartDimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ProtoPartManufacturerToModel(m *inventoryV1.Manufacturer) *model.PartManufacturer {
	if m == nil {
		return nil
	}
	return &model.PartManufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ProtoMetadataValueToModel(v *inventoryV1.MetadataValue) *model.PartMetadataValue {
	if v == nil {
		return nil
	}

	switch k := v.Kind.(type) {
	case *inventoryV1.MetadataValue_StringValue:
		return model.MetaString(k.StringValue)
	case *inventoryV1.MetadataValue_Int_64Value:
		return model.MetaInt64(k.Int_64Value)
	case *inventoryV1.MetadataValue_DoubleValue:
		return model.MetaDouble(k.DoubleValue)
	case *inventoryV1.MetadataValue_BoolValue:
		if k.BoolValue == nil {
			return &model.PartMetadataValue{Kind: model.MetadataKindBool}
		}
		return model.MetaBool(k.BoolValue.Value)
	default:
		return &model.PartMetadataValue{Kind: model.MetadataKindUnknown}
	}
}

func ProtoMetadataToModel(m map[string]*inventoryV1.MetadataValue) map[string]*model.PartMetadataValue {
	if m == nil {
		return nil
	}

	out := make(map[string]*model.PartMetadataValue, len(m))
	for key, val := range m {
		out[key] = ProtoMetadataValueToModel(val)
	}
	return out
}

func ProtoPartsToModel(part []*inventoryV1.Part) []model.Part {
	var parts []model.Part
	for _, p := range part {
		parts = append(parts, model.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      ProtoPartCategoryToModel(p.Category),
			Dimensions:    ProtoPartDimensionsToModel(p.Dimensions),
			Manufacturer:  ProtoPartManufacturerToModel(p.Manufacturer),
			Tags:          p.Tags,
			Metadata:      ProtoMetadataToModel(p.Metadata),
			CreatedAt:     lo.ToPtr(p.CreatedAt.AsTime()),
			UpdatedAt:     lo.ToPtr(p.UpdatedAt.AsTime()),
		})
	}
	return parts
}

func PartCategoryToProto(c model.PartCategory) inventoryV1.Category {
	switch c {
	case model.PartCategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.PartCategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.PartCategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.PartCategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
	}
}

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))

	for _, category := range filter.Categories {
		categories = append(categories, PartCategoryToProto(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
