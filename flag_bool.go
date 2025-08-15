package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*BoolFlag)(nil)

// BoolFlag represents a command-line flag that accepts boolean values.
// It provides automatic binding to environment variables via Viper and supports
// custom validation through ValidateFunc or Validator fields.
//
// BoolFlag supports all standard flag features:
//   - Required flags (will cause command execution to fail if not provided)
//   - Persistent flags (available to subcommands)
//   - Shorthand notation (single character aliases)
//   - Custom Viper keys for configuration binding
//   - Validation with custom functions or validators
//
// Boolean flags have special behavior:
//   - They can be used without a value: --verbose (sets to true)
//   - They can be explicitly set: --verbose=true or --verbose=false
//   - Environment variables accept: true/false, 1/0, yes/no, on/off
//
// Example usage:
//
//	verboseFlag := &BoolFlag{
//		Name:      "verbose",
//		Shorthand: "v",
//		Usage:     "Enable verbose output",
//		Value:     false,
//		ValidateFunc: func(verbose bool) error {
//			// Custom validation logic if needed
//			return nil
//		},
//	}
//	verboseFlag.Register(cmd)
//
// Environment variable binding:
// With CobraOnInitialize("MYAPP", cmd), a flag named "verbose" will
// automatically bind to the environment variable "MYAPP_VERBOSE".
type BoolFlag FlagBase[bool]

// pBoolFlag is an alias for a pointer to FlagBase[bool].
type pBoolFlag = *FlagBase[bool]

func (s *BoolFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	flags.BoolP(s.Name, s.Shorthand, s.Value, s.Usage)

	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pBoolFlag(s).getViperKey()}
}

// GetBool retrieves the current boolean value of the flag.
// This method automatically binds the flag to Viper on first call and returns
// the value from Viper, which may come from command-line arguments, environment
// variables, or configuration files.
//
// Note: This method does NOT perform validation. Use GetBoolE() if you need
// validation to be executed.
//
// Returns the boolean value, which may be the default value if the flag was not set.
func (s *BoolFlag) GetBool() bool {
	viperKey := pBoolFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetBool(viperKey)
}

// GetBoolE retrieves the current boolean value of the flag with validation.
// This method automatically binds the flag to Viper on first call, retrieves
// the value, and then applies any configured validation (ValidateFunc or Validator).
//
// Validation behavior:
//   - If ValidateFunc is set, it is called with the boolean value
//   - If ValidateFunc is nil but Validator is set, Validator.Validate() is called
//   - If neither is set, no validation is performed
//
// Returns:
//   - On success: the boolean value and nil error
//   - On validation failure: false and the validation error
//
// Use this method when you need to ensure the flag value meets your validation criteria.
func (s *BoolFlag) GetBoolE() (bool, error) {
	viperKey := pBoolFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetBool(viperKey)

	if result, err := pBoolFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
