package converter

import (
	"github.com/google/uuid"

	"github.com/delyke/go_workspace_example/order/internal/model"
	orderV1 "github.com/delyke/go_workspace_example/shared/pkg/openapi/order/v1"
	paymentV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(m orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case orderV1.PaymentMethodPAYMENTMETHODUNKNOWNUNSPECIFIED:
		fallthrough
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNKNOWN_UNSPECIFIED
	}
}

func ModelOrderStatusToOpenApi(m model.OrderStatus) orderV1.OrderStatus {
	switch m {
	case model.OrderStatusPENDINGPAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func ModelOrderPaymentMethodToOpenApi(m model.PaymentMethod) orderV1.PaymentMethod {
	switch m {
	case model.PaymentMethodUNKNOWN:
		return orderV1.PaymentMethodPAYMENTMETHODUNKNOWNUNSPECIFIED
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodPAYMENTMETHODCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodPAYMENTMETHODSBP
	case model.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodPAYMENTMETHODCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		return orderV1.PaymentMethodPAYMENTMETHODUNKNOWNUNSPECIFIED
	}
}

func OpenApiPaymentMethodToModelOrderPayment(m orderV1.PaymentMethod) model.PaymentMethod {
	switch m {
	case orderV1.PaymentMethodPAYMENTMETHODCARD:
		return model.PaymentMethodCARD
	case orderV1.PaymentMethodPAYMENTMETHODCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case orderV1.PaymentMethodPAYMENTMETHODUNKNOWNUNSPECIFIED:
		return model.PaymentMethodUNKNOWN
	case orderV1.PaymentMethodPAYMENTMETHODSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodPAYMENTMETHODINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func PartUUIDSOpenApiToModel(partUUIDS []uuid.UUID) []string {
	var partUUIDSOpenApi []string
	for _, id := range partUUIDS {
		partUUIDSOpenApi = append(partUUIDSOpenApi, id.String())
	}
	return partUUIDSOpenApi
}

func ModelOrderToOpenApiOrder(m *model.Order) (orderV1.OrderDto, error) {
	var partUuids []uuid.UUID
	if m.PartUuids != nil {
		partUuids = make([]uuid.UUID, len(m.PartUuids))
		copy(partUuids, m.PartUuids)
	}

	var transactionUUID orderV1.OptUUID

	if m.TransactionUUID != nil {
		transactionUUID = orderV1.OptUUID{
			Value: *m.TransactionUUID,
			Set:   true,
		}
	} else {
		transactionUUID = orderV1.OptUUID{
			Set: false,
		}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if m.PaymentMethod != nil {
		paymentMethod = orderV1.OptPaymentMethod{
			Value: ModelOrderPaymentMethodToOpenApi(*m.PaymentMethod),
		}
	}

	return orderV1.OrderDto{
		OrderUUID:       m.UUID,
		UserUUID:        m.UserUUID,
		PartUuids:       partUuids,
		TotalPrice:      m.TotalPrice,
		TransactionUUID: transactionUUID,
		Status:          ModelOrderStatusToOpenApi(m.OrderStatus),
		PaymentMethod:   paymentMethod,
	}, nil
}
