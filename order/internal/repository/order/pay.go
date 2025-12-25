package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/converter"
)

func (r *repository) Pay(
	_ context.Context,
	uuid string,
	method model.PaymentMethod,
	txUUID string,
	status model.OrderStatus,
) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[uuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}
	order.TransactionUUID = &txUUID
	order.OrderStatus = string(status)
	order.PaymentMethod = lo.ToPtr(string(method))
	return converter.RepoToOrderModel(order)
}
