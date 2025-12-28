package model

import "github.com/google/uuid"

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "PAYMENT_METHOD_UNKNOWN_UNSPECIFIED"
	PaymentMethodCARD          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

type Order struct {
	UUID            uuid.UUID   `json:"uuid"`
	UserUUID        uuid.UUID   `json:"user_uuid"`
	PartUuids       []uuid.UUID `json:"part_uuids"`
	TotalPrice      float64
	TransactionUUID *uuid.UUID `json:"transaction_uuid"`
	OrderStatus     OrderStatus
	PaymentMethod   *PaymentMethod
}
