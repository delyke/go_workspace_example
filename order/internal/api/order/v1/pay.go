package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/delyke/go_workspace_example/order/internal/converter"
	"github.com/delyke/go_workspace_example/order/internal/model"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	txID, err := a.orderService.Pay(ctx, params.OrderUUID.String(), converter.OpenApiPaymentMethodToModelOrderPayment(req.PaymentMethod))
	if err == nil {
		resUUID, err := uuid.Parse(txID)
		if err != nil {
			return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
		}

		return &orderV1.PayOrderResponse{
			TransactionUUID: resUUID,
		}, nil
	}

	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		return &orderV1.NotFoundError{Code: 404, Message: "order not found"}, nil
	case errors.Is(err, model.ErrPaymentBadGateway):
		return &orderV1.BadGatewayError{Code: 502, Message: "Bad Gateway"}, nil
	case errors.Is(err, model.ErrPaymentServiceUnavailable):
		return &orderV1.ServiceUnavailableError{Code: 503, Message: "Service Unavailable"}, nil
	case errors.Is(err, model.ErrPaymentServiceDeadlineExceeded):
		return &orderV1.GatewayTimeoutError{Code: 504, Message: "Gateway Timeout"}, nil
	case errors.Is(err, model.ErrPaymentInternalServerError):
		return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
	default:
		return &orderV1.InternalServerError{Code: 500, Message: "Internal Server Error"}, nil
	}
}
