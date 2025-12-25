package order

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}
	return order, nil
}
