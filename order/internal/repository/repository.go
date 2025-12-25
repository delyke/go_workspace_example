package repository

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) (*model.Order, error)
	Get(ctx context.Context, uuid string) (*model.Order, error)
	Pay(ctx context.Context, uuid string, method model.PaymentMethod, txUUID string, status model.OrderStatus) (*model.Order, error)
	Cancel(ctx context.Context, uuid string, status model.OrderStatus) error
}
