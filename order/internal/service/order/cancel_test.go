package order

import "github.com/delyke/go_workspace_example/order/internal/model"

func (s *ServiceSuite) TestCancelOrderGetRepoError() {
	var (
		order   = s.generateRandomOrder()
		repoErr = model.ErrOrderNotFound
	)

	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(nil, repoErr).Once()
	err := s.service.Cancel(s.ctx, order.UUID.String())
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancelOrderPayedError() {
	var (
		order   = s.generateRandomOrder()
		needErr = model.ErrOrderPayed
	)
	order.OrderStatus = model.OrderStatusPAID
	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	err := s.service.Cancel(s.ctx, order.UUID.String())
	s.Require().Error(err)
	s.Require().ErrorIs(err, needErr)
}

func (s *ServiceSuite) TestCancelOrderCancelRepoError() {
	var (
		order       = s.generateRandomOrder()
		needRepoErr = model.ErrOrderNotFound
	)
	order.OrderStatus = model.OrderStatusPENDINGPAYMENT
	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.orderRepository.On("Cancel", s.ctx, order.UUID.String(), model.OrderStatusCANCELLED).Return(needRepoErr).Once()
	err := s.service.Cancel(s.ctx, order.UUID.String())
	s.Require().Error(err)
	s.Require().ErrorIs(err, needRepoErr)
}

func (s *ServiceSuite) TestCancelOrderSuccess() {
	order := s.generateRandomOrder()
	order.OrderStatus = model.OrderStatusPENDINGPAYMENT
	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.orderRepository.On("Cancel", s.ctx, order.UUID.String(), model.OrderStatusCANCELLED).Return(nil).Once()
	err := s.service.Cancel(s.ctx, order.UUID.String())
	s.Require().NoError(err)
	s.Require().Empty(err)
}
