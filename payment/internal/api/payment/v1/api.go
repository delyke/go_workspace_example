package v1

import (
	"github.com/delyke/go_workspace_example/payment/internal/service"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewApi(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
