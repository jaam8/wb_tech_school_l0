package errors

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	ErrOrderNotFound      = errors.New("order not found")
	ErrEmptyOrderUID      = errors.New("empty order uid")
	ErrOrderItemsNotFound = errors.New("order items not found")
)
