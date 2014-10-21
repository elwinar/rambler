package commands

import (
	"github.com/elwinar/cast"
	"github.com/elwinar/cobra"
	"github.com/elwinar/viper"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
)

type command func(*cobra.Command, []string)

var flags = make(map[string]*pflag.Flag)

func override(key string, flag *pflag.Flag) {
	if flag != nil {
		flags[key] = flag
	}
}

func do(f command) command {
	return func(cmd *cobra.Command, args []string) {
		var configuration = cmd.Flags().Lookup("configuration")
		if configuration != nil {
			viper.SetConfigFile(configuration.Value.String())
		}

		viper.ReadInConfig()

		for key, flag := range flags {
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

		jww.SetStdoutThreshold(jww.LevelInfo)
		jww.SetLogThreshold(jww.LevelInfo)

		if viper.GetBool("quiet") {
			jww.SetStdoutThreshold(jww.LevelError)
			jww.SetLogThreshold(jww.LevelError)
		}

		if viper.GetBool("verbose") {
			jww.SetStdoutThreshold(jww.LevelTrace)
			jww.SetLogThreshold(jww.LevelTrace)
		}

		f(cmd, args)
	}
}
