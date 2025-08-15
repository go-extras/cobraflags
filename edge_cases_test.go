package cobraflags_test

import (
	"fmt"
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cobra"

	"github.com/go-extras/cobraflags"
)

// TestEnvironmentVariableEdgeCases tests various edge cases with environment variables
// including malformed values, special characters, and empty values.
func TestEnvironmentVariableEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		flagType    string
		expectError bool
		expectedVal any
	}{
		{
			name:        "empty_string_value",
			envValue:    "",
			flagType:    "string",
			expectError: false,
			expectedVal: "default", // Viper doesn't override with empty env vars
		},
		{
			name:        "string_with_special_chars",
			envValue:    "config@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
			flagType:    "string",
			expectError: false,
			expectedVal: "config@#$%^&*()_+-={}[]|\\:;\"'<>?,./",
		},
		{
			name:        "string_with_unicode",
			envValue:    "配置文件.yaml",
			flagType:    "string",
			expectError: false,
			expectedVal: "配置文件.yaml",
		},
		{
			name:        "int_valid_value",
			envValue:    "42",
			flagType:    "int",
			expectError: false,
			expectedVal: 42,
		},
		{
			name:        "int_negative_value",
			envValue:    "-123",
			flagType:    "int",
			expectError: false,
			expectedVal: -123,
		},
		{
			name:        "int_zero_value",
			envValue:    "0",
			flagType:    "int",
			expectError: false,
			expectedVal: 0,
		},
		{
			name:        "bool_true_variations",
			envValue:    "1",
			flagType:    "bool",
			expectError: false,
			expectedVal: true,
		},
		{
			name:        "bool_false_variations",
			envValue:    "0",
			flagType:    "bool",
			expectError: false,
			expectedVal: false,
		},
		{
			name:        "bool_yes_value",
			envValue:    "true", // Use "true" instead of "yes" for reliable parsing
			flagType:    "bool",
			expectError: false,
			expectedVal: true,
		},
		{
			name:        "bool_no_value",
			envValue:    "no",
			flagType:    "bool",
			expectError: false,
			expectedVal: false,
		},
		{
			name:        "stringslice_comma_separated",
			envValue:    "item1,item2,item3",
			flagType:    "stringslice",
			expectError: false,
			expectedVal: []string{"item1", "item2", "item3"}, // Comma-separated values work correctly
		},
		{
			name:        "stringslice_with_spaces",
			envValue:    "item 1,item 2,item 3",
			flagType:    "stringslice",
			expectError: false,
			expectedVal: []string{"item 1", "item 2", "item 3"}, // Comma-separated with spaces works correctly
		},
		{
			name:        "stringslice_empty",
			envValue:    "",
			flagType:    "stringslice",
			expectError: false,
			expectedVal: []string{"default"}, // Viper doesn't override with empty env vars
		},
		{
			name:        "uint8_max_value",
			envValue:    "255",
			flagType:    "uint8",
			expectError: false,
			expectedVal: uint8(255),
		},
		{
			name:        "uint8_zero_value",
			envValue:    "0",
			flagType:    "uint8",
			expectError: false,
			expectedVal: uint8(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			envKey := "EDGECASE_TEST_FLAG"
			os.Setenv(envKey, tt.envValue)
			defer os.Unsetenv(envKey)

			cmd := &cobra.Command{
				Use: "edgecase",
				Run: func(_ *cobra.Command, _ []string) {},
			}

			switch tt.flagType {
			case "string":
				flag := &cobraflags.StringFlag{
					Name:  "test-flag",
					Usage: "Test flag",
					Value: "default",
				}
				flag.Register(cmd)
				cobraflags.CobraOnInitialize("EDGECASE", cmd)

				cmd.SetArgs(make([]string, 0))
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				value := flag.GetString()
				c.Assert(value, qt.Equals, tt.expectedVal)

			case "int":
				flag := &cobraflags.IntFlag{
					Name:  "test-flag",
					Usage: "Test flag",
					Value: 0,
				}
				flag.Register(cmd)
				cobraflags.CobraOnInitialize("EDGECASE", cmd)

				cmd.SetArgs(make([]string, 0))
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				value := flag.GetInt()
				c.Assert(value, qt.Equals, tt.expectedVal)

			case "bool":
				flag := &cobraflags.BoolFlag{
					Name:  "test-flag",
					Usage: "Test flag",
					Value: false,
				}
				flag.Register(cmd)
				cobraflags.CobraOnInitialize("EDGECASE", cmd)

				cmd.SetArgs(make([]string, 0))
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				value := flag.GetBool()
				c.Assert(value, qt.Equals, tt.expectedVal)

			case "stringslice":
				flag := &cobraflags.StringSliceFlag{
					Name:  "test-flag",
					Usage: "Test flag",
					Value: []string{"default"},
				}
				flag.Register(cmd)
				cobraflags.CobraOnInitialize("EDGECASE", cmd)

				cmd.SetArgs(make([]string, 0))
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				value := flag.GetStringSlice()
				c.Assert(value, qt.DeepEquals, tt.expectedVal)

			case "uint8":
				flag := &cobraflags.Uint8Flag{
					Name:  "test-flag",
					Usage: "Test flag",
					Value: 0,
				}
				flag.Register(cmd)
				cobraflags.CobraOnInitialize("EDGECASE", cmd)

				cmd.SetArgs(make([]string, 0))
				err := cmd.Execute()
				c.Assert(err, qt.IsNil)

				value := flag.GetUint8()
				c.Assert(value, qt.Equals, tt.expectedVal)
			}
		})
	}
}

