package cobraflags_test

import (
	"errors"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cobra"

	"github.com/go-extras/cobraflags"
)

// TestNoErrorFunction tests the behavior of the noError function
// Note: Since noError panics on error, we test scenarios where it should NOT panic
func TestNoErrorFunction_HappyPath(t *testing.T) {
	c := qt.New(t)

	// Test that normal flag registration doesn't cause panics
	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	// Test various flag types to ensure noError doesn't panic during normal operation
	stringFlag := &cobraflags.StringFlag{
		Name:  "string-flag",
		Usage: "String flag",
		Value: "default",
	}

	intFlag := &cobraflags.IntFlag{
		Name:  "int-flag",
		Usage: "Int flag",
		Value: 42,
	}

	boolFlag := &cobraflags.BoolFlag{
		Name:  "bool-flag",
		Usage: "Bool flag",
		Value: false,
	}

	stringSliceFlag := &cobraflags.StringSliceFlag{
		Name:  "slice-flag",
		Usage: "String slice flag",
		Value: []string{"default"},
	}

	uint8Flag := &cobraflags.Uint8Flag{
		Name:  "uint8-flag",
		Usage: "Uint8 flag",
		Value: 128,
	}

	// Register all flags - this should not panic
	c.Assert(func() {
		stringFlag.Register(cmd)
		intFlag.Register(cmd)
		boolFlag.Register(cmd)
		stringSliceFlag.Register(cmd)
		uint8Flag.Register(cmd)
	}, qt.Not(qt.PanicMatches), ".*")

	// Execute command - this should not panic
	cmd.SetArgs([]string{
		"--string-flag", "test",
		"--int-flag", "100",
		"--bool-flag",
		"--slice-flag", "item1,item2",
		"--uint8-flag", "200",
	})

	c.Assert(func() {
		err := cmd.Execute()
		c.Assert(err, qt.IsNil)
	}, qt.Not(qt.PanicMatches), ".*")

	// Verify values are accessible without panic
	c.Assert(func() {
		_ = stringFlag.GetString()
		_ = intFlag.GetInt()
		_ = boolFlag.GetBool()
		_ = stringSliceFlag.GetStringSlice()
		_ = uint8Flag.GetUint8()
	}, qt.Not(qt.PanicMatches), ".*")
}

// TestFlagRegistrationErrors tests error conditions during flag registration
func TestFlagRegistrationErrors(t *testing.T) {
	c := qt.New(t)

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	// Test duplicate flag registration
	flag1 := &cobraflags.StringFlag{
		Name:  "duplicate",
		Usage: "First flag",
		Value: "first",
	}

	flag2 := &cobraflags.StringFlag{
		Name:  "duplicate",
		Usage: "Second flag",
		Value: "second",
	}

	// Register first flag
	flag1.Register(cmd)

	// Registering second flag with same name should cause a panic due to pflag error
	c.Assert(func() {
		flag2.Register(cmd)
	}, qt.PanicMatches, ".*")
}

// TestRequiredFlagErrors tests error conditions with required flags
func TestRequiredFlagErrors(t *testing.T) {
	c := qt.New(t)

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	// Test required flag that is not provided
	requiredFlag := &cobraflags.StringFlag{
		Name:     "required",
		Usage:    "Required flag",
		Value:    "default",
		Required: true,
	}

	requiredFlag.Register(cmd)

	// Execute without providing required flag should return error
	cmd.SetArgs(make([]string, 0))
	err := cmd.Execute()
	c.Assert(err, qt.IsNotNil)
	c.Assert(err.Error(), qt.Matches, ".*required.*")
}

// TestValidationErrors tests various validation error scenarios
func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		flagType      string
		value         string
		expectedError string
	}{
		{
			name:          "string_validation_error",
			flagType:      "string",
			value:         "invalid",
			expectedError: "string validation failed",
		},
		{
			name:          "int_validation_error",
			flagType:      "int",
			value:         "-1",
			expectedError: "int validation failed",
		},
		{
			name:          "bool_validation_error",
			flagType:      "bool",
			value:         "true",
			expectedError: "bool validation failed",
		},
		{
			name:          "stringslice_validation_error",
			flagType:      "stringslice",
			value:         "item1,item2",
			expectedError: "stringslice validation failed",
		},
		{
			name:          "uint8_validation_error",
			flagType:      "uint8",
			value:         "100",
			expectedError: "uint8 validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := &cobra.Command{
				Use: "test",
				Run: func(_ *cobra.Command, _ []string) {},
			}

			switch tt.flagType {
			case "string":
				flag := &cobraflags.StringFlag{
					Name:  "test",
					Usage: "Test flag",
					Value: "default",
					ValidateFunc: func(_ string) error {
						return errors.New(tt.expectedError)
					},
				}
				flag.Register(cmd)
				cmd.SetArgs([]string{"--test", tt.value})
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				_, err = flag.GetStringE()
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)

			case "int":
				flag := &cobraflags.IntFlag{
					Name:  "test",
					Usage: "Test flag",
					Value: 0,
					ValidateFunc: func(_ int) error {
						return errors.New(tt.expectedError)
					},
				}
				flag.Register(cmd)
				cmd.SetArgs([]string{"--test", tt.value})
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				_, err = flag.GetIntE()
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)

			case "bool":
				flag := &cobraflags.BoolFlag{
					Name:  "test",
					Usage: "Test flag",
					Value: false,
					ValidateFunc: func(_ bool) error {
						return errors.New(tt.expectedError)
					},
				}
				flag.Register(cmd)
				cmd.SetArgs([]string{"--test=" + tt.value})
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				_, err = flag.GetBoolE()
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)

			case "stringslice":
				flag := &cobraflags.StringSliceFlag{
					Name:  "test",
					Usage: "Test flag",
					Value: []string{"default"},
					ValidateFunc: func(_ []string) error {
						return errors.New(tt.expectedError)
					},
				}
				flag.Register(cmd)
				cmd.SetArgs([]string{"--test", tt.value})
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				_, err = flag.GetStringSliceE()
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)

			case "uint8":
				flag := &cobraflags.Uint8Flag{
					Name:  "test",
					Usage: "Test flag",
					Value: 0,
					ValidateFunc: func(_ uint8) error {
						return errors.New(tt.expectedError)
					},
				}
				flag.Register(cmd)
				cmd.SetArgs([]string{"--test", tt.value})
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				_, err = flag.GetUint8E()
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)
			}
		})
	}
}

// TestViperKeyEdgeCases tests edge cases with custom Viper keys
func TestViperKeyEdgeCases(t *testing.T) {
	c := qt.New(t)

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	// Test flag with nested Viper key
	flag := &cobraflags.StringFlag{
		Name:     "config-file",
		ViperKey: "app.config.file.path",
		Usage:    "Configuration file path",
		Value:    "default.yaml",
	}

	flag.Register(cmd)

	// Verify that the flag works with nested Viper key
	cmd.SetArgs([]string{"--config-file", "custom.yaml"})
	err := cmd.Execute()
	c.Assert(err, qt.IsNil)

	value := flag.GetString()
	c.Assert(value, qt.Equals, "custom.yaml")

	// Test GetStringE as well
	valueE, err := flag.GetStringE()
	c.Assert(err, qt.IsNil)
	c.Assert(valueE, qt.Equals, "custom.yaml")
}
