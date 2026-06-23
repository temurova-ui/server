package errs

import (
	"errors"
	"fmt"
)

var(
	ErrInvalidInput = errors.New("invalid input data")
	// ErrUserNotFound = errors.New("user  not found")
	ErrEmailConflict = errors.New("email already registered")
	ErrInternalServer = errors.New("internal server error")
)

func NewInvalidInputError(msg string) error{
	return fmt.Errorf("%w, %s", ErrInvalidInput, msg)
}