package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	txID, err := a.paymentService.PayOrder(ctx, req.GetUuid(), req.GetUserUuid(), string(req.GetPaymentMethod()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &paymentV1.PayOrderResponse{TransactionUuid: txID}, nil
}
