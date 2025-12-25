package v1

import (
	"context"
	"errors"

	"github.com/delyke/go_workspace_example/order/internal/converter"
	"github.com/delyke/go_workspace_example/order/internal/model"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error) {
	ord, err := a.orderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
		}
		return nil, err
	}
	ansOrder, err := converter.ModelOrderToOpenApiOrder(ord)
	if err != nil {
		return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
	}
	return &ansOrder, nil
}
