package errors

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderItemsNotFound = errors.New("order items not found")
)
