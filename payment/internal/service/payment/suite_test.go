package payment

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	ctx     context.Context //nolint:containedctx
	service *service
	faker   *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = NewService()
	s.faker = gofakeit.New(44)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
