package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"rambler/config"
)

var (
	Rambler = &cobra.Command{
		Use:   "rambler",
		Short: "Rambler is a simple and language-independant SQL schema migration tool",
	}
)

func init() {
	viper.SetDefault("quiet", false)
	Rambler.PersistentFlags().BoolP("quiet", "q", false, "supress all output")
	config.BindPFlag("quiet", Rambler.PersistentFlags().Lookup("quiet"))

	viper.SetDefault("verbose", false)
	Rambler.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	config.BindPFlag("verbose", Rambler.PersistentFlags().Lookup("verbose"))

	Rambler.PersistentFlags().StringP("configuration", "c", "rambler.json", "configuration file")
	
	viper.SetDefault("driver", "mysql")
	Rambler.PersistentFlags().String("driver", "mysql", "database driver to use")
	config.BindPFlag("driver", Rambler.PersistentFlags().Lookup("driver"))
	
	viper.SetDefault("protocol", "tcp")
	Rambler.PersistentFlags().String("protocol", "tcp", "host to connect to")
	config.BindPFlag("protocol", Rambler.PersistentFlags().Lookup("protocol"))

	viper.SetDefault("host", "localhost")
	Rambler.PersistentFlags().StringP("host", "h", "localhost", "host to connect to")
	config.BindPFlag("host", Rambler.PersistentFlags().Lookup("host"))
	
	viper.SetDefault("port", 3306)
	Rambler.PersistentFlags().Int("port", 3306, "host to connect to")
	config.BindPFlag("host", Rambler.PersistentFlags().Lookup("host"))

	viper.SetDefault("user", "root")
	Rambler.PersistentFlags().StringP("user", "u", "root", "user to connect as")
	config.BindPFlag("user", Rambler.PersistentFlags().Lookup("user"))

	viper.SetDefault("password", "")
	Rambler.PersistentFlags().StringP("password", "p", "", "password to connect with")
	config.BindPFlag("password", Rambler.PersistentFlags().Lookup("password"))

	viper.SetDefault("database", "")
	Rambler.PersistentFlags().StringP("database", "d", "", "database to use")
	config.BindPFlag("database", Rambler.PersistentFlags().Lookup("database"))

	viper.SetDefault("migrations", ".")
	Rambler.PersistentFlags().StringP("migrations", "m", ".", "migrations directory")
	config.BindPFlag("migrations", Rambler.PersistentFlags().Lookup("migrations"))
	
	viper.SetConfigName("rambler")
	viper.AddConfigPath(".")
}
