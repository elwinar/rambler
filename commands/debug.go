package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"rambler/lib"
)

var (
	Debug = &cobra.Command{
		Use:   "debug",
		Short: "Print debug informations",
		Run:   debug,
	}
)

func init() {
	Rambler.AddCommand(Debug)
}

func debug (cmd *cobra.Command, args []string) {
	lib.Init(cmd.Flags().Lookup("configuration"))
	
	cmd.DebugFlags()
	viper.Debug()
}
