package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/order/mapper"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	inventoryAddress  = "localhost:50051"
	paymentAddress    = "localhost:50052"
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) GetOrderByUUID(uuid uuid.UUID) *orderV1.OrderDto {
	s.mu.Lock()
	defer s.mu.Unlock()
	order := s.orders[uuid.String()]
	return order
}

var ErrOrderNotFound = errors.New("order not found")

func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.orders[order.OrderUUID.String()]; !ok {
		return ErrOrderNotFound
	}

	s.orders[order.OrderUUID.String()] = order
	return nil
}

type OrderHandler struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}

func (h *OrderHandler) CancelOrder(_ context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrderByUUID(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
	}
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{Code: 409, Message: "order is paid"}, nil
	}
	order.Status = orderV1.OrderStatusCANCELLED
	err := h.storage.UpdateOrder(order)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		}
	}
	return &orderV1.CancelOrderNoContent{}, nil
}

func (h *OrderHandler) CreateOrder(ctx context.Context, params *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	partUuidsString := make([]string, len(params.PartUuids))
	for i, u := range params.PartUuids {
		partUuidsString[i] = u.String()
	}

	listParts, err := h.inventoryClient.ListParts(
		ctx,
		&inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartsFilter{
				Uuids: partUuidsString,
			},
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return &orderV1.BadGatewayError{Code: 502, Message: "Bad Gateway"}, nil
		}

		switch st.Code() {
		case codes.Unavailable:
			return &orderV1.ServiceUnavailableError{Code: 503, Message: "Service Unavailable"}, nil
		case codes.DeadlineExceeded:
			return &orderV1.GatewayTimeoutError{Code: 504, Message: "Gateway Timeout"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
		}
	}

	if len(listParts.Parts) != len(params.PartUuids) {
		return &orderV1.NotFoundError{Code: 404, Message: "Parts not found"}, nil
	}

	var totalPrice float64

	for _, part := range listParts.Parts {
		totalPrice += part.Price
	}

	order := &orderV1.OrderDto{
		OrderUUID:  uuid.New(),
		UserUUID:   params.UserUUID,
		PartUuids:  params.PartUuids,
		TotalPrice: totalPrice,
		Status:     orderV1.OrderStatusPENDINGPAYMENT,
	}

	h.storage.orders[order.OrderUUID.String()] = order
	return &orderV1.CreateOrderResponse{
		UUID:       order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

func (h *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	order := h.storage.GetOrderByUUID(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
	}
	return order, nil
}

func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrderByUUID(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
	}

	res, err := h.paymentClient.PayOrder(
		ctx,
		&paymentV1.PayOrderRequest{
			PaymentMethod: mapper.PaymentMethodToProto(req.GetPaymentMethod()),
			Uuid:          order.GetOrderUUID().String(),
			UserUuid:      order.GetUserUUID().String(),
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return &orderV1.BadGatewayError{Code: 502, Message: "Bad Gateway"}, nil
		}

		switch st.Code() {
		case codes.Unavailable:
			return &orderV1.ServiceUnavailableError{Code: 503, Message: "Service Unavailable"}, nil
		case codes.DeadlineExceeded:
			return &orderV1.GatewayTimeoutError{Code: 504, Message: "Gateway Timeout"}, nil
		default:
			return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
		}
	}
	txUUID, err := uuid.Parse(res.TransactionUuid)
	if err != nil {
		return &orderV1.BadGatewayError{Code: 502, Message: "Invalid transaction uuid from payment service"}, nil
	}

	order.PaymentMethod = orderV1.OptPaymentMethod{
		Value: req.PaymentMethod,
	}

	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.OptUUID{
		Value: txUUID,
		Set:   true,
	}
	err = h.storage.UpdateOrder(order)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: order.TransactionUUID.Value,
	}, nil
}

func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewOrderStorage()

	connInventory, err := grpc.NewClient(
		inventoryAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to create inventory client: %v", err)
		return
	}
	defer func() {
		if cerr := connInventory.Close(); cerr != nil {
			log.Printf("Failed to close inventory client: %v", cerr)
		}
	}()

	inventoryClient := inventoryV1.NewInventoryServiceClient(connInventory)

	connPayment, err := grpc.NewClient(
		paymentAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed to create payment client: %v", err)
		return
	}
	defer func() {
		if cerr := connPayment.Close(); cerr != nil {
			log.Printf("Failed to close payment client: %v", cerr)
		}
	}()
	paymentClient := paymentV1.NewPaymentServiceClient(connPayment)

	orderHandler := NewOrderHandler(storage, inventoryClient, paymentClient)

	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
