package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"rambler/config"
	"rambler/lib"
)

var (
	Rollback = &cobra.Command{
		Use:   "rollback",
		Short: "Rollback a database",
		Run:   rollback,
	}
)

func init() {
	Rambler.AddCommand(Rollback)
}

func rollback (cmd *cobra.Command, args []string) {
	config.Init(cmd.Flags().Lookup("configuration"))
	
	err := lib.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Done")
}