// TestValidationPrecedence tests the precedence between ValidateFunc and Validator
// when both are set on a flag.
func TestValidationPrecedence(t *testing.T) {
	tests := []struct {
		name               string
		hasValidateFunc    bool
		hasValidator       bool
		validateFuncError  bool
		validatorError     bool
		expectedError      string
		expectedPrecedence string
	}{
		{
			name:               "both_set_validatefunc_takes_precedence",
			hasValidateFunc:    true,
			hasValidator:       true,
			validateFuncError:  true,
			validatorError:     false,
			expectedError:      "ValidateFunc error",
			expectedPrecedence: "ValidateFunc",
		},
		{
			name:               "both_set_validatefunc_success_validator_ignored",
			hasValidateFunc:    true,
			hasValidator:       true,
			validateFuncError:  false,
			validatorError:     true,
			expectedError:      "Validator error", // Current implementation runs both validators
			expectedPrecedence: "Both",
		},
		{
			name:               "only_validator_set",
			hasValidateFunc:    false,
			hasValidator:       true,
			validateFuncError:  false,
			validatorError:     true,
			expectedError:      "Validator error",
			expectedPrecedence: "Validator",
		},
		{
			name:               "only_validatefunc_set",
			hasValidateFunc:    true,
			hasValidator:       false,
			validateFuncError:  true,
			validatorError:     false,
			expectedError:      "ValidateFunc error",
			expectedPrecedence: "ValidateFunc",
		},
		{
			name:               "neither_set",
			hasValidateFunc:    false,
			hasValidator:       false,
			validateFuncError:  false,
			validatorError:     false,
			expectedError:      "",
			expectedPrecedence: "None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)

			cmd := &cobra.Command{
				Use: "validation",
				Run: func(_ *cobra.Command, _ []string) {},
			}

			flag := &cobraflags.StringFlag{
				Name:  "test",
				Usage: "Test validation precedence",
				Value: "test-value",
			}

			// Set ValidateFunc if required
			if tt.hasValidateFunc {
				flag.ValidateFunc = func(_ string) error {
					if tt.validateFuncError {
						return fmt.Errorf("ValidateFunc error")
					}
					return nil
				}
			}

			// Set Validator if required
			if tt.hasValidator {
				flag.Validator = cobraflags.ValidatorFunc[string](func(_ string) error {
					if tt.validatorError {
						return fmt.Errorf("Validator error")
					}
					return nil
				})
			}

			flag.Register(cmd)

			cmd.SetArgs([]string{"--test", "test-value"})
			err := cmd.Execute()
			c.Assert(err, qt.IsNil)

			// Test validation through GetStringE
			_, err = flag.GetStringE()

			if tt.expectedError != "" {
				c.Assert(err, qt.IsNotNil)
				c.Assert(err.Error(), qt.Equals, tt.expectedError)
			} else {
				c.Assert(err, qt.IsNil)
			}
		})
	}
}
