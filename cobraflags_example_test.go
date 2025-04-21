package cobraflags_test

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/go-extras/cobraflags"
)

// ExampleIntFlag demonstrates how to use an IntFlag with Cobra commands
func ExampleIntFlag() {
	myFlag := &cobraflags.IntFlag{
		Name:     "count",
		Usage:    "Number of items to process",
		Value:    10,
		Required: true,
	}

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process items",
		Run: func(_cmd *cobra.Command, _args []string) {
			value := myFlag.GetInt()
			fmt.Printf("Processing %d items\n", value)
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--count", "10"})
	_ = cmd.Execute()

	// Output:
	// Processing 10 items
}

// ExampleStringFlag demonstrates how to use a StringFlag with Cobra commands
func ExampleStringFlag() {
	myFlag := &cobraflags.StringFlag{
		Name:     "message",
		Usage:    "Message to display",
		Value:    "Hello, World!",
		Required: true,
	}

	cmd := &cobra.Command{
		Use:   "greet",
		Short: "Greeting application",
		Run: func(_cmd *cobra.Command, _args []string) {
			message := myFlag.GetString()
			fmt.Printf("Message: %s\n", message)
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--message", "Hello, Cobra!"})
	_ = cmd.Execute()

	// Output:
	// Message: Hello, Cobra!
}

// ExampleBoolFlag demonstrates how to use a BoolFlag with Cobra commands
func ExampleBoolFlag() {
	myFlag := &cobraflags.BoolFlag{
		Name:     "verbose",
		Usage:    "Enable verbose output",
		Value:    false,
		Required: false,
	}

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run the application",
		Run: func(_cmd *cobra.Command, _args []string) {
			if myFlag.GetBool() {
				fmt.Println("Verbose mode enabled")
			} else {
				fmt.Println("Verbose mode disabled")
			}
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--verbose"})
	_ = cmd.Execute()

	// Output:
	// Verbose mode enabled
}

// ExampleStringSliceFlag demonstrates how to use a StringSliceFlag with Cobra commands
func ExampleStringSliceFlag() {
	myFlag := &cobraflags.StringSliceFlag{
		Name:     "items",
		Usage:    "List of items to process",
		Value:    []string{"item1", "item2"},
		Required: false,
	}

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process items",
		Run: func(_cmd *cobra.Command, _args []string) {
			items := myFlag.GetStringSlice()
			fmt.Printf("Processing items: %v\n", items)
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--items", "item3,item4"})
	_ = cmd.Execute()

	// Output:
	// Processing items: [item3 item4]
}

// ExampleUint8Flag demonstrates how to use a Uint8Flag with Cobra commands
func ExampleUint8Flag() {
	myFlag := &cobraflags.Uint8Flag{
		Name:     "level",
		Usage:    "Level of detail (0-255)",
		Value:    128,
		Required: false,
	}

	cmd := &cobra.Command{
		Use:   "setlevel",
		Short: "Set the level",
		Run: func(_cmd *cobra.Command, _args []string) {
			level := myFlag.GetUint8()
			fmt.Printf("Setting level to %d\n", level)
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--level", "200"})
	_ = cmd.Execute()

	// Output:
	// Setting level to 200
}

// ExampleIntFlag_withValidation demonstrates adding validation to an IntFlag
func ExampleIntFlag_withValidation() {
	myFlag := &cobraflags.IntFlag{
		Name:     "count",
		Usage:    "Number of items to process (must be positive)",
		Value:    5,
		Required: true,
		ValidateFunc: func(value int) error {
			if value <= 0 {
				return fmt.Errorf("count must be positive")
			}
			return nil
		},
	}

	cmd := &cobra.Command{
		Use:   "process",
		Short: "Process items",
		Run: func(_cmd *cobra.Command, _args []string) {
			value, err := myFlag.GetIntE()
			if err != nil {
				fmt.Printf("Validation error: %v\n", err)
				return
			}
			fmt.Printf("Processing %d items\n", value)
		},
	}

	myFlag.Register(cmd)
	cmd.SetArgs([]string{"--count", "-5"})
	_ = cmd.Execute()

	// Output:
	// Validation error: count must be positive
}

// ExampleRegister demonstrates how to register multiple flags at once
func ExampleRegister() {
	countFlag := &cobraflags.IntFlag{
		Name:  "count",
		Usage: "Number of items",
		Value: 3,
	}

	verboseFlag := &cobraflags.BoolFlag{
		Name:  "verbose",
		Usage: "Enable verbose output",
		Value: true,
	}

	cmd := &cobra.Command{
		Use:   "app",
		Short: "Example application",
		Run: func(_cmd *cobra.Command, _args []string) {
			count := countFlag.GetInt()
			verbose := verboseFlag.GetBool()

			fmt.Printf("Count: %d, Verbose: %v\n", count, verbose)
		},
	}

	cobraflags.Register(cmd, countFlag, verboseFlag)
	cmd.SetArgs(make([]string, 0))
	_ = cmd.Execute()

	// Output:
	// Count: 3, Verbose: true
}

// ExampleCobraOnInitialize demonstrates environment variable binding
func ExampleCobraOnInitialize() {
	// Set environment variable for demo
	os.Setenv("MYAPP_MESSAGE", "from environment")
	defer os.Unsetenv("MYAPP_MESSAGE")

	messageFlag := &cobraflags.StringFlag{
		Name:  "message",
		Usage: "Message to display",
		Value: "default message",
	}

	cmd := &cobra.Command{
		Use:   "greet",
		Short: "Greeting application",
		Run: func(_cmd *cobra.Command, _args []string) {
			message := messageFlag.GetString()
			fmt.Printf("Message: %s\n", message)
		},
	}

	messageFlag.Register(cmd)
	cobraflags.CobraOnInitialize("MYAPP", cmd)
	cmd.SetArgs(make([]string, 0))
	_ = cmd.Execute()

	// Output:
	// Message: from environment
}
