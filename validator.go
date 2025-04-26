package cobraflags

import (
	"fmt"
)

// Validator is an interface that defines a method for validating a value.
// It is used to validate the values of flags.
type Validator interface {
	Validate(any) error
}

// ValidatorFunc implements the Validator interface.
var _ Validator = (*ValidatorFunc[any])(nil)

// ValidatorFunc is a function type that implements the Validator interface.
// This function exists just for demonstration and testing purposes only.
// Use ValidateFunc field in FlagBase instead.
// Note, T must be the same type as the flag value.
type ValidatorFunc[T any] func(T) error

// Validate calls the ValidatorFunc itself to validate the value.
func (f ValidatorFunc[T]) Validate(value any) error {
	v, ok := value.(T)
	if !ok {
		return fmt.Errorf("invalid value type, expected %T, got %T", v, value)
	}
	return f(v)
}
