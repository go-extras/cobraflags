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

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pStringSliceFlag(s).getViperKey()}
}

func (s *StringSliceFlag) GetStringSlice() []string {
	viperKey := pStringSliceFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetStringSlice(viperKey)
}

func (s *StringSliceFlag) GetStringSliceE() ([]string, error) {
	viperKey := pStringSliceFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetStringSlice(viperKey)

	if result, err := pStringSliceFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
