package config

import (
	"github.com/elwinar/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"rambler/lib"
)

var (
	pflags map[string]*pflag.Flag
)

func init() {
	pflags = make(map[string]*pflag.Flag)
}

func Init(configuration *pflag.Flag) {
	if configuration != nil {
		viper.SetConfigFile(configuration.Value.String())
	}
	
	viper.ReadInConfig()
	
	for key, flag := range pflags {
		if flag != nil && flag.Changed {
			switch flag.Value.Type() {
				case "int", "int8", "int16", "int32", "int64":
					viper.Set(key, cast.ToInt(flag.Value.String()))
				case "bool":
					viper.Set(key, cast.ToBool(flag.Value.String()))
				default:
					viper.Set(key, flag.Value.String())
			}
		}
	}
	
	lib.SetQuiet(viper.GetBool("quiet"))
}

func BindPFlag(key string, flag *pflag.Flag) {
	pflags[key] = flag
}
