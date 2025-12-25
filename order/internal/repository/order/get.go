package order

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, ok := r.orders[uuid]
	if !ok {
		return nil, model.ErrOrderNotFound
	}
	return converter.RepoToOrderModel(order)
}
