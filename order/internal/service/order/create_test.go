package order

import (
	"errors"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *ServiceSuite) TestCreateListPartsServiceError() {
	var (
		userUUID   = s.faker.UUID()
		partUUIDs  = []string{s.faker.UUID(), s.faker.UUID()}
		notGrpcErr = errors.New("not grpc error")
		needError  = model.ErrInventoryInternalServerError
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return([]model.Part{}, notGrpcErr)

	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, needError)
}

func (s *ServiceSuite) TestCreateListPartsServiceUnavailableError() {
	var (
		userUUID   = s.faker.UUID()
		partUUIDs  = []string{s.faker.UUID(), s.faker.UUID()}
		notGrpcErr = grpcStatus.Error(codes.Unavailable, "inventory down")
		needError  = model.ErrInventoryServiceUnavailable
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return([]model.Part{}, notGrpcErr)

	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, needError)
}

func (s *ServiceSuite) TestCreateListPartsServiceDeadlineError() {
	var (
		userUUID   = s.faker.UUID()
		partUUIDs  = []string{s.faker.UUID(), s.faker.UUID()}
		notGrpcErr = grpcStatus.Error(codes.DeadlineExceeded, "deadline exceeded")
		needError  = model.ErrInventoryServiceDeadlineExceeded
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return([]model.Part{}, notGrpcErr)

	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, needError)
}

func (s *ServiceSuite) TestCreateListPartsServiceBadGatewayError() {
	var (
		userUUID   = s.faker.UUID()
		partUUIDs  = []string{s.faker.UUID(), s.faker.UUID()}
		notGrpcErr = grpcStatus.Error(codes.Aborted, "grpc random")
		needError  = model.ErrInventoryBadGateway
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return([]model.Part{}, notGrpcErr)

	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, needError)
}

func (s *ServiceSuite) TestCreateListPartsNotFoundError() {
	var (
		userUUID  = s.faker.UUID()
		partUUIDs = []string{s.faker.UUID(), s.faker.UUID()}
		parts     = []model.Part{
			{
				Price: s.faker.Price(0.0, 5000.0),
			},
		}
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return(parts, nil)
	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, model.ErrInventoryPartNotFound)
}

func (s *ServiceSuite) TestCreateRepoCreateError() {
	var (
		userUUID  = s.faker.UUID()
		partUUIDs = []string{s.faker.UUID(), s.faker.UUID()}
		parts     = []model.Part{
			{
				Price: s.faker.Price(0.0, 5000.0),
			},
			{
				Price: s.faker.Price(0.0, 5000.0),
			},
		}
		repoErr = errors.New("repo error")
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return(parts, nil)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(nil, repoErr).Once()
	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().Error(err)
	s.Require().Empty(ansUUID)
	s.Require().Empty(price)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		userUUID  = s.faker.UUID()
		partUUIDs = []string{s.faker.UUID(), s.faker.UUID()}
		parts     = []model.Part{
			{
				Price: s.faker.Price(0.0, 5000.0),
			},
			{
				Price: s.faker.Price(0.0, 5000.0),
			},
		}
		createdOrder = s.generateRandomOrder()
	)

	s.inventoryClient.On("ListParts", s.ctx, model.PartsFilter{
		UUIDs: partUUIDs,
	}).Return(parts, nil)
	s.orderRepository.On("Create", s.ctx, mock.Anything).Return(createdOrder, nil).Once()
	ansUUID, price, err := s.service.Create(s.ctx, userUUID, partUUIDs)
	s.Require().NoError(err)
	s.Require().NotEmpty(ansUUID)
	s.Require().Equal(createdOrder.UUID.String(), ansUUID)
	s.Require().Equal(price, createdOrder.TotalPrice)
}
