package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*StringSliceFlag)(nil)

// StringSliceFlag represents a command-line flag that accepts multiple string values.
// It provides automatic binding to environment variables via Viper and supports
// custom validation through ValidateFunc or Validator fields.
//
// StringSliceFlag supports all standard flag features:
//   - Required flags (will cause command execution to fail if not provided)
//   - Persistent flags (available to subcommands)
//   - Shorthand notation (single character aliases)
//   - Custom Viper keys for configuration binding
//   - Validation with custom functions or validators
//
// String slice flags accept multiple values in several ways:
//   - Multiple flag instances: --item value1 --item value2
//   - Comma-separated values: --item value1,value2,value3
//   - Environment variables as comma-separated strings
//
// Example usage:
//
//	tagsFlag := &StringSliceFlag{
//		Name:      "tags",
//		Shorthand: "t",
//		Usage:     "Tags to apply (can be specified multiple times)",
//		Value:     []string{"default"},
//		ValidateFunc: func(tags []string) error {
//			if len(tags) == 0 {
//				return fmt.Errorf("at least one tag must be specified")
//			}
//			return nil
//		},
//	}
//	tagsFlag.Register(cmd)
//
// Environment variable binding:
// With CobraOnInitialize("MYAPP", cmd), a flag named "tags" will
// automatically bind to the environment variable "MYAPP_TAGS".
type StringSliceFlag FlagBase[[]string]

// pStringSliceFlag is an alias for a pointer to FlagBase[[]string].
type pStringSliceFlag = *FlagBase[[]string]

func (s *StringSliceFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.StringSlice(s.Name, s.Value, s.Usage)
	} else {
		flags.StringSliceP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pStringSliceFlag(s).getViperKey()}
}

// GetStringSlice retrieves the current string slice value of the flag.
// This method automatically binds the flag to Viper on first call and returns
// the value from Viper, which may come from command-line arguments, environment
// variables, or configuration files.
//
// Note: This method does NOT perform validation. Use GetStringSliceE() if you need
// validation to be executed.
//
// Returns the string slice value, which may be the default value if the flag was not set.
func (s *StringSliceFlag) GetStringSlice() []string {
	viperKey := pStringSliceFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetStringSlice(viperKey)
}

// GetStringSliceE retrieves the current string slice value of the flag with validation.
// This method automatically binds the flag to Viper on first call, retrieves
// the value, and then applies any configured validation (ValidateFunc or Validator).
//
// Validation behavior:
//   - If ValidateFunc is set, it is called with the string slice value
//   - If ValidateFunc is nil but Validator is set, Validator.Validate() is called
//   - If neither is set, no validation is performed
//
// Returns:
//   - On success: the string slice value and nil error
//   - On validation failure: nil slice and the validation error
//
// Use this method when you need to ensure the flag value meets your validation criteria.
func (s *StringSliceFlag) GetStringSliceE() ([]string, error) {
	viperKey := pStringSliceFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetStringSlice(viperKey)

	if result, err := pStringSliceFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
