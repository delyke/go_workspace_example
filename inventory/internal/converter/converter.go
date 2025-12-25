package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func ModelPartListToProto(parts []*model.Part) []*inventoryV1.Part {
	protoParts := make([]*inventoryV1.Part, len(parts))
	for i, part := range parts {
		protoParts[i] = ModelPartToProto(part)
	}
	return protoParts
}

func ProtoFiltersToModel(f *inventoryV1.PartsFilter) *model.PartsFilter {
	return &model.PartsFilter{
		UUIDs:                 f.GetUuids(),
		Tags:                  f.GetTags(),
		ManufacturerCountries: f.GetManufacturerCountries(),
		Categories:            ProtoCategoriesToModel(f.GetCategories()),
		Names:                 f.GetNames(),
	}
}

func ProtoCategoriesToModel(categories []inventoryV1.Category) []model.PartCategory {
	mc := make([]model.PartCategory, len(categories))
	for i, cat := range categories {
		mc[i] = ProtoCategoryToModel(cat)
	}
	return mc
}

func ProtoCategoryToModel(c inventoryV1.Category) model.PartCategory {
	switch c {
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

func ModelPartToProto(part *model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ModelPartCategoryToProto(part.Category),
		Dimensions:    ModelPartDimensionsToProto(part.Dimensions),
		Manufacturer:  ModelPartManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ModelPartMetadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(*part.CreatedAt),
		UpdatedAt:     timestamppb.New(*part.UpdatedAt),
	}
}

func ModelPartMetadataToProto(m map[string]*model.PartMetadataValue) map[string]*inventoryV1.MetadataValue {
	if m == nil {
		return nil
	}
	res := make(map[string]*inventoryV1.MetadataValue, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		switch v.Kind {
		case model.MetadataKindString:
			if v.String == nil {
				continue
			}
			res[k] = &inventoryV1.MetadataValue{
				Kind: &inventoryV1.MetadataValue_StringValue{
					StringValue: *v.String,
				},
			}
		case model.MetadataKindInt64:
			if v.Int64 == nil {
				continue
			}
			res[k] = &inventoryV1.MetadataValue{
				Kind: &inventoryV1.MetadataValue_Int_64Value{
					Int_64Value: *v.Int64,
				},
			}
		case model.MetadataKindDouble:
			if v.Double == nil {
				continue
			}
			res[k] = &inventoryV1.MetadataValue{
				Kind: &inventoryV1.MetadataValue_DoubleValue{
					DoubleValue: *v.Double,
				},
			}
		case model.MetadataKindBool:
			if v.Bool == nil {
				continue
			}
			res[k] = &inventoryV1.MetadataValue{
				Kind: &inventoryV1.MetadataValue_BoolValue{
					BoolValue: wrapperspb.Bool(*v.Bool),
				},
			}
		default:
			continue
		}
	}
	return res
}

func ModelPartManufacturerToProto(m *model.PartManufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ModelPartDimensionsToProto(d *model.PartDimensions) *inventoryV1.Dimensions {
	if d == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ModelPartCategoryToProto(category model.PartCategory) inventoryV1.Category {
	switch category {
	case model.PartCategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.PartCategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.PartCategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.PartCategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNKNOWN_UNSPECIFIED
	}
}
