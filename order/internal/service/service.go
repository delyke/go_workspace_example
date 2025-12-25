package service

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

type OrderService interface {
	Cancel(ctx context.Context, orderID string) error
	Create(ctx context.Context, userUUID string, partUUIDs []string) (string, float64, error)
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
}
