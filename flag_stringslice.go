package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*StringSliceFlag)(nil)

// StringSliceFlag is a flag that holds a string value.
type StringSliceFlag FlagBase[[]string]

// pStringSliceFlag is an alias for a pointer to FlagBase[[]string].
type pStringSliceFlag = *FlagBase[[]string]

func (s *StringSliceFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.StringSlice(s.Name, s.Value, s.Usage)
	} else {
		flags.StringSliceP(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)
}

func (s *StringSliceFlag) GetStringSlice() []string {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	return viper.GetStringSlice(s.Name)
}

func (s *StringSliceFlag) GetStringSliceE() ([]string, error) {
	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(s.Name, s.flag))
	})

	v := viper.GetStringSlice(s.Name)

	if result, err := pStringSliceFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
