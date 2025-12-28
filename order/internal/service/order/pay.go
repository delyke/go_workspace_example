package order

import (
	"context"

	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *service) Pay(
	ctx context.Context,
	orderUUID string,
	paymentMethod model.PaymentMethod,
) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	txID, err := s.paymentClient.PayOrder(
		ctx,
		order.UUID.String(),
		order.UserUUID.String(),
		paymentMethod,
	)
	if err != nil {
		st, ok := grpcStatus.FromError(err)
		if !ok {
			return "", model.ErrPaymentInternalServerError
		}

		switch st.Code() {
		case codes.Unavailable:
			return "", model.ErrPaymentServiceUnavailable
		case codes.DeadlineExceeded:
			return "", model.ErrPaymentServiceDeadlineExceeded
		default:
			return "", model.ErrPaymentBadGateway
		}
	}

	payedOrder, err := s.orderRepository.Pay(
		ctx,
		order.UUID.String(),
		paymentMethod,
		txID,
		model.OrderStatusPAID,
	)
	if err != nil {
		return "", err
	}

	var trxUUID string
	if payedOrder.TransactionUUID != nil {
		trxUUID = payedOrder.TransactionUUID.String()
	}

	return trxUUID, nil
}
