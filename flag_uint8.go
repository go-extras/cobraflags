package cobraflags

import (
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var _ Flag = (*Uint8Flag)(nil)

// Uint8Flag is a flag that holds an uint8 value.
type Uint8Flag FlagBase[uint8]

// pUint8Flag is an alias for a pointer to FlagBase[uint8].
type pUint8Flag = *FlagBase[uint8]

func (s *Uint8Flag) Register(cmd *cobra.Command) {
	var flags *pflag.FlagSet
	if s.Persistent {
		flags = cmd.PersistentFlags()
	} else {
		flags = cmd.Flags()
	}
	if s.Shorthand == "" {
		flags.Uint8(s.Name, s.Value, s.Usage)
	} else {
		flags.Uint8P(s.Name, s.Shorthand, s.Value, s.Usage)
	}
	if s.Required {
		noError(cmd.MarkFlagRequired(s.Name))
	}
	s.flag = flags.Lookup(s.Name)

	if s.flag.Annotations == nil {
		s.flag.Annotations = make(map[string][]string)
	}
	s.flag.Annotations[viperKeyAnnotation] = []string{pUint8Flag(s).getViperKey()}
}

func (s *Uint8Flag) GetUint8() uint8 {
	viperKey := pUint8Flag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	return cast.ToUint8(viper.GetUint16(viperKey))
}

func (s *Uint8Flag) GetUint8E() (uint8, error) {
	viperKey := pUint8Flag(s).getViperKey()

	s.bindOnce.Do(func() {
		noError(viper.BindPFlag(viperKey, s.flag))
	})

	u16 := viper.GetUint16(viperKey)
	v := cast.ToUint8(u16)

	if result, err := pUint8Flag(s).validate(v); err != nil {
		return result, err
	}

	return v, nil
}
