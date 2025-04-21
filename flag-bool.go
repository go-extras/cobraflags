package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*BoolFlag)(nil)

// BoolFlag is a flag that holds a boolean value.
type BoolFlag FlagBase[bool]

func (s *BoolFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.Bool(s.Name, s.Value, s.Usage)
	} else {
		flags.BoolP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)
}

func (s *BoolFlag) GetBool() bool {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	return viper.GetBool(s.Name)
}

func (s *BoolFlag) GetBoolE() (bool, error) {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	v := viper.GetBool(s.Name)

	if s.ValidateFunc != nil {
		err := s.ValidateFunc(v)
		if err != nil {
			return false, err
		}
	}

	return v, nil
}
