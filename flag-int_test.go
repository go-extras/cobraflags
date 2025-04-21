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

	c.Assert(err, qt.Not(qt.IsNil))
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
