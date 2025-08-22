package errors

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotFount      = errors.New("order not found")
)
