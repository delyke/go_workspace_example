package converter

import (
	"fmt"

	"github.com/delyke/go_workspace_example/order/internal/model"
	repoModel "github.com/delyke/go_workspace_example/order/internal/repository/model"
)

func paymentMethodPtrToStringPtr(pm *model.PaymentMethod) *string {
	if pm == nil {
		return nil
	}
	s := string(*pm)
	return &s
}

func orderStatusToString(os model.OrderStatus) string {
	return string(os)
}

func OrderToRepoModel(order *model.Order) *repoModel.Order {
	return &repoModel.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		OrderStatus:     orderStatusToString(order.OrderStatus),
		PaymentMethod:   paymentMethodPtrToStringPtr(order.PaymentMethod),
	}
}

func paymentMethodStringPtrToEnumPtr(s *string) (*model.PaymentMethod, error) {
	if s == nil {
		return nil, nil
	}

	pm := model.PaymentMethod(*s)

	switch pm {
	case model.PaymentMethodUNKNOWN,
		model.PaymentMethodCARD,
		model.PaymentMethodSBP,
		model.PaymentMethodCREDITCARD,
		model.PaymentMethodINVESTORMONEY:
		return &pm, nil
	default:
		return nil, fmt.Errorf("unknown payment method: %s", *s)
	}
}

func orderStatusStringToEnum(s string) (model.OrderStatus, error) {
	os := model.OrderStatus(s)

	switch os {
	case model.OrderStatusPENDINGPAYMENT,
		model.OrderStatusPAID,
		model.OrderStatusCANCELLED:
		return os, nil
	default:
		return "", fmt.Errorf("unknown order status: %s", s)
	}
}

func RepoToOrderModel(order *repoModel.Order) (*model.Order, error) {
	orderStatus, err := orderStatusStringToEnum(order.OrderStatus)
	if err != nil {
		return nil, err
	}

	paymentMethod, err := paymentMethodStringPtrToEnumPtr(order.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		UUID:            order.UUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		OrderStatus:     orderStatus,
		PaymentMethod:   paymentMethod,
	}, nil
}
