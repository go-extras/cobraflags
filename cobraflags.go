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

// FlagBase is a generic base struct for all flag types that provides common functionality
// for flag registration, validation, and value retrieval. It uses Go generics to ensure
// type safety while sharing common behavior across different flag types.
//
// The validation system supports two approaches:
//   - ValidateFunc: A simple function that takes the flag's value type and returns an error
//   - Validator: An interface-based validator that can be reused across different flag types
//
// When both ValidateFunc and Validator are set, ValidateFunc takes precedence and Validator is ignored.
//
// The ViperKey field allows using different configuration keys than flag names for Viper binding.
// If ViperKey is empty, the flag will fall back to using its Name for Viper binding.
// This enables:
//   - Using different configuration keys than flag names
//   - Supporting nested configuration structures (e.g., "app.config.file")
//   - Maintaining backward compatibility when renaming flags
//
// Example usage:
//
//	flag := &StringFlag{
//		Name:     "config-file",
//		ViperKey: "app.config.file", // Custom Viper key
//		Usage:    "Path to configuration file",
//		Value:    "default.yaml",
//		ValidateFunc: func(path string) error {
//			if !strings.HasSuffix(path, ".yaml") {
//				return fmt.Errorf("config file must be a YAML file")
//			}
//			return nil
//		},
//	}
type FlagBase[T any] struct {
	Name         string        // Flag name used for command line arguments
	ViperKey     string        // Custom Viper configuration key (falls back to Name if empty)
	Shorthand    string        // Single character shorthand for the flag
	Usage        string        // Help text for the flag
	Required     bool          // Whether the flag is required
	Persistent   bool          // Whether the flag is persistent across subcommands
	Value        T             // Default value
	ValidateFunc func(T) error // Custom validation function (takes precedence over Validator)
	Validator    Validator     // Custom validator implementing the Validator interface

	flag     *pflag.Flag
	bindOnce sync.Once

	flagGetter
	flagGetterE
}

// validate applies custom validation logic if defined and returns the value or an error if validation fails.
//
// Validation precedence (in order):
//  1. ValidateFunc - if set, this function is called and Validator is ignored
//  2. Validator - if set and ValidateFunc is nil, the Validate method is called
//  3. No validation - if neither is set, the value is returned as-is
//
// Returns:
//   - On success: the original value and nil error
//   - On validation failure: zero value of type T and the validation error
//
// This method is called internally by GetE methods to ensure validation
// occurs before returning values to the caller.
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

// getViperKey returns the Viper configuration key to use for this flag.
//
// Behavior:
//   - If ViperKey is set (non-empty), returns ViperKey
//   - If ViperKey is empty, falls back to using Name
//
// This allows flags to use different configuration keys than their command-line names,
// enabling nested configuration structures and backward compatibility.
//
// Example:
//
//	Flag with Name="config-file" and ViperKey="app.config.file"
//	will bind to the "app.config.file" key in Viper instead of "config-file".
func (s *FlagBase[T]) getViperKey() string {
	if s.ViperKey != "" {
		return s.ViperKey
	}
	return s.Name
}

// Register registers multiple flags with the given cobra command in a single call.
// This is a convenience function that calls Register() on each flag individually.
//
// Example:
//
//	countFlag := &IntFlag{Name: "count", Value: 10}
//	verboseFlag := &BoolFlag{Name: "verbose", Value: false}
//	Register(cmd, countFlag, verboseFlag)
func Register(cmd *cobra.Command, flags ...Flag) {
	for _, flag := range flags {
		flag.Register(cmd)
	}
}

// RegisterMap registers flags from a map with the given cobra command.
// The map keys are ignored; only the flag values are registered.
// This is useful when you have flags organized in a map structure.
//
// Example:
//
//	flags := map[string]Flag{
//		"count":   &IntFlag{Name: "count", Value: 10},
//		"verbose": &BoolFlag{Name: "verbose", Value: false},
//	}
//	RegisterMap(cmd, flags)
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
