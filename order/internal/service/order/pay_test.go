package order

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderGetError() {
	var (
		uuid          = s.faker.UUID()
		repoErr       = model.ErrOrderNotFound
		paymentMethod = s.getRandomPaymentMethod()
	)
	s.orderRepository.On("Get", s.ctx, uuid).Return(nil, repoErr).Once()

	ansOrder, err := s.service.Pay(s.ctx, uuid, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentClientInternalServerError() {
	var (
		order               = s.generateRandomOrder()
		paymentMethod       = s.getRandomPaymentMethod()
		paymentServiceError = errors.New("not grpc error")
		payReturnErr        = model.ErrPaymentInternalServerError
	)

	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return("", paymentServiceError).Once()
	ansOrder, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, payReturnErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentClientUnavailableError() {
	var (
		order               = s.generateRandomOrder()
		paymentMethod       = s.getRandomPaymentMethod()
		paymentServiceError = grpcStatus.Error(codes.Unavailable, "payment down")
		payReturnErr        = model.ErrPaymentServiceUnavailable
	)

	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return("", paymentServiceError).Once()
	ansOrder, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, payReturnErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentClientDeadlineExceededError() {
	var (
		order               = s.generateRandomOrder()
		paymentMethod       = s.getRandomPaymentMethod()
		paymentServiceError = grpcStatus.Error(codes.DeadlineExceeded, "deadline exceeded")
		payReturnErr        = model.ErrPaymentServiceDeadlineExceeded
	)

	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return("", paymentServiceError).Once()
	ansOrder, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, payReturnErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentClientBadGatewayError() {
	var (
		order               = s.generateRandomOrder()
		paymentMethod       = s.getRandomPaymentMethod()
		paymentServiceError = grpcStatus.Error(codes.Aborted, "random error from grpc")
		payReturnErr        = model.ErrPaymentBadGateway
	)

	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return("", paymentServiceError).Once()
	ansOrder, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, payReturnErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentRepositoryPayError() {
	var (
		order         = s.generateRandomOrder()
		paymentMethod = s.getRandomPaymentMethod()
		txID          = s.faker.UUID()
		repoErr       = model.ErrOrderNotFound
	)
	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return(txID, nil).Once()
	s.orderRepository.
		On("Pay", s.ctx, order.UUID.String(), paymentMethod, txID, model.OrderStatusPAID).
		Return(nil, repoErr).Once()

	ansOrder, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(ansOrder)
}

func (s *ServiceSuite) TestPaymentRepositoryPaySuccess() {
	var (
		order         = s.generateRandomOrder()
		paymentMethod = s.getRandomPaymentMethod()
		txID          = s.faker.UUID()
	)

	genUUID, err := uuid.Parse(txID)
	s.Require().NoError(err)
	order.TransactionUUID = lo.ToPtr(genUUID)

	log.Printf("%#v", order)
	s.orderRepository.On("Get", s.ctx, order.UUID.String()).Return(order, nil).Once()
	s.paymentClient.
		On("PayOrder", s.ctx, order.UUID.String(), order.UserUUID.String(), paymentMethod).
		Return(txID, nil).Once()
	s.orderRepository.
		On("Pay", s.ctx, order.UUID.String(), paymentMethod, txID, model.OrderStatusPAID).
		Return(order, nil).Once()

	ansTxUUID, err := s.service.Pay(s.ctx, order.UUID.String(), paymentMethod)
	s.Require().NoError(err)

	var txPrepare string
	if order.TransactionUUID != nil {
		txPrepare = order.TransactionUUID.String()
	}

	s.Require().Equal(txPrepare, ansTxUUID)
}
