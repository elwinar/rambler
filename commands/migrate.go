package commands

import (
	"log"
	"github.com/spf13/cobra"
	"rambler/lib"
)

var (
	Migrate = &cobra.Command{
		Use:   "migrate",
		Short: "Migrate a database",
		Run:   migrate,
	}
)

func init() {
	Rambler.AddCommand(Migrate)
}

func migrate (cmd *cobra.Command, args []string) {
	err := lib.Init(cmd.Flags().Lookup("configuration"))
	if err != nil {
		log.Println(err)
		return
	}
	
	files, err := lib.GetMigrationsFiles()
	if err != nil {
		log.Println(err)
		return
	}
	
	log.Println(files)
	
	log.Println("Done")
}
