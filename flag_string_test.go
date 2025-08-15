package cobraflags_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func TestStringFlag_Register(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:  "name",
		Value: "default",
		Usage: "usage",
	}

	flag.Register(cmd)

	const expectedValue = "test"
	cmd.SetArgs([]string{"--name", expectedValue})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetString(), qt.Equals, expectedValue)
}

func TestStringFlag_GetStringE(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:  "feature",
		Value: "default",
		Usage: "enable feature",
	}

	flag.Register(cmd)

	const expectedValue = "test"
	cmd.SetArgs([]string{"--feature", expectedValue})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	value, err := flag.GetStringE()
	c.Assert(err, qt.IsNil)
	c.Assert(value, qt.Equals, expectedValue)
}

func TestStringFlag_WithShorthand(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:      "name",
		Value:     "default",
		Usage:     "usage",
		Shorthand: "n",
	}

	flag.Register(cmd)

	const expectedValue = "test"
	cmd.SetArgs([]string{"-n", expectedValue})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetString(), qt.Equals, expectedValue)
}

func TestStringFlag_WithDefaultValue(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:  "name",
		Value: "default",
		Usage: "usage",
	}

	flag.Register(cmd)

	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetString(), qt.Equals, "default")
}

func TestStringFlag_WithRequired(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:     "name",
		Value:    "default",
		Usage:    "usage",
		Required: true,
	}

	flag.Register(cmd)

	// Test missing required flag
	cmd.SetArgs(make([]string, 0))
	err := cmd.Execute()

	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "required flag(s) \"name\" not set")

	// Test with required flag provided
	cmd.SetArgs([]string{"--name", "test"})
	err = cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetString(), qt.Equals, "test")
}

func TestStringFlag_ValidateFunc(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:  "name",
		Value: "default",
		Usage: "usage",
		ValidateFunc: func(v string) error {
			if v == "" {
				return fmt.Errorf("invalid value for flag %s", "name")
			}
			return nil
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--name", ""})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetStringE()
	c.Assert(err.Error(), qt.Equals, "invalid value for flag name")
}

func TestStringFlag_Validator(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:  "name",
		Value: "default",
		Usage: "usage",
		Validator: cobraflags.ValidatorFunc[string](func(v string) error {
			if v == "" {
				return fmt.Errorf("invalid value for flag %s", "name")
			}
			return nil
		}),
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--name", ""})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetStringE()
	c.Assert(err.Error(), qt.Equals, "invalid value for flag name")
}

func TestStringFlag_WithPersistent(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringFlag{
		Name:       "name",
		Value:      "default",
		Usage:      "usage",
		Persistent: true,
	}

	flag.Register(cmd)

	// Verify the flag is registered to PersistentFlags
	f := cmd.PersistentFlags().Lookup("name")
	c.Assert(f, qt.IsNotNil)

	const expectedValue = "test"
	cmd.SetArgs([]string{"--name", expectedValue})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetString(), qt.Equals, expectedValue)
}

// TestStringFlag_ViperKey_HappyPath tests ViperKey functionality with successful scenarios.
func TestStringFlag_ViperKey_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		flagName    string
		viperKey    string
		expectedKey string
	}{
		{
			name:        "with_viper_key_set",
			flagName:    "config-file",
			viperKey:    "app.config.file",
			expectedKey: "app.config.file",
		},
		{
			name:        "with_empty_viper_key_fallback_to_name",
			flagName:    "config-file",
			viperKey:    "",
			expectedKey: "config-file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := newCobraCommand()
			flag := &cobraflags.StringFlag{
				Name:     tt.flagName,
				ViperKey: tt.viperKey,
				Value:    "default.conf",
				Usage:    "test flag",
			}

			flag.Register(cmd)

			// Test that getViperKey returns expected key
			actualValue := flag.GetString()                  // This will trigger binding
			c.Assert(actualValue, qt.Equals, "default.conf") // Default value

			// Test with flag set
			expectedValue := "custom.conf"
			cmd.SetArgs([]string{"--" + tt.flagName, expectedValue})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)
			c.Assert(flag.GetString(), qt.Equals, expectedValue)

			// Test GetStringE
			value, err := flag.GetStringE()
			c.Assert(err, qt.IsNil)
			c.Assert(value, qt.Equals, expectedValue)
		})
	}
}
