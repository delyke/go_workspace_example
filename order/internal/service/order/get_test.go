package order

import "github.com/delyke/go_workspace_example/order/internal/model"

func (s *ServiceSuite) TestGetSuccess() {
	randomOrder := s.generateRandomOrder()

	s.orderRepository.On("Get", s.ctx, randomOrder.UUID.String()).Return(randomOrder, nil).Once()

	ansOrder, err := s.service.Get(s.ctx, randomOrder.UUID.String())
	s.Require().NoError(err)
	s.Require().Equal(randomOrder, ansOrder)
}

func (s *ServiceSuite) TestGetFailure() {
	var (
		repoErr = model.ErrOrderNotFound
		uuid    = s.faker.UUID()
	)

	s.orderRepository.On("Get", s.ctx, uuid).Return(nil, repoErr).Once()

	ansOrder, err := s.service.Get(s.ctx, uuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(ansOrder)
}
