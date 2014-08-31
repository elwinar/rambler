package commands

import (
	"log"
	"github.com/spf13/cobra"
	"rambler/config"
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
	config.Init(cmd.Flags().Lookup("configuration"))
	
	err := lib.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	
	if !lib.HasMigrationTable() {
		err := lib.CreateMigrationTable()
		if err != nil {
			log.Println(err)
			return
		}
	}
	
	_, err = lib.GetMigrationsFiles()
	if err != nil {
		log.Println(err)
		return
	}
	
	log.Println("Done")
}
