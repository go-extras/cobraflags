package cobraflags

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// initOnceMap stores sync.Once instances per command to prevent multiple initializations
// of the same command while allowing different commands to be initialized independently
var initOnceMap = make(map[*cobra.Command]*sync.Once)
var initOnceMutex sync.Mutex

var noEnvFlags = map[string]bool{
	"help": true,
}

// CobraOnInitialize initializes Cobra command(s) with automatic environment variable binding.
// This function sets up Viper to automatically detect and bind environment variables
// to command flags based on the provided prefix. It should be called after registering
// flags but before executing the Cobra command.
//
// Environment Variable Mapping:
// Flags are automatically mapped to environment variables using the pattern:
// {envPrefix}_{FLAG_NAME} where FLAG_NAME is the flag name converted to uppercase
// with hyphens replaced by underscores.
//
// Examples:
//   - Flag "config-file" with prefix "MYAPP" → "MYAPP_CONFIG_FILE"
//   - Flag "port" with prefix "SERVER" → "SERVER_PORT"
//   - Flag "verbose" with prefix "CLI" → "CLI_VERBOSE"
//
// Custom Viper Keys:
// If a flag has a custom ViperKey set, the environment variable will be based
// on the ViperKey instead of the flag name, following the same transformation rules.
//
// Parameters:
//   - envPrefix: Environment variable prefix (without trailing underscore)
//   - command: Root Cobra command to initialize (subcommands are processed recursively)
//
// Usage Example:
//
//	cmd := &cobra.Command{Use: "myapp"}
//	flag := &StringFlag{Name: "config", Value: "config.yaml"}
//	flag.Register(cmd)
//	CobraOnInitialize("MYAPP", cmd)  // Binds to MYAPP_CONFIG
//	cmd.Execute()
//
// Note: This function modifies the help function to ensure initialization occurs
// before help is displayed, and uses sync.Once to prevent multiple initializations.
func CobraOnInitialize(envPrefix string, command *cobra.Command) {
	// Get or create a sync.Once for this specific command
	initOnceMutex.Lock()
	initOnce, exists := initOnceMap[command]
	if !exists {
		initOnce = &sync.Once{}
		initOnceMap[command] = initOnce
	}
	initOnceMutex.Unlock()

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

		envVarName := strings.ToUpper(envPrefix + "_" + strings.ReplaceAll(strings.ReplaceAll(viperKey, ".", "_"), "-", "_"))
		newUsage := fmt.Sprintf("%s [env: %s]", f.Usage, envVarName)
		f.Usage = newUsage

		if viper.IsSet(viperKey) && viper.GetString(viperKey) != "" {
			_ = cmd.Flags().Set(f.Name, viper.GetString(viperKey)) // Set flag value from environment variable.
		}
	})
}
