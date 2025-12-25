package part

import (
	"errors"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
)

func (s *ServiceSuite) TestListSuccess() {
	var (
		partCategories = s.getPartCategoriesList()
		uuid           = s.faker.UUID()
		name           = s.faker.Name()
		description    = s.faker.ProductDescription()
		price          = s.faker.Price(1000, 5000)
		stockQuantity  = s.faker.Number(1, 100)
		category       = partCategories[s.faker.Number(0, len(partCategories)-1)]
		manufacturer   = s.randomManufacturerOrNil()
		dimensions     = s.randomPartDimensionsOrNil()
		tags           = s.randomTags(1, 5)
		metadata       = s.randomMetadata(0, 5)
		createdAt      = s.faker.Date()
		updatedAt      = s.faker.Date()

		part = model.Part{
			Uuid:          uuid,
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: int64(stockQuantity),
			Category:      category,
			Manufacturer:  manufacturer,
			Dimensions:    dimensions,
			Tags:          tags,
			Metadata:      metadata,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
		}

		parts   = []*model.Part{&part}
		filters = &model.PartsFilter{
			UUIDs:                 []string{},
			Names:                 []string{},
			Categories:            []model.PartCategory{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}
	)
	s.partRepository.On("ListParts", s.ctx, filters).Return(parts, nil)
	ansList, err := s.service.ListParts(s.ctx, filters)
	s.Require().NoError(err)
	s.Require().Equal(parts, ansList)
}

func (s *ServiceSuite) TestListFail() {
	var (
		filters = &model.PartsFilter{
			UUIDs:                 []string{},
			Names:                 []string{},
			Categories:            []model.PartCategory{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}

		repoErr = errors.New("unexpected error")
	)
	s.partRepository.On("ListParts", s.ctx, filters).Return(nil, repoErr)

	ansList, err := s.service.ListParts(s.ctx, filters)

	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(ansList)
}
