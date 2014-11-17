package commands

import (
	"github.com/elwinar/cobra"
	jww "github.com/elwinar/jwalterweatherman"
	"github.com/elwinar/rambler/lib"
	"github.com/elwinar/viper"
)

type command func(*cobra.Command, []string)

// do is the preparation function that gather all common instruction for commands
// It should not contain any logic at all, only environment configuration
func do(f command) command {
	return func(cmd *cobra.Command, args []string) {
		var err error
		
		var configuration = cmd.Flags().Lookup("configuration")
		if configuration != nil {
			viper.SetConfigFile(configuration.Value.String())
		}

		err = viper.ReadInConfig()
		if err != nil {
			jww.ERROR.Println("Unable to load configuration file:", err)
			return
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
		
		// Load the working environment configuration
		jww.TRACE.Println("Loading configuration")
		err = lib.LoadEnvironment(cmd)
		if err != nil {
			jww.ERROR.Println("Error while loading configuration:", err)
			return
		}

		f(cmd, args)
	}
}
