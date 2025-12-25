package order

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"

	clientsMocks "github.com/delyke/go_workspace_example/order/internal/client/grpc/mocks"
	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context //nolint:containedctx
	orderRepository *mocks.OrderRepository
	inventoryClient *clientsMocks.InventoryClient
	paymentClient   *clientsMocks.PaymentClient
	service         *service
	faker           *gofakeit.Faker
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = clientsMocks.NewInventoryClient(s.T())
	s.paymentClient = clientsMocks.NewPaymentClient(s.T())
	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)
	s.faker = gofakeit.New(43)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (s *ServiceSuite) generateRandomOrder() *model.Order {
	return &model.Order{
		UUID:            s.faker.UUID(),
		UserUUID:        s.faker.UUID(),
		PartUuids:       []string{s.faker.UUID(), s.faker.UUID(), s.faker.UUID()},
		TotalPrice:      s.faker.Price(1.24, 6000.0),
		TransactionUUID: s.getTransactionUUIDOrNil(),
		OrderStatus:     s.getRandomOrderStatus(),
		PaymentMethod:   s.getRandomPaymentMethodOrNil(),
	}
}

func (s *ServiceSuite) generateTransactionUUID() string {
	return s.faker.UUID()
}

func (s *ServiceSuite) getTransactionUUIDOrNil() *string {
	if s.faker.Bool() {
		return nil
	}
	return lo.ToPtr(s.generateTransactionUUID())
}

func (s *ServiceSuite) getRandomOrderStatus() model.OrderStatus {
	orderStatusList := s.getOrderStatusLists()
	return orderStatusList[s.faker.Number(0, len(orderStatusList)-1)]
}

func (s *ServiceSuite) getOrderStatusLists() []model.OrderStatus {
	return []model.OrderStatus{
		model.OrderStatusPENDINGPAYMENT,
		model.OrderStatusPAID,
		model.OrderStatusCANCELLED,
	}
}

func (s *ServiceSuite) getPaymentMethodList() []model.PaymentMethod {
	return []model.PaymentMethod{
		model.PaymentMethodUNKNOWN,
		model.PaymentMethodSBP,
		model.PaymentMethodCARD,
		model.PaymentMethodCREDITCARD,
		model.PaymentMethodINVESTORMONEY,
	}
}

func (s *ServiceSuite) getRandomPaymentMethod() model.PaymentMethod {
	paymentMethodList := s.getPaymentMethodList()
	return paymentMethodList[s.faker.Number(0, len(paymentMethodList)-1)]
}

func (s *ServiceSuite) getRandomPaymentMethodOrNil() *model.PaymentMethod {
	if s.faker.Bool() {
		return nil
	}
	return lo.ToPtr(s.getRandomPaymentMethod())
}
