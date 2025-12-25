package order

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderID string) error {
	order, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return err
	}

	if order.OrderStatus == model.OrderStatusPAID {
		return model.ErrOrderPayed
	}

	err = s.orderRepository.Cancel(ctx, orderID, model.OrderStatusCANCELLED)
	if err != nil {
		return err
	}
	return nil
}
