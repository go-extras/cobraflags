package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*StringFlag)(nil)

// StringFlag is a flag that holds a string value.
type StringFlag FlagBase[string]

func (s *StringFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.String(s.Name, s.Value, s.Usage)
	} else {
		flags.StringP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)
}

func (s *StringFlag) GetString() string {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	return viper.GetString(s.Name)
}

func (s *StringFlag) GetStringE() (string, error) {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	v := viper.GetString(s.Name)

	if s.ValidateFunc != nil {
		err := s.ValidateFunc(v)
		if err != nil {
			return "", err
		}
	}

	return v, nil
}
