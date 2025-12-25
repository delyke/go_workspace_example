package v1

import (
	"context"

	"github.com/delyke/go_workspace_example/order/internal/client/converter"
	"github.com/delyke/go_workspace_example/order/internal/model"
	generatedPaymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	resp, err := c.generatedClient.PayOrder(
		ctx,
		&generatedPaymentV1.PayOrderRequest{
			PaymentMethod: converter.ModelPaymentMethodToProto(paymentMethod),
			Uuid:          orderUUID,
			UserUuid:      userUUID,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.TransactionUuid, nil
}
