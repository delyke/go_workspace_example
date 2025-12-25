package part

import (
	"github.com/delyke/go_workspace_example/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
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
	)
	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(&part, nil)
	ansPart, err := s.service.GetPart(s.ctx, part.Uuid)
	s.Require().NoError(err)
	s.Require().Equal(part.Uuid, ansPart.Uuid)
}

func (s *ServiceSuite) TestGetError() {
	var (
		uuid    = s.faker.UUID()
		repoErr = s.faker.Error()
	)
	s.partRepository.On("GetPart", s.ctx, uuid).Return(nil, repoErr)
	ansPart, err := s.service.GetPart(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(ansPart)
}
