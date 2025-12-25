package part

import "github.com/delyke/go_workspace_example/inventory/internal/model"

func (s *RepositorySuite) TestGetNotFound() {
	var (
		uuid    = s.faker.UUID()
		needErr = model.ErrPartNotFound
	)

	part, err := s.repository.GetPart(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().Empty(part)
	s.Require().ErrorIs(err, needErr)
}

func (s *RepositorySuite) TestGetSuccess() {
	uuid := "11111111-1111-1111-1111-111111111111"
	part, err := s.repository.GetPart(s.ctx, uuid)
	s.Require().NoError(err)
	s.Require().NotEmpty(part)
	s.Require().Equal(part.Uuid, uuid)
}
