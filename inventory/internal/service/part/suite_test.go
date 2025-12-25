package part

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx            context.Context //nolint:containedctx
	partRepository *mocks.PartRepository
	service        *service
	faker          *gofakeit.Faker
}

// Подготавливает какие-то данные до теста
func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.partRepository = mocks.NewPartRepository(s.T())
	s.service = NewService(s.partRepository)
	s.faker = gofakeit.New(42)
}

// Тут можно что-то подчистить после того как тесты завершены
func (s *ServiceSuite) TearDownTest() {}

// Без этого не заработает
func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) randomTags(min, max int) []string {
	if min < 0 {
		min = 0
	}
	if max < min {
		max = min
	}

	n := 0
	if max > 0 {
		n = s.faker.Number(min, max)
	}

	if n == 0 {
		return nil
	}

	dict := []string{
		"space", "orbital", "mil-spec", "lightweight", "heavy-duty", "certified",
		"prototype", "production", "refurbished", "v2", "v3", "export",
		"sealed", "vacuum-rated", "high-temp", "low-temp",
	}

	seen := make(map[string]struct{}, n)
	out := make([]string, 0, n)

	limit := n * 3
	for i := 0; i < limit && len(out) < n; i++ {
		tag := dict[s.faker.Number(0, len(dict)-1)]
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

func (s *ServiceSuite) randomMetadata(minKeys, maxKeys int) map[string]*model.PartMetadataValue {
	if minKeys < 0 {
		minKeys = 0
	}
	if maxKeys < minKeys {
		maxKeys = minKeys
	}

	n := 0
	if maxKeys > 0 {
		n = s.faker.Number(minKeys, maxKeys)
	}
	if n == 0 {
		return nil
	}

	possibleKeys := []string{
		"serial", "batch", "material", "revision", "certified", "max_temp_c",
		"min_temp_c", "pressure_bar", "voltage_v", "power_kw", "notes",
	}

	meta := make(map[string]*model.PartMetadataValue, n)
	seen := make(map[string]struct{}, n)

	limit := n * 4
	for i := 0; i < limit && len(meta) < n; i++ {
		key := possibleKeys[s.faker.Number(0, len(possibleKeys)-1)]
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		meta[key] = s.randomMetadataValueForKey(key)
	}

	if len(meta) == 0 {
		return nil
	}
	return meta
}

func (s *ServiceSuite) randomMetadataValueForKey(key string) *model.PartMetadataValue {
	switch key {
	case "serial", "batch":
		s := s.faker.UUID()
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindString,
			String: &s,
		}
	case "material":
		s := []string{"titanium", "aluminum", "steel", "carbon-fiber"}[gofakeit.Number(0, 3)]
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindString,
			String: &s,
		}
	case "revision":
		i := int64(s.faker.Number(1, 10))
		return &model.PartMetadataValue{
			Kind:  model.MetadataKindInt64,
			Int64: &i,
		}
	case "certified":
		b := s.faker.Bool()
		return &model.PartMetadataValue{
			Kind: model.MetadataKindBool,
			Bool: &b,
		}
	case "max_temp_c":
		v := s.faker.Float64Range(80, 1200)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindDouble,
			Double: &v,
		}
	case "min_temp_c":
		v := s.faker.Float64Range(-200, 20)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindDouble,
			Double: &v,
		}
	case "pressure_bar":
		v := s.faker.Float64Range(0.5, 500)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindDouble,
			Double: &v,
		}
	case "voltage_v":
		v := s.faker.Float64Range(3.3, 1000)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindDouble,
			Double: &v,
		}
	case "power_kw":
		v := s.faker.Float64Range(0.1, 2000)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindDouble,
			Double: &v,
		}
	default:
		s := s.faker.LoremIpsumSentence(8)
		return &model.PartMetadataValue{
			Kind:   model.MetadataKindString,
			String: &s,
		}
	}
}

func (s *ServiceSuite) randomManufacturerOrNil() *model.PartManufacturer {
	if s.faker.Bool() {
		return nil
	}
	return s.randomManufacturer()
}

func (s *ServiceSuite) randomManufacturer() *model.PartManufacturer {
	return &model.PartManufacturer{
		Name:    s.faker.Company(),
		Country: s.faker.Country(),
		Website: s.faker.URL(),
	}
}

func (s *ServiceSuite) randomPartDimensions() *model.PartDimensions {
	return &model.PartDimensions{
		Length: s.faker.Float64Range(0.2, 12.0),   // м
		Width:  s.faker.Float64Range(0.2, 5.0),    // м
		Height: s.faker.Float64Range(0.2, 4.0),    // м
		Weight: s.faker.Float64Range(1.0, 5000.0), // кг
	}
}

func (s *ServiceSuite) randomPartDimensionsOrNil() *model.PartDimensions {
	if s.faker.Bool() {
		return nil
	}
	return s.randomPartDimensions()
}

func (s *ServiceSuite) getPartCategoriesList() []model.PartCategory {
	partCategories := []model.PartCategory{
		model.PartCategoryEngine,
		model.PartCategoryFuel,
		model.PartCategoryPorthole,
		model.PartCategoryWing,
	}
	return partCategories
}
