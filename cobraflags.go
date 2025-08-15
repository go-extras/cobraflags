package cobraflags

import (
	"log/slog"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const viperKeyAnnotation = "viper-key"

// flagGetter is an interface for getting flag values.
type flagGetter interface {
	GetString() string
	GetBool() bool
	GetInt() int
	GetUint8() uint8
	GetStringSlice() []string
}

// flagGetterE is an interface for getting flag values together with validation.
type flagGetterE interface {
	GetStringE() (string, error)
	GetBoolE() (bool, error)
	GetIntE() (int, error)
	GetUint8E() (uint8, error)
	GetStringSliceE() ([]string, error)
}

// Flag is an interface for a flag that can be registered with a cobra command.
type Flag interface {
	// Register registers the flag with the given cobra command.
	Register(*cobra.Command)

	flagGetter
	flagGetterE
}

// FlagBase is a base struct for flags.
type FlagBase[T any] struct {
	Name         string        // Flag name used for command line arguments
	ViperKey     string        // Custom Viper configuration key (falls back to Name if empty)
	Shorthand    string        // Single character shorthand for the flag
	Usage        string        // Help text for the flag
	Required     bool          // Whether the flag is required
	Persistent   bool          // Whether the flag is persistent across subcommands
	Value        T             // Default value
	ValidateFunc func(T) error // Custom validation function
	Validator    Validator     // Custom validator implementing the Validator interface

	flag     *pflag.Flag
	bindOnce sync.Once

	flagGetter
	flagGetterE
}

// validate applies custom validation logic if defined and returns the value or an error if validation fails.
// If no custom validation is defined, it returns the value and nil.
// If both ValidateFunc and Validator are defined, ValidateFunc takes precedence.
// If validation fails, it returns the zero value of the type and the error.
func (s *FlagBase[T]) validate(v T) (result T, err error) {
	if s.ValidateFunc != nil {
		err = s.ValidateFunc(v)
		if err != nil {
			return result, err
		}
	}

	if s.Validator != nil {
		err = s.Validator.Validate(v)
		if err != nil {
			return result, err
		}
	}

	return v, nil
}

// getViperKey returns the Viper key to use, falling back to Name if ViperKey is empty.
func (s *FlagBase[T]) getViperKey() string {
	if s.ViperKey != "" {
		return s.ViperKey
	}
	return s.Name
}

// Register registers the given flags with the given cobra command.
func Register(cmd *cobra.Command, flags ...Flag) {
	for _, flag := range flags {
		flag.Register(cmd)
	}
}

// RegisterMap registers the given flags with the given cobra command.
func RegisterMap(cmd *cobra.Command, flags map[string]Flag) {
	for _, flag := range flags {
		flag.Register(cmd)
	}
}

func noError(err error) {
	if err != nil {
		slog.With("error", err).Error("unexpected error")
		panic(err)
	}
}
