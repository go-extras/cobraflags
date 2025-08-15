package cobraflags

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*Uint8Flag)(nil)

// Uint8Flag represents a command-line flag that accepts unsigned 8-bit integer values (0-255).
// It provides automatic binding to environment variables via Viper and supports
// custom validation through ValidateFunc or Validator fields.
//
// Uint8Flag supports all standard flag features:
//   - Required flags (will cause command execution to fail if not provided)
//   - Persistent flags (available to subcommands)
//   - Shorthand notation (single character aliases)
//   - Custom Viper keys for configuration binding
//   - Validation with custom functions or validators
//
// Uint8 flags accept values in the range 0-255. Values outside this range
// will be automatically clamped by the underlying cast.ToUint8() function.
//
// Example usage:
//
//	priorityFlag := &Uint8Flag{
//		Name:      "priority",
//		Shorthand: "p",
//		Usage:     "Task priority level (0-255)",
//		Value:     128,
//		ValidateFunc: func(priority uint8) error {
//			if priority == 0 {
//				return fmt.Errorf("priority must be greater than 0")
//			}
//			return nil
//		},
//	}
//	priorityFlag.Register(cmd)
//
// Environment variable binding:
// With CobraOnInitialize("MYAPP", cmd), a flag named "priority" will
// automatically bind to the environment variable "MYAPP_PRIORITY".
type Uint8Flag FlagBase[uint8]

// pUint8Flag is an alias for a pointer to FlagBase[uint8].
type pUint8Flag = *FlagBase[uint8]

func (s *Uint8Flag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.Uint8(s.Name, s.Value, s.Usage)
	} else {
		flags.Uint8P(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pUint8Flag(s).getViperKey()}
}

// GetUint8 retrieves the current uint8 value of the flag.
// This method automatically binds the flag to Viper on first call and returns
// the value from Viper, which may come from command-line arguments, environment
// variables, or configuration files.
//
// Note: This method does NOT perform validation. Use GetUint8E() if you need
// validation to be executed.
//
// The value is retrieved as uint16 from Viper and then cast to uint8 using
// spf13/cast.ToUint8(), which handles overflow by clamping to the uint8 range.
//
// Returns the uint8 value, which may be the default value if the flag was not set.
func (s *Uint8Flag) GetUint8() uint8 {
	viperKey := pUint8Flag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return cast.ToUint8(viper.GetUint16(viperKey))
}

// GetUint8E retrieves the current uint8 value of the flag with validation.
// This method automatically binds the flag to Viper on first call, retrieves
// the value, and then applies any configured validation (ValidateFunc or Validator).
//
// Validation behavior:
//   - If ValidateFunc is set, it is called with the uint8 value
//   - If ValidateFunc is nil but Validator is set, Validator.Validate() is called
//   - If neither is set, no validation is performed
//
// The value is retrieved as uint16 from Viper and then cast to uint8 using
// spf13/cast.ToUint8(), which handles overflow by clamping to the uint8 range.
//
// Returns:
//   - On success: the uint8 value and nil error
//   - On validation failure: 0 and the validation error
//
// Use this method when you need to ensure the flag value meets your validation criteria.
func (s *Uint8Flag) GetUint8E() (uint8, error) {
	viperKey := pUint8Flag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	u16 := viper.GetUint16(viperKey)
	v := cast.ToUint8(u16)

	if result, err := pUint8Flag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
