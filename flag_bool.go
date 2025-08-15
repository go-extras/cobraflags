package cobraflags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*BoolFlag)(nil)

// BoolFlag is a flag that holds a boolean value.
type BoolFlag FlagBase[bool]

// pBoolFlag is an alias for a pointer to FlagBase[bool].
type pBoolFlag = *FlagBase[bool]

func (s *BoolFlag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}

	flags.BoolP(s.Name, s.Shorthand, s.Value, s.Usage)

	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pBoolFlag(s).getViperKey()}
}

func (s *BoolFlag) GetBool() bool {
	viperKey := pBoolFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return viper.GetBool(viperKey)
}

func (s *BoolFlag) GetBoolE() (bool, error) {
	viperKey := pBoolFlag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	v := viper.GetBool(viperKey)

	if result, err := pBoolFlag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
