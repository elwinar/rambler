package commands

import (
	"github.com/elwinar/cobra"
	"github.com/elwinar/viper"
	jww "github.com/spf13/jwalterweatherman"
	"path/filepath"
)

func init() {
	Rambler.AddCommand(Debug)
}

var Debug = &cobra.Command{
	Use:   "debug",
	Short: "Print cli debug informations",
	Run:   do(func (cmd *cobra.Command, args []string) {
		jww.INFO.Println("driver:", viper.GetString("driver"))
		jww.INFO.Println("protocol:", viper.GetString("protocol"))
		jww.INFO.Println("host:", viper.GetString("host"))
		jww.INFO.Println("port:", viper.GetString("port"))
		jww.INFO.Println("user:", viper.GetString("user"))
		jww.INFO.Println("password:", viper.GetString("password"))
		jww.INFO.Println("database:", viper.GetString("database"))
		path, _ := filepath.Abs(viper.GetString("migrations"))
		jww.INFO.Println("migrations:", path)
	}),
}
