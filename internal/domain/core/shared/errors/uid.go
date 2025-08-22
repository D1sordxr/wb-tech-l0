package errors

import "errors"

var (
	ErrOrderUIDInvalidLength = errors.New("order UID must be 20 characters long")
	ErrOrderUIDInvalidSuffix = errors.New("order UID must end with 'test'")
	ErrOrderUIDInvalidChars  = errors.New("order UID can only contain lowercase letters and digits")
)
