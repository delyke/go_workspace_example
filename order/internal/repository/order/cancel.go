package order

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (r *repository) Cancel(_ context.Context, uuid string, status model.OrderStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}
	order.OrderStatus = string(status)
	return nil
}
