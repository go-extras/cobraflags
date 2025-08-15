package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*IntFlag)(nil)

// IntFlag represents a command-line flag that accepts integer values.
// It provides automatic binding to environment variables via Viper and supports
// custom validation through ValidateFunc or Validator fields.
//
// IntFlag supports all standard flag features:
//   - Required flags (will cause command execution to fail if not provided)
//   - Persistent flags (available to subcommands)
//   - Shorthand notation (single character aliases)
//   - Custom Viper keys for configuration binding
//   - Validation with custom functions or validators
//
// Example usage:
//
//	portFlag := &IntFlag{
//		Name:      "port",
//		Shorthand: "p",
//		Usage:     "Server port number",
//		Value:     8080,
//		Required:  true,
//		ValidateFunc: func(port int) error {
//			if port < 1 || port > 65535 {
//				return fmt.Errorf("port must be between 1 and 65535")
//			}
//			return nil
//		},
//	}
//	portFlag.Register(cmd)
//
// Environment variable binding:
// With CobraOnInitialize("MYAPP", cmd), a flag named "port" will
// automatically bind to the environment variable "MYAPP_PORT".
type IntFlag FlagBase[int]

// pIntFlag is an alias for a pointer to FlagBase[int].
type pIntFlag = *FlagBase[int]

func (s *IntFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.Int(s.Name, s.Value, s.Usage)
	} else {
		flags.IntP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pIntFlag(s).getViperKey()}
}

// GetInt retrieves the current integer value of the flag.
// This method automatically binds the flag to Viper on first call and returns
// the value from Viper, which may come from command-line arguments, environment
// variables, or configuration files.
//
// Note: This method does NOT perform validation. Use GetIntE() if you need
// validation to be executed.
//
// Returns the integer value, which may be the default value if the flag was not set.
func (s *IntFlag) GetInt() int {
	viperKey := pIntFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetInt(viperKey)
}

// GetIntE retrieves the current integer value of the flag with validation.
// This method automatically binds the flag to Viper on first call, retrieves
// the value, and then applies any configured validation (ValidateFunc or Validator).
//
// Validation behavior:
//   - If ValidateFunc is set, it is called with the integer value
//   - If ValidateFunc is nil but Validator is set, Validator.Validate() is called
//   - If neither is set, no validation is performed
//
// Returns:
//   - On success: the integer value and nil error
//   - On validation failure: 0 and the validation error
//
// Use this method when you need to ensure the flag value meets your validation criteria.
func (s *IntFlag) GetIntE() (int, error) {
	viperKey := pIntFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetInt(viperKey)

	if result, err := pIntFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
