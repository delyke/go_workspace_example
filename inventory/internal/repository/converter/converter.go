package converter

import (
	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

func ModelPartToRepo(part *model.Part) *repoModel.Part {
	return &repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ModelPartCategoryToRepo(part.Category),
		Dimensions:    ModelPartDimensionsToRepo(part.Dimensions),
		Manufacturer:  ModelPartManufacturerToRepo(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ModelPartMetadataToRepo(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ModelPartMetadataKindToRepo(kind model.PartMetadataKind) repoModel.PartMetadataKind {
	switch kind {
	case model.MetadataKindString:
		return repoModel.MetadataKindString
	case model.MetadataKindInt64:
		return repoModel.MetadataKindInt64
	case model.MetadataKindDouble:
		return repoModel.MetadataKindDouble
	case model.MetadataKindBool:
		return repoModel.MetadataKindBool
	default:
		return repoModel.MetadataKindUnknown
	}
}

func ModelPartMetadataToRepo(m map[string]*model.PartMetadataValue) map[string]*repoModel.PartMetadataValue {
	if len(m) == 0 {
		return nil
	}
	result := make(map[string]*repoModel.PartMetadataValue, len(m))
	for k, v := range m {
		result[k] = &repoModel.PartMetadataValue{
			Kind:   ModelPartMetadataKindToRepo(v.Kind),
			String: v.String,
			Int64:  v.Int64,
			Double: v.Double,
			Bool:   v.Bool,
		}
	}
	return result
}

func ModelPartManufacturerToRepo(manufacturer *model.PartManufacturer) *repoModel.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &repoModel.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ModelPartDimensionsToRepo(dimensions *model.PartDimensions) *repoModel.PartDimensions {
	if dimensions == nil {
		return nil
	}
	return &repoModel.PartDimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ModelPartCategoryToRepo(category model.PartCategory) repoModel.PartCategory {
	switch category {
	case model.PartCategoryEngine:
		return repoModel.PartCategoryEngine
	case model.PartCategoryFuel:
		return repoModel.PartCategoryFuel
	case model.PartCategoryPorthole:
		return repoModel.PartCategoryPorthole
	case model.PartCategoryWing:
		return repoModel.PartCategoryWing
	default:
		return repoModel.PartCategoryUnknown
	}
}

func RepoListPartsToModel(parts []*repoModel.Part) []*model.Part {
	modelParts := make([]*model.Part, len(parts))
	for i, part := range parts {
		modelParts[i] = lo.ToPtr(RepoPartToModel(*part))
	}
	return modelParts
}

func RepoPartToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      RepoPartCategoryToModel(part.Category),
		Dimensions:    RepoDimensionsToModel(part.Dimensions),
		Manufacturer:  RepoManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      RepoMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func RepoPartCategoryToModel(category repoModel.PartCategory) model.PartCategory {
	switch category {
	case repoModel.PartCategoryUnknown:
		return model.PartCategoryUnknown
	case repoModel.PartCategoryEngine:
		return model.PartCategoryEngine
	case repoModel.PartCategoryFuel:
		return model.PartCategoryFuel
	case repoModel.PartCategoryPorthole:
		return model.PartCategoryPorthole
	case repoModel.PartCategoryWing:
		return model.PartCategoryWing
	default:
		return model.PartCategoryUnknown
	}
}

func RepoDimensionsToModel(dimensions *repoModel.PartDimensions) *model.PartDimensions {
	if dimensions == nil {
		return nil
	}
	return &model.PartDimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func RepoManufacturerToModel(manufacturer *repoModel.PartManufacturer) *model.PartManufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.PartManufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func RepoMetadataToModel(metadata map[string]*repoModel.PartMetadataValue) map[string]*model.PartMetadataValue {
	if len(metadata) == 0 {
		return nil
	}
	result := make(map[string]*model.PartMetadataValue, len(metadata))
	for k, v := range metadata {
		result[k] = &model.PartMetadataValue{
			Kind:   RepoMetadataKindToModel(v.Kind),
			String: v.String,
			Int64:  v.Int64,
			Double: v.Double,
			Bool:   v.Bool,
		}
	}
	return result
}

func RepoMetadataKindToModel(kind repoModel.PartMetadataKind) model.PartMetadataKind {
	switch kind {
	case repoModel.MetadataKindUnknown:
		return model.MetadataKindUnknown
	case repoModel.MetadataKindString:
		return model.MetadataKindString
	case repoModel.MetadataKindInt64:
		return model.MetadataKindInt64
	case repoModel.MetadataKindDouble:
		return model.MetadataKindDouble
	case repoModel.MetadataKindBool:
		return model.MetadataKindBool
	default:
		return model.MetadataKindUnknown
	}
}
