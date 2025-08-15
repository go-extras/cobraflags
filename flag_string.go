package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*StringFlag)(nil)

// StringFlag represents a command-line flag that accepts string values.
// It provides automatic binding to environment variables via Viper and supports
// custom validation through ValidateFunc or Validator fields.
//
// StringFlag supports all standard flag features:
//   - Required flags (will cause command execution to fail if not provided)
//   - Persistent flags (available to subcommands)
//   - Shorthand notation (single character aliases)
//   - Custom Viper keys for configuration binding
//   - Validation with custom functions or validators
//
// Example usage:
//
//	configFlag := &StringFlag{
//		Name:      "config",
//		Shorthand: "c",
//		Usage:     "Path to configuration file",
//		Value:     "config.yaml",
//		Required:  true,
//		ValidateFunc: func(path string) error {
//			if !strings.HasSuffix(path, ".yaml") {
//				return fmt.Errorf("config file must be a YAML file")
//			}
//			return nil
//		},
//	}
//	configFlag.Register(cmd)
//
// Environment variable binding:
// With CobraOnInitialize("MYAPP", cmd), a flag named "config-file" will
// automatically bind to the environment variable "MYAPP_CONFIG_FILE".
type StringFlag FlagBase[string]

// pStringFlag is an alias for a pointer to FlagBase[string].
type pStringFlag = *FlagBase[string]

func (s *StringFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.String(s.Name, s.Value, s.Usage)
	} else {
		flags.StringP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pStringFlag(s).getViperKey()}
}

// GetString retrieves the current string value of the flag.
// This method automatically binds the flag to Viper on first call and returns
// the value from Viper, which may come from command-line arguments, environment
// variables, or configuration files.
//
// Note: This method does NOT perform validation. Use GetStringE() if you need
// validation to be executed.
//
// Returns the string value, which may be the default value if the flag was not set.
func (s *StringFlag) GetString() string {
	viperKey := pStringFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetString(viperKey)
}

// GetStringE retrieves the current string value of the flag with validation.
// This method automatically binds the flag to Viper on first call, retrieves
// the value, and then applies any configured validation (ValidateFunc or Validator).
//
// Validation behavior:
//   - If ValidateFunc is set, it is called with the string value
//   - If ValidateFunc is nil but Validator is set, Validator.Validate() is called
//   - If neither is set, no validation is performed
//
// Returns:
//   - On success: the string value and nil error
//   - On validation failure: empty string and the validation error
//
// Use this method when you need to ensure the flag value meets your validation criteria.
func (s *StringFlag) GetStringE() (string, error) {
	viperKey := pStringFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetString(viperKey)

	if result, err := pStringFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
