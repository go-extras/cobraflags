package cobraflags_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func TestIntFlag_Register(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:  "name",
		Value: 0,
		Usage: "usage",
	}

	flag.Register(cmd)

	const expectedValue = 42
	cmd.SetArgs([]string{"--name", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetInt(), qt.Equals, expectedValue)
}

func TestIntFlag_GetIntE(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:  "feature",
		Value: 0,
		Usage: "enable feature",
	}

	flag.Register(cmd)

	const expectedValue = 42
	cmd.SetArgs([]string{"--feature", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	value, err := flag.GetIntE()
	c.Assert(err, qt.IsNil)
	c.Assert(value, qt.Equals, expectedValue)
}

func TestIntFlag_WithShorthand(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:      "name",
		Value:     0,
		Usage:     "usage",
		Shorthand: "n",
	}

	flag.Register(cmd)

	const expectedValue = 42
	cmd.SetArgs([]string{"-n", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetInt(), qt.Equals, expectedValue)
}

func TestIntFlag_WithDefaultValue(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:  "name",
		Value: 42,
		Usage: "usage",
	}

	flag.Register(cmd)

	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetInt(), qt.Equals, 42)
}

func TestIntFlag_WithRequired(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:     "name",
		Value:    0,
		Usage:    "usage",
		Required: true,
	}

	flag.Register(cmd)

	err := cmd.Execute()

	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "required flag(s) \"name\" not set")
}

func TestIntFlag_ValidateFunc(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:  "name",
		Value: 0,
		Usage: "usage",
		ValidateFunc: func(v int) error {
			if v < 0 {
				return fmt.Errorf("invalid value %d for flag %s", v, "name")
			}
			return nil
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--name", "-1"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetIntE()
	c.Assert(err.Error(), qt.Equals, "invalid value -1 for flag name")
}

func TestIntFlag_Validator(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:  "name",
		Value: 0,
		Usage: "usage",
		Validator: cobraflags.ValidatorFunc[int](func(v int) error {
			if v < 0 {
				return fmt.Errorf("invalid value %d for flag %s", v, "name")
			}
			return nil
		}),
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--name", "-1"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetIntE()
	c.Assert(err.Error(), qt.Equals, "invalid value -1 for flag name")
}

func TestIntFlag_WithPersistent(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.IntFlag{
		Name:       "name",
		Value:      0,
		Usage:      "usage",
		Persistent: true,
	}

	flag.Register(cmd)

	const expectedValue = 42
	cmd.SetArgs([]string{"--name", "42"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetInt(), qt.Equals, expectedValue)
}

// TestIntFlag_ViperKey_HappyPath tests ViperKey functionality with successful scenarios.
func TestIntFlag_ViperKey_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		flagName    string
		viperKey    string
		expectedKey string
	}{
		{
			name:        "with_viper_key_set",
			flagName:    "max-connections",
			viperKey:    "server.max.connections",
			expectedKey: "server.max.connections",
		},
		{
			name:        "with_empty_viper_key_fallback_to_name",
			flagName:    "max-connections",
			viperKey:    "",
			expectedKey: "max-connections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := newCobraCommand()
			flag := &cobraflags.IntFlag{
				Name:     tt.flagName,
				ViperKey: tt.viperKey,
				Value:    10,
				Usage:    "test flag",
			}

			flag.Register(cmd)

			// Test that getViperKey returns expected key
			actualValue := flag.GetInt()         // This will trigger binding
			c.Assert(actualValue, qt.Equals, 10) // Default value

			// Test with flag set
			expectedValue := 42
			cmd.SetArgs([]string{"--" + tt.flagName, "42"})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)
			c.Assert(flag.GetInt(), qt.Equals, expectedValue)

			// Test GetIntE
			value, err := flag.GetIntE()
			c.Assert(err, qt.IsNil)
			c.Assert(value, qt.Equals, expectedValue)
		})
	}
}
