package cobraflags_test

import (
	"errors"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func TestUint8Flag_Register(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:  "level",
		Value: 0,
		Usage: "set level",
	}

	flag.Register(cmd)

	const expectedValue uint8 = 42
	cmd.SetArgs([]string{"--level", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetUint8(), qt.Equals, expectedValue)
}

func TestUint8Flag_GetUint8E(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:  "level",
		Value: 0,
		Usage: "set level",
	}

	flag.Register(cmd)

	const expectedValue uint8 = 42
	cmd.SetArgs([]string{"--level", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	value, err := flag.GetUint8E()
	c.Assert(err, qt.IsNil)
	c.Assert(value, qt.Equals, expectedValue)
}

func TestUint8Flag_WithShorthand(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:      "level",
		Value:     0,
		Usage:     "set level",
		Shorthand: "l",
	}

	flag.Register(cmd)

	const expectedValue uint8 = 42
	cmd.SetArgs([]string{"-l", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetUint8(), qt.Equals, expectedValue)
}

func TestUint8Flag_WithDefaultValue(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:  "level",
		Value: 10,
		Usage: "set level",
	}

	flag.Register(cmd)

	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetUint8(), qt.Equals, uint8(10))
}

func TestUint8Flag_WithRequired(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:     "level",
		Value:    0,
		Usage:    "set level",
		Required: true,
	}

	flag.Register(cmd)

	// Test missing required flag
	cmd.SetArgs(make([]string, 0))
	err := cmd.Execute()

	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "required flag(s) \"level\" not set")

	// Test with required flag provided
	cmd.SetArgs([]string{"--level", "42"})
	err = cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetUint8(), qt.Equals, uint8(42))
}

func TestUint8Flag_ValidateFunc(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:  "level",
		Value: 0,
		Usage: "set level",
		ValidateFunc: func(v uint8) error {
			if v > 100 {
				return errors.New("level must be <= 100")
			}
			return nil
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--level", "150"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetUint8E()
	c.Assert(err.Error(), qt.Equals, "level must be <= 100")
}

func TestUint8Flag_Validator(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:  "level",
		Value: 0,
		Usage: "set level",
		Validator: cobraflags.ValidatorFunc[uint8](func(v uint8) error {
			if v > 100 {
				return errors.New("level must be <= 100")
			}
			return nil
		}),
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--level", "150"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetUint8E()
	c.Assert(err.Error(), qt.Equals, "level must be <= 100")
}

func TestUint8Flag_WithPersistent(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.Uint8Flag{
		Name:       "level",
		Value:      0,
		Usage:      "set level",
		Persistent: true,
	}

	flag.Register(cmd)

	// Verify the flag is registered to PersistentFlags
	f := cmd.PersistentFlags().Lookup("level")
	c.Assert(f, qt.IsNotNil)

	const expectedValue uint8 = 42
	cmd.SetArgs([]string{"--level", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetUint8(), qt.Equals, expectedValue)
}

// TestUint8Flag_ViperKey_HappyPath tests ViperKey functionality with successful scenarios.
func TestUint8Flag_ViperKey_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		flagName    string
		viperKey    string
		expectedKey string
	}{
		{
			name:        "with_viper_key_set",
			flagName:    "log-level",
			viperKey:    "logging.level",
			expectedKey: "logging.level",
		},
		{
			name:        "with_empty_viper_key_fallback_to_name",
			flagName:    "log-level",
			viperKey:    "",
			expectedKey: "log-level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := newCobraCommand()
			flag := &cobraflags.Uint8Flag{
				Name:     tt.flagName,
				ViperKey: tt.viperKey,
				Value:    1,
				Usage:    "test flag",
			}

			flag.Register(cmd)

			// Test that getViperKey returns expected key
			actualValue := flag.GetUint8()             // This will trigger binding
			c.Assert(actualValue, qt.Equals, uint8(1)) // Default value

			// Test with flag set
			expectedValue := uint8(5)
			cmd.SetArgs([]string{"--" + tt.flagName, "5"})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)
			c.Assert(flag.GetUint8(), qt.Equals, expectedValue)

			// Test GetUint8E
			value, err := flag.GetUint8E()
			c.Assert(err, qt.IsNil)
			c.Assert(value, qt.Equals, expectedValue)
		})
	}
}
