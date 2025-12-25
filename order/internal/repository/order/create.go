package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order *model.Order) (*model.Order, error) {
	newUUID := uuid.New()

	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := converter.OrderToRepoModel(order)
	repoOrder.UUID = newUUID.String()
	r.orders[newUUID.String()] = repoOrder
	return converter.RepoToOrderModel(repoOrder)
}
