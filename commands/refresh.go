package commands

import (
	"github.com/spf13/cobra"
	"log"
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
	err := lib.Init(cmd.Flags().Lookup("configuration"))
	if err != nil {
		log.Println(err)
		return
	}
	
	log.Println("Done")
}
