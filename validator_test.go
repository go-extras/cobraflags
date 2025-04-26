package cobraflags_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

// TestValidatorFunc_Success tests that ValidatorFunc validates a value successfully.
func TestValidatorFunc_Success(t *testing.T) {
	c := qt.New(t)

	validator := cobraflags.ValidatorFunc[int](func(v int) error {
		if v < 0 {
			return fmt.Errorf("value must be non-negative")
		}
		return nil
	})

	err := validator.Validate(10)
	c.Assert(err, qt.IsNil)
}

// TestValidatorFunc_Failure tests that ValidatorFunc returns an error for invalid values.
func TestValidatorFunc_Failure(t *testing.T) {
	c := qt.New(t)

	validator := cobraflags.ValidatorFunc[int](func(v int) error {
		if v < 0 {
			return fmt.Errorf("value must be non-negative")
		}
		return nil
	})

	err := validator.Validate(-5)
	c.Assert(err, qt.Not(qt.IsNil))
	c.Assert(err.Error(), qt.Equals, "value must be non-negative")
}

// TestValidatorFunc_InvalidType tests that ValidatorFunc returns an error for invalid types.
func TestValidatorFunc_InvalidType(t *testing.T) {
	c := qt.New(t)

	validator := cobraflags.ValidatorFunc[int](func(v int) error {
		if v < 0 {
			return fmt.Errorf("value must be non-negative")
		}
		return nil
	})

	err := validator.Validate("invalid")
	c.Assert(err, qt.Not(qt.IsNil))
	c.Assert(err.Error(), qt.Matches, "invalid value type, expected.*")
}
