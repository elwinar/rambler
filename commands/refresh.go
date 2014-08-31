package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"rambler/config"
	"rambler/lib"
)

var (
	Refresh = &cobra.Command{
		Use:   "refresh",
		Short: "Refresh a database",
		Run:   refresh,
	}
)

func init() {
	Rambler.AddCommand(Refresh)
}

func refresh (cmd *cobra.Command, args []string) {
	config.Init(cmd.Flags().Lookup("configuration"))
	
	err := lib.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Done")
}
