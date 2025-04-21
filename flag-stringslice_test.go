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
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	c.Assert(err, qt.Not(qt.IsNil))
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
	c.Assert(f, qt.Not(qt.IsNil))

	expectedValue := []string{"item1", "item2"}
	cmd.SetArgs([]string{"--items", "item1,item2"})
	err := cmd.Execute()

	c.Assert(err, qt.IsNil)
	c.Assert(flag.GetStringSlice(), qt.DeepEquals, expectedValue)
}
