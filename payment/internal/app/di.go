package app

import (
	"context"

	paymentApiV1 "github.com/delyke/go_workspace_example/payment/internal/api/payment/v1"
	"github.com/delyke/go_workspace_example/payment/internal/service"
	paymentService "github.com/delyke/go_workspace_example/payment/internal/service/payment"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentService service.PaymentService
	paymentAPIV1   paymentV1.PaymentServiceServer
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (c *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if c.paymentService == nil {
		c.paymentService = paymentService.NewService()
	}
	return c.paymentService
}

func (c *diContainer) PaymentAPIV1(ctx context.Context) paymentV1.PaymentServiceServer {
	if c.paymentAPIV1 == nil {
		c.paymentAPIV1 = paymentApiV1.NewApi(c.PaymentService(ctx))
	}
	return c.paymentAPIV1
}
