package part

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	ctx        context.Context //nolint:containedctx
	repository *repository
	faker      *gofakeit.Faker
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.faker = gofakeit.New(45)
	s.repository = NewRepository()
}

func (s *RepositorySuite) TearDownTest() {}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
