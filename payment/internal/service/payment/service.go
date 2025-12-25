package payment

import (
	def "github.com/delyke/go_workspace_example/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
