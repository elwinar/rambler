package commands

import (
	"github.com/elwinar/cobra"
	"github.com/elwinar/viper"
	jww "github.com/spf13/jwalterweatherman"
)

type command func(*cobra.Command, []string)

func do(f command) command {
	return func(cmd *cobra.Command, args []string) {
		var configuration = cmd.Flags().Lookup("configuration")
		if configuration != nil {
			viper.SetConfigFile(configuration.Value.String())
		}

		viper.ReadInConfig()

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
