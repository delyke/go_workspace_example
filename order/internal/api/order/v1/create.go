package v1

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/google/uuid"

	"github.com/delyke/go_workspace_example/order/internal/converter"
	"github.com/delyke/go_workspace_example/order/internal/model"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, params *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	convertedPartUUIDS := converter.PartUUIDSOpenApiToModel(params.PartUuids)
	oUUID, totalPrice, err := a.orderService.Create(ctx, params.UserUUID.String(), convertedPartUUIDS)

	if err == nil {
		resUUID, err := uuid.Parse(oUUID)
		if err != nil {
			return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
		}
		return &orderV1.CreateOrderResponse{
			UUID:       resUUID,
			TotalPrice: totalPrice,
		}, nil
	}

	switch {
	case errors.Is(err, model.ErrInventoryBadGateway):
		return &orderV1.BadGatewayError{Code: 502, Message: "Bad Gateway"}, nil
	case errors.Is(err, model.ErrInventoryServiceUnavailable):
		return &orderV1.ServiceUnavailableError{Code: 503, Message: "Service Unavailable"}, nil
	case errors.Is(err, model.ErrInventoryServiceDeadlineExceeded):
		return &orderV1.GatewayTimeoutError{Code: 504, Message: "Gateway Timeout"}, nil
	case errors.Is(err, model.ErrInventoryInternalServerError):
		return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
	case errors.Is(err, model.ErrInventoryPartNotFound):
		return &orderV1.NotFoundError{Code: 404, Message: "Parts not found"}, nil
	default:
		return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
	}
}
