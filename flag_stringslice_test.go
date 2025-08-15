package cobraflags_test

import (
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"

	"github.com/go-extras/cobraflags"
)

func TestStringSliceFlag_Register(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:  "items",
		Value: []string{"default1", "default2"},
		Usage: "usage",
	}

	flag.Register(cmd)

	expectedValue := []string{"item1", "item2"}
	cmd.SetArgs([]string{"--items", "item1,item2"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, expectedValue)
}

func TestStringSliceFlag_GetStringSliceE(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:  "items",
		Value: []string{"default1", "default2"},
		Usage: "usage",
	}

	flag.Register(cmd)

	expectedValue := []string{"item1", "item2"}
	cmd.SetArgs([]string{"--items", "item1,item2"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	value, err := flag.GetStringSliceE()
	c.Assert(err, qt.IsNil)
	c.Assert(value, qt.DeepEquals, expectedValue)
}

func TestStringSliceFlag_WithShorthand(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:      "items",
		Value:     []string{"default1", "default2"},
		Usage:     "usage",
		Shorthand: "i",
	}

	flag.Register(cmd)

	expectedValue := []string{"item1", "item2"}
	cmd.SetArgs([]string{"-i", "item1,item2"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, expectedValue)
}

func TestStringSliceFlag_WithDefaultValue(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:  "items",
		Value: []string{"default1", "default2"},
		Usage: "usage",
	}

	flag.Register(cmd)

	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, []string{"default1", "default2"})
}

func TestStringSliceFlag_WithRequired(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:     "items",
		Value:    []string{"default1", "default2"},
		Usage:    "usage",
		Required: true,
	}

	flag.Register(cmd)

	// Test missing required flag
	cmd.SetArgs(make([]string, 0))
	err := cmd.Execute()

	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Equals, "required flag(s) \"items\" not set")

	// Test with required flag provided
	cmd.SetArgs([]string{"--items", "item1,item2"})
	err = cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, []string{"item1", "item2"})
}

func TestStringSliceFlag_ValidateFunc(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:  "items",
		Value: []string{"default1", "default2"},
		Usage: "usage",
		ValidateFunc: func(v []string) error {
			if len(v) == 0 {
				return fmt.Errorf("invalid value for flag %s", "items")
			}
			return nil
		},
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--items", ""})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetStringSliceE()
	c.Assert(err.Error(), qt.Equals, "invalid value for flag items")
}

func TestStringSliceFlag_Validator(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:  "items",
		Value: []string{"default1", "default2"},
		Usage: "usage",
		Validator: cobraflags.ValidatorFunc[[]string](func(v []string) error {
			if len(v) == 0 {
				return fmt.Errorf("invalid value for flag %s", "items")
			}
			return nil
		}),
	}

	flag.Register(cmd)

	cmd.SetArgs([]string{"--items", ""})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)

	_, err = flag.GetStringSliceE()
	c.Assert(err.Error(), qt.Equals, "invalid value for flag items")
}

func TestStringSliceFlag_WithPersistent(t *testing.T) {
	c := qt.New(t)

	cmd := newCobraCommand()
	flag := &cobraflags.StringSliceFlag{
		Name:       "items",
		Value:      []string{"default1", "default2"},
		Usage:      "usage",
		Persistent: true,
	}

	flag.Register(cmd)

	// Verify the flag is registered to PersistentFlags
	f := cmd.PersistentFlags().Lookup("items")
	c.Assert(f, qt.IsNotNil)

	expectedValue := []string{"item1", "item2"}
	cmd.SetArgs([]string{"--items", "item1,item2"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, expectedValue)
}

// TestStringSliceFlag_ViperKey_HappyPath tests ViperKey functionality with successful scenarios.
func TestStringSliceFlag_ViperKey_HappyPath(t *testing.T) {
	tests := []struct {
		name        string
		flagName    string
		viperKey    string
		expectedKey string
	}{
		{
			name:        "with_viper_key_set",
			flagName:    "allowed-hosts",
			viperKey:    "security.allowed.hosts",
			expectedKey: "security.allowed.hosts",
		},
		{
			name:        "with_empty_viper_key_fallback_to_name",
			flagName:    "allowed-hosts",
			viperKey:    "",
			expectedKey: "allowed-hosts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := newCobraCommand()
			flag := &cobraflags.StringSliceFlag{
				Name:     tt.flagName,
				ViperKey: tt.viperKey,
				Value:    []string{"localhost"},
				Usage:    "test flag",
			}

			flag.Register(cmd)

			// Test that getViperKey returns expected key
			actualValue := flag.GetStringSlice()                        // This will trigger binding
			c.Assert(actualValue, qt.DeepEquals, []string{"localhost"}) // Default value

			// Test with flag set
			expectedValue := []string{"host1", "host2"}
			cmd.SetArgs([]string{"--" + tt.flagName, "host1,host2"})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)
			c.Assert(flag.GetStringSlice(), qt.DeepEquals, expectedValue)

			// Test GetStringSliceE
			value, err := flag.GetStringSliceE()
			c.Assert(err, qt.IsNil)
			c.Assert(value, qt.DeepEquals, expectedValue)
		})
	}
}
