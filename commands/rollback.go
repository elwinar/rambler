package commands

import (
	"github.com/spf13/cobra"
	"log"
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
	
	for _, file := range files {
		sections, err := file.Down()
		log.Println(sections, err)
	}
	
	log.Println("Done")
}
