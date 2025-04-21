package cobraflags_test

import (
	"errors"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func TestBoolFlag_Register(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:  "name",
		Value: false,
		Usage: "usage",
	}

	flag.Register(cmd)

	const expectedValue = true
	cmd.SetArgs([]string{"--name", boolToString(expectedValue)})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetBool(), qt.Equals, expectedValue)
}

func TestBoolFlag_GetBoolE(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:  "feature",
		Value: false,
		Usage: "enable feature",
	}

	flag.Register(cmd)

	const expectedValue = true
	cmd.SetArgs([]string{"--feature", boolToString(expectedValue)})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	value, err := flag.GetBoolE()
	c.Assert(err, qt.IsNil)
	c.Assert(value, qt.Equals, expectedValue)
}

func TestBoolFlag_WithShorthand(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:      "verbose",
		Shorthand: "v",
		Value:     false,
		Usage:     "verbose output",
	}

	flag.Register(cmd)

	const expectedValue = true
	cmd.SetArgs([]string{"-v", boolToString(expectedValue)})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetBool(), qt.Equals, expectedValue)
}

func TestBoolFlag_WithRequired(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:     "confirm",
		Value:    false,
		Usage:    "confirm action",
		Required: true,
	}

	flag.Register(cmd)

	// Test missing required flag
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	c.Assert(err, qt.Not(qt.IsNil))

	// Test with required flag provided
	cmd.SetArgs([]string{"--confirm", "true"})
	err = cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetBool(), qt.Equals, true)
}

func TestBoolFlag_Persistent(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:       "debug",
		Value:      false,
		Usage:      "enable debug",
		Persistent: true,
	}

	flag.Register(cmd)

	// Verify the flag is registered to PersistentFlags
	f := cmd.PersistentFlags().Lookup("debug")
	c.Assert(f, qt.Not(qt.IsNil))

	const expectedValue = true
	cmd.SetArgs([]string{"--debug", boolToString(expectedValue)})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetBool(), qt.Equals, expectedValue)
}

func TestBoolFlag_WithValidation(t *testing.T) {
	c := qt.New(t)

	var validationCalled bool

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:  "feature",
		Value: false,
		Usage: "enable feature",
		ValidateFunc: func(value bool) error {
			return errors.New("invalid value")
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--feature", "true"})
	err := cmd.Execute()
	c.Assert(err, qt.IsNil)

	// GetBool doesn't call validation
	_ = flag.GetBool()
	c.Assert(validationCalled, qt.IsFalse)

	// GetBoolE calls validation
	_, err = flag.GetBoolE()
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "invalid value")
}
