package converter

import (
	"github.com/delyke/go_workspace_example/order/internal/model"
	generatedPaymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

func ModelPaymentMethodToProto(m model.PaymentMethod) generatedPaymentV1.PaymentMethod {
	switch m {
	case model.PaymentMethodCARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDITCARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}
