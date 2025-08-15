# CobraFlags

[![Build Status](https://github.com/go-extras/cobraflags/actions/workflows/go-test.yml/badge.svg)](https://github.com/go-extras/cobraflags/actions/workflows/go-test.yml)

`cobraflags` is a Go module that provides an integration layer between [Cobra](https://github.com/spf13/cobra)
CLI applications and [Viper](https://github.com/spf13/viper) configuration management. It automates the binding
of environment variables to Cobra command flags, simplifying the process of managing configurations in CLI applications.

## Features

- Automatic binding of environment variables to Cobra flags.
- Support for persistent and non-persistent flags.
- Validation of flag values using custom validation functions.
- Easy integration with Cobra commands and subcommands.

## Installation

To install the package, use:

```bash
go get github.com/go-extras/cobraflags
```

## Usage

### Registering Flags

You can define and register flags using the provided `IntFlag` and `StringFlag` types. For example:

```go
import (
	"github.com/spf13/cobra"
	"github.com/go-extras/cobraflags"
)

var myFlag = &cobraflags.IntFlag{
	Name:     "example-flag",
	Usage:    "An example integer flag",
	Value:    42,
	Required: true,
}

func main() {
	cmd := &cobra.Command{
		Use:   "myapp",
		Short: "An example application",
		Run: func(cmd *cobra.Command, args []string) {
			value := myFlag.GetInt()
			fmt.Println("Flag value:", value)
		},
	}

	myFlag.Register(cmd)
	cobraflags.CobraOnInitialize("MYAPP", cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
```

### Environment Variable Binding

Flags are automatically bound to environment variables using the provided prefix. For example,
a flag named `example-flag` with the prefix `MYAPP` will be bound to the environment variable `MYAPP_EXAMPLE_FLAG`.

### Custom Viper Keys

By default, flags use their name as the Viper configuration key. You can customize this by setting the `ViperKey` field:

```go
myFlag := &cobraflags.StringFlag{
    Name:     "config-file",
    ViperKey: "app.config.file", // Custom Viper key
    Usage:    "Path to configuration file",
    Value:    "config.yaml",
}
```

This allows you to:
- Use different configuration keys than flag names
- Support nested configuration structures
- Maintain backward compatibility when renaming flags

If `ViperKey` is empty, the flag will fall back to using its `Name` for Viper binding.

### Validation

You can add custom validation logic for flags using the `ValidateFunc` field:

```go
myFlag.ValidateFunc = func(value int) error {
	if value < 0 {
		return fmt.Errorf("value must be non-negative")
	}
	return nil
}
```

You can also use the `Validator` field to provide a custom validator that implements the `cobraflags.Validator`
interface:

```go
myFlag.Validator = cobraflags.ValidatorFunc[int](func(value int) error {
	if value < 0 {
		return fmt.Errorf("value must be non-negative")
	}
	return nil
})
```

_Note: cobraflags.ValidatorFunc is used for demonstration purposes only, use your own validators_.

## Documentation

For detailed documentation, refer to the source code and comments in the package.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
