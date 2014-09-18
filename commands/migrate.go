package commands

import (
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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
		jww.ERROR.Println(err)
		return
	}
	
	files, err := lib.GetMigrationsFiles()
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	
// 	rows, err := lib.GetMigrationsRows()
// 	if err != nil {
// 		jww.ERROR.Println(err)
// 		return
// 	}
	
	for _, file := range files {
		jww.TRACE.Println(file.File)
		sections, err := file.Up()
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
		err = lib.Execute(sections)
		if err != nil {
			jww.ERROR.Println(err)
			return
		}
	}
	
	jww.INFO.Println("Done")
}
