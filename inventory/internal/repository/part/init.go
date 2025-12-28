package part

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

func (r *repository) Init(ctx context.Context) error {
	now := time.Now()
	partList := []*repoModel.Part{
		{
			Uuid:          "11111111-1111-1111-1111-111111111111",
			Name:          "Main Engine X1",
			Description:   "Primary propulsion engine",
			Price:         1200000,
			StockQuantity: 5,
			Category:      repoModel.PartCategoryEngine,
			Dimensions: &repoModel.PartDimensions{
				Length: 4.5,
				Width:  2.1,
				Height: 2.0,
				Weight: 1800,
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    "Orbital Dynamics",
				Country: "Germany",
				Website: "https://orbital-dynamics.de",
			},
			Tags: []string{"main", "engine", "booster"},
			Metadata: map[string]any{
				"fuel_type": "liquid",
				"reusable":  true,
			},
			CreatedAt: lo.ToPtr(now),
			UpdatedAt: lo.ToPtr(now),
		},
		{
			Uuid:          "22222222-2222-2222-2222-222222222222",
			Name:          "Fuel Tank A9",
			Description:   "High-pressure fuel container",
			Price:         320000,
			StockQuantity: 12,
			Category:      repoModel.PartCategoryFuel,
			Dimensions: &repoModel.PartDimensions{
				Length: 3.0,
				Width:  1.8,
				Height: 1.8,
				Weight: 900,
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    "CosmoFuel",
				Country: "USA",
				Website: "https://cosmofuel.com",
			},
			Tags: []string{"fuel", "tank"},
			Metadata: map[string]any{
				"capacity_liters": int64(5000),
			},
			CreatedAt: lo.ToPtr(now),
			UpdatedAt: lo.ToPtr(now),
		},
		{
			Uuid:          "33333333-3333-3333-3333-333333333333",
			Name:          "Observation Porthole",
			Description:   "Reinforced glass porthole",
			Price:         85000,
			StockQuantity: 20,
			Category:      repoModel.PartCategoryPorthole,
			Dimensions: &repoModel.PartDimensions{
				Length: 1.2,
				Width:  1.2,
				Height: 0.2,
				Weight: 80,
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    "SpaceGlass",
				Country: "France",
				Website: "https://spaceglass.fr",
			},
			Tags: []string{"window", "glass"},
			Metadata: map[string]any{
				"radiation_protected": true,
			},
			CreatedAt: lo.ToPtr(now),
			UpdatedAt: lo.ToPtr(now),
		},
		{
			Uuid:          "44444444-4444-4444-4444-444444444444",
			Name:          "Wing Module L",
			Description:   "Left aerodynamic wing",
			Price:         210000,
			StockQuantity: 7,
			Category:      repoModel.PartCategoryWing,
			Dimensions: &repoModel.PartDimensions{
				Length: 6.0,
				Width:  2.5,
				Height: 0.8,
				Weight: 600,
			},
			Manufacturer: &repoModel.PartManufacturer{
				Name:    "AeroSpace Ltd",
				Country: "UK",
				Website: "https://aerospace.co.uk",
			},
			Tags: []string{"wing", "left"},
			Metadata: map[string]any{
				"material": "carbon",
			},
			CreatedAt: lo.ToPtr(now),
			UpdatedAt: lo.ToPtr(now),
		},
		{
			Uuid:          "55555555-5555-5555-5555-555555555555",
			Name:          "Wing Module R",
			Description:   "Right aerodynamic wing",
			Price:         210000,
			StockQuantity: 7,
			Category:      repoModel.PartCategoryWing,
			Tags:          []string{"wing", "right"},
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
		{
			Uuid:          "66666666-6666-6666-6666-666666666666",
			Name:          "Auxiliary Engine B2",
			Description:   "Secondary maneuvering engine",
			Price:         480000,
			StockQuantity: 4,
			Category:      repoModel.PartCategoryEngine,
			Tags:          []string{"engine", "aux"},
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
		{
			Uuid:          "77777777-7777-7777-7777-777777777777",
			Name:          "Fuel Valve V1",
			Description:   "Fuel flow regulator",
			Price:         15000,
			StockQuantity: 40,
			Category:      repoModel.PartCategoryFuel,
			Tags:          []string{"fuel", "valve"},
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
		{
			Uuid:          "88888888-8888-8888-8888-888888888888",
			Name:          "Thermal Porthole",
			Description:   "Heat resistant porthole",
			Price:         92000,
			StockQuantity: 10,
			Category:      repoModel.PartCategoryPorthole,
			Tags:          []string{"window", "thermal"},
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
		{
			Uuid:          "99999999-9999-9999-9999-999999999999",
			Name:          "Fuel Pump P3",
			Description:   "High efficiency pump",
			Price:         60000,
			StockQuantity: 15,
			Category:      repoModel.PartCategoryFuel,
			Tags:          []string{"fuel", "pump"},
			CreatedAt:     lo.ToPtr(now),
			UpdatedAt:     lo.ToPtr(now),
		},
	}

	for _, part := range partList {
		err := r.CreatePart(ctx, lo.ToPtr(converter.RepoPartToModel(*part)))
		if err != nil {
			return err
		}
	}
	return nil
}
