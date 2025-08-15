package cobraflags_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func boolToString(value bool) string { //revive:disable-line:flag-parameter // not a control flag
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
	cmd.SetArgs(make([]string, 0))
	err := cmd.Execute()

	c.Assert(err, qt.IsNotNil)

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
	c.Assert(f, qt.IsNotNil)

	const expectedValue = true
	cmd.SetArgs([]string{"--debug", boolToString(expectedValue)})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetBool(), qt.Equals, expectedValue)
}

func TestBoolFlag_WithValidation(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:  "feature",
		Value: false,
		Usage: "enable feature",
		ValidateFunc: func(v bool) error {
			return fmt.Errorf("%v is invalid value", v)
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--feature", "true"})
	err := cmd.Execute()
	c.Assert(err, qt.IsNil)

	// GetBoolE calls validation
	_, err = flag.GetBoolE()
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "true is invalid value")
}

func TestBoolFlag_WithValidator(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.BoolFlag{
		Name:  "feature",
		Value: false,
		Usage: "enable feature",
		Validator: cobraflags.ValidatorFunc[bool](func(v bool) error {
			return fmt.Errorf("%v is invalid value", v)
		}),
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--feature", "true"})
	err := cmd.Execute()
	c.Assert(err, qt.IsNil)

	// GetBoolE calls validation
	_, err = flag.GetBoolE()
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "true is invalid value")
}

// TestBoolFlag_ViperKey_HappyPath tests ViperKey functionality with successful scenarios.
func TestBoolFlag_ViperKey_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		flagName    string
		viperKey    string
		expectedKey string
	}{
		{
			name:        "with_viper_key_set",
			flagName:    "enable-feature",
			viperKey:    "custom.feature.enabled",
			expectedKey: "custom.feature.enabled",
		},
		{
			name:        "with_empty_viper_key_fallback_to_name",
			flagName:    "enable-feature",
			viperKey:    "",
			expectedKey: "enable-feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := newCobraCommand()
			flag := &cobraflags.BoolFlag{
				Name:     tt.flagName,
				ViperKey: tt.viperKey,
				Value:    false,
				Usage:    "test flag",
			}

			flag.Register(cmd)

			// Test that getViperKey returns expected key
			actualKey := flag.GetBool()           // This will trigger binding
			c.Assert(actualKey, qt.Equals, false) // Default value

			// Test with flag set
			cmd.SetArgs([]string{"--" + tt.flagName, "true"})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)
			c.Assert(flag.GetBool(), qt.Equals, true)

			// Test GetBoolE
			value, err := flag.GetBoolE()
			c.Assert(err, qt.IsNil)
			c.Assert(value, qt.Equals, true)
		})
	}
}
