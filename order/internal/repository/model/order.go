package model

import "github.com/google/uuid"

type Order struct {
	UUID            uuid.UUID
	UserUUID        uuid.UUID
	PartUuids       []uuid.UUID
	TotalPrice      float64
	TransactionUUID *uuid.UUID
	OrderStatus     string
	PaymentMethod   *string
}
