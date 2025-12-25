package v1

import (
	"context"
	"errors"

	"github.com/delyke/go_workspace_example/order/internal/model"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID.String())
	if err == nil {
		return &orderV1.CancelOrderNoContent{}, nil
	}

	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
	case errors.Is(err, model.ErrOrderPayed):
		return &orderV1.ConflictError{Code: 409, Message: "order is paid"}, nil
	default:
		return &orderV1.InternalServerError{Code: 500, Message: err.Error()}, nil
	}
}
