package model

type Order struct {
	UUID            string
	UserUUID        string
	PartUuids       []string
	TotalPrice      float64
	TransactionUUID *string
	OrderStatus     string
	PaymentMethod   *string
}
