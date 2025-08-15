package cobraflags_test

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/spf13/cobra"

	"github.com/go-extras/cobraflags"
)

// TestCobraOnInitialize_ComplexHierarchy tests environment variable binding
// with complex command hierarchies including subcommands and nested subcommands.
func TestCobraOnInitialize_ComplexHierarchy(t *testing.T) {
	c := qt.New(t)

	// Set up environment variables
	envVars := map[string]string{
		"TESTAPP_GLOBAL_CONFIG": "global.yaml",
		"TESTAPP_SERVER_PORT":   "9090",
		"TESTAPP_DB_HOST":       "localhost",
		"TESTAPP_DB_PORT":       "5432",
		"TESTAPP_VERBOSE":       "true",
	}

	// Set environment variables and defer cleanup
	for key, value := range envVars {
		os.Setenv(key, value)
	}
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()

	// Create flags
	globalConfigFlag := &cobraflags.StringFlag{
		Name:  "global-config",
		Usage: "Global configuration file",
		Value: "default.yaml",
	}

	verboseFlag := &cobraflags.BoolFlag{
		Name:       "verbose",
		Usage:      "Enable verbose output",
		Value:      false,
		Persistent: true, // Available to all subcommands
	}

	serverPortFlag := &cobraflags.IntFlag{
		Name:  "server-port",
		Usage: "Server port",
		Value: 8080,
	}

	dbHostFlag := &cobraflags.StringFlag{
		Name:  "db-host",
		Usage: "Database host",
		Value: "127.0.0.1",
	}

	dbPortFlag := &cobraflags.IntFlag{
		Name:  "db-port",
		Usage: "Database port",
		Value: 3306,
	}

	// Create command hierarchy
	rootCmd := &cobra.Command{
		Use:   "testapp",
		Short: "Test application",
		Run: func(_ *cobra.Command, _ []string) {
			// Root command execution
		},
	}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Server commands",
		Run: func(_ *cobra.Command, _ []string) {
			// Server command execution
		},
	}

	dbCmd := &cobra.Command{
		Use:   "db",
		Short: "Database commands",
		Run: func(_ *cobra.Command, _ []string) {
			// DB command execution
		},
	}

	// Register flags with appropriate commands
	globalConfigFlag.Register(rootCmd)
	verboseFlag.Register(rootCmd) // Persistent flag
	serverPortFlag.Register(serverCmd)
	dbHostFlag.Register(dbCmd)
	dbPortFlag.Register(dbCmd)

	// Add subcommands
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(dbCmd)

	// Initialize with environment variable binding
	cobraflags.CobraOnInitialize("TESTAPP", rootCmd)

	// Test root command with global flags
	rootCmd.SetArgs(make([]string, 0))
	err := rootCmd.Execute()
	c.Assert(err, qt.IsNil)

	// Verify environment variable binding
	c.Assert(globalConfigFlag.GetString(), qt.Equals, "global.yaml")
	c.Assert(verboseFlag.GetBool(), qt.Equals, true)

	// Test server subcommand
	rootCmd.SetArgs([]string{"server"})
	err = rootCmd.Execute()
	c.Assert(err, qt.IsNil)

	c.Assert(serverPortFlag.GetInt(), qt.Equals, 9090)
	c.Assert(verboseFlag.GetBool(), qt.Equals, true) // Persistent flag

	// Test nested subcommand (db under server)
	rootCmd.SetArgs([]string{"server", "db"})
	err = rootCmd.Execute()
	c.Assert(err, qt.IsNil)

	c.Assert(dbHostFlag.GetString(), qt.Equals, "localhost")
	c.Assert(dbPortFlag.GetInt(), qt.Equals, 5432)
	c.Assert(verboseFlag.GetBool(), qt.Equals, true) // Persistent flag available
}

// TestCobraOnInitialize_SubcommandInheritance tests that persistent flags
// are properly inherited by subcommands and environment variables work correctly.
func TestCobraOnInitialize_SubcommandInheritance(t *testing.T) {
	c := qt.New(t)

	// Set environment variables
	c.Setenv("INHERIT_VERBOSE", "true")
	c.Setenv("INHERIT_CONFIG", "test.yaml")

	// Create persistent and non-persistent flags
	verboseFlag := &cobraflags.BoolFlag{
		Name:       "verbose",
		Usage:      "Enable verbose output",
		Value:      false,
		Persistent: true,
	}

	configFlag := &cobraflags.StringFlag{
		Name:       "config",
		Usage:      "Configuration file",
		Value:      "default.yaml",
		Persistent: false, // Not persistent
	}

	subConfigFlag := &cobraflags.StringFlag{
		Name:  "sub-config",
		Usage: "Subcommand configuration",
		Value: "sub.yaml",
	}

	// Create command hierarchy
	rootCmd := &cobra.Command{
		Use: "inherit",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	subCmd := &cobra.Command{
		Use: "sub",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	// Register flags
	verboseFlag.Register(rootCmd)
	configFlag.Register(rootCmd)
	subConfigFlag.Register(subCmd)

	rootCmd.AddCommand(subCmd)

	// Initialize
	cobraflags.CobraOnInitialize("INHERIT", rootCmd)

	// Test that persistent flag is available in subcommand
	rootCmd.SetArgs([]string{"sub"})
	err := rootCmd.Execute()
	c.Assert(err, qt.IsNil)

	// Persistent flag should be available and bound to env var
	c.Assert(verboseFlag.GetBool(), qt.Equals, true)

	// Non-persistent flag should still work from root
	rootCmd.SetArgs(make([]string, 0))
	err = rootCmd.Execute()
	c.Assert(err, qt.IsNil)
	c.Assert(configFlag.GetString(), qt.Equals, "test.yaml")
}

// TestConcurrentFlagAccess tests that flag operations are thread-safe
// when accessed concurrently from multiple goroutines.
func TestConcurrentFlagAccess(t *testing.T) {
	c := qt.New(t)

	// Create a flag for concurrent testing
	testFlag := &cobraflags.StringFlag{
		Name:  "concurrent-test",
		Usage: "Flag for concurrent access testing",
		Value: "default",
	}

	cmd := &cobra.Command{
		Use: "concurrent",
		Run: func(_ *cobra.Command, _ []string) {},
	}

	testFlag.Register(cmd)
	cobraflags.CobraOnInitialize("CONCURRENT", cmd)

	// Set a value for the flag
	cmd.SetArgs([]string{"--concurrent-test", "test-value"})
	err := cmd.Execute()
	c.Assert(err, qt.IsNil)

	// Test concurrent reads
	const numGoroutines = 100
	const numReadsPerGoroutine = 10

	results := make(chan string, numGoroutines*numReadsPerGoroutine)
	errors := make(chan error, numGoroutines*numReadsPerGoroutine)

	// Launch multiple goroutines that read the flag value
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numReadsPerGoroutine; j++ {
				// Test both Get and GetE methods
				value := testFlag.GetString()
				results <- value

				valueE, err := testFlag.GetStringE()
				if err != nil {
					errors <- err
				} else {
					results <- valueE
				}
			}
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines*numReadsPerGoroutine*2; i++ {
		select {
		case value := <-results:
			c.Assert(value, qt.Equals, "test-value")
		case err := <-errors:
			c.Assert(err, qt.IsNil, qt.Commentf("Unexpected error in concurrent access: %v", err))
		}
	}
}
