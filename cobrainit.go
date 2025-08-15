package cobraflags

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var initOnce sync.Once

var noEnvFlags = map[string]bool{
	"help": true,
}

// CobraOnInitialize initializes Cobra command(s) with configurations
// derived from environment variables. It sets up Viper to automatically
// detect and bind these variables based on the provided environment
// variable prefix. This function should be called before executing
// the Cobra command to ensure all configurations are properly loaded.
//
// Parameters:
// - envPrefix: A string prefix that environment variables must have to be considered by Viper.
// - commands: A slice of Cobra commands to be initialized.
//
// This function also ensures that all flags for the provided commands
// are preset with values from the corresponding environment variables if they exist.
func CobraOnInitialize(envPrefix string, command *cobra.Command) {
	cobraInit := func() {
		initOnce.Do(func() {
			visited := make(map[*pflag.Flag]bool)
			viper.AutomaticEnv()                          // Enable automatic detection of environment variables.
			viper.SetEnvPrefix(envPrefix)                 // Set the prefix for environment variables.
			replacer := strings.NewReplacer("-", "_")     // Create a replacer for environment variable names.
			viper.SetEnvKeyReplacer(replacer)             // Set the replacer for Viper.
			PostInitCommands(envPrefix, visited, command) // Initialize commands with environment variable values.
		})
	}

	fn := command.HelpFunc()
	command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		cobraInit()
		fn(cmd, args)
	})

	cobra.OnInitialize(cobraInit)
}

// PostInitCommands iterates through the given slice of Cobra commands
// and recursively initializes them and their subcommands. This includes
// binding each command's flags to corresponding environment variables
// using Viper.
//
// Parameters:
// - commands: A slice of Cobra commands to be initialized.
//
// This function is called recursively for each command that contains subcommands,
// ensuring that the entire command tree is covered.
func PostInitCommands(envPrefix string, flags map[*pflag.Flag]bool, commands ...*cobra.Command) {
	for _, cmd := range commands {
		PresetRequiredFlags(envPrefix, flags, cmd) // Bind environment variables to command flags.
		if cmd.HasSubCommands() {
			PostInitCommands(envPrefix, flags, cmd.Commands()...) // Recursively initialize subcommands.
		}
	}
}

// PresetRequiredFlags binds each flag of the given Cobra command
// to a corresponding environment variable, if such a variable is set.
// This function uses Viper to read the environment variable that matches
// the flag name and sets the flag's value accordingly.
//
// Parameters:
// - cmd: The Cobra command whose flags are to be initialized.
//
// This function iterates through all flags of the given command,
// binding them to environment variables and setting their values if applicable.
func PresetRequiredFlags(envPrefix string, flags map[*pflag.Flag]bool, cmd *cobra.Command) {
	_ = viper.BindPFlags(cmd.Flags()) // Bind the command's flags to Viper.
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if flags[f] {
			return
		}

		flags[f] = true

		if noEnvFlags[f.Name] {
			return
		}

		viperKey := f.Name
		if annotations := f.Annotations[viperKeyAnnotation]; len(annotations) > 0 {
			viperKey = annotations[0]
		}

		envVarName := strings.ToUpper(envPrefix + "_" + strings.ReplaceAll(viperKey, "-", "_"))
		newUsage := fmt.Sprintf("%s [env: %s]", f.Usage, envVarName)
		f.Usage = newUsage

		if viper.IsSet(viperKey) && viper.GetString(viperKey) != "" {
			_ = cmd.Flags().Set(viperKey, viper.GetString(viperKey)) // Set flag value from environment variable.
		}
	})
}
