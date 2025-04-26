package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*IntFlag)(nil)

// IntFlag is a flag that holds an integer value.
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
}

func (s *IntFlag) GetInt() int {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	return viper.GetInt(s.Name)
}

func (s *IntFlag) GetIntE() (int, error) {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	v := viper.GetInt(s.Name)

	if result, err := pIntFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
