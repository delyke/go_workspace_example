package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderPayed    = errors.New("order has been payed")
)

var (
	ErrInventoryBadGateway              = errors.New("inventory bad gateway")
	ErrInventoryServiceUnavailable      = errors.New("inventory service unavailable")
	ErrInventoryServiceDeadlineExceeded = errors.New("inventory service deadline exceeded")
	ErrInventoryInternalServerError     = errors.New("inventory service internal server error")
	ErrInventoryPartNotFound            = errors.New("inventory part not found")
)

var (
	ErrPaymentBadGateway              = errors.New("payment bad gateway")
	ErrPaymentServiceUnavailable      = errors.New("payment service unavailable")
	ErrPaymentServiceDeadlineExceeded = errors.New("payment service deadline exceeded")
	ErrPaymentInternalServerError     = errors.New("payment service internal server error")
)
