package commands

import (
	"github.com/elwinar/cobra"
	"github.com/elwinar/viper"
)

func init() {
	// Setup the configuration file lookup
	viper.SetConfigName("rambler")
	viper.AddConfigPath("/etc/rambler")
	viper.AddConfigPath("$HOME/.rambler")
	viper.AddConfigPath(".")
	
	// Add the configuration flag to choose the configuration on the command line
	Rambler.PersistentFlags().StringP("configuration", "c", "", "read the configuration from the given file")
	
	// Set the default configuration
	viper.SetDefault("quiet", false)
	viper.SetDefault("verbose", false)
	viper.SetDefault("driver", "mysql")
	viper.SetDefault("protocol", "tcp")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 3306)
	viper.SetDefault("user", "root")
	viper.SetDefault("password", "")
	viper.SetDefault("database", "")
	viper.SetDefault("migrations", ".")
	
	// Add ubiquitous flags to the main command
	Rambler.PersistentFlags().BoolP("quiet", "q", false, "supress all output")
	Rambler.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	Rambler.PersistentFlags().String("driver", "mysql", "database driver to use")
	Rambler.PersistentFlags().String("protocol", "tcp", "host to connect to")
	Rambler.PersistentFlags().StringP("host", "h", "localhost", "host to connect to")
	Rambler.PersistentFlags().Int("port", 3306, "host to connect to")
	Rambler.PersistentFlags().StringP("user", "u", "root", "user to connect as")
	Rambler.PersistentFlags().StringP("password", "p", "", "password to connect with")
	Rambler.PersistentFlags().StringP("database", "d", "", "database to use")
	Rambler.PersistentFlags().StringP("migrations", "m", ".", "migrations directory")
	
	// Set overrides from the command-line to viper
	
	viper.BindPFlag("quiet", Rambler.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("verbose", Rambler.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("driver", Rambler.PersistentFlags().Lookup("driver"))
	viper.BindPFlag("protocol", Rambler.PersistentFlags().Lookup("protocol"))
	viper.BindPFlag("host", Rambler.PersistentFlags().Lookup("host"))
	viper.BindPFlag("host", Rambler.PersistentFlags().Lookup("host"))
	viper.BindPFlag("user", Rambler.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", Rambler.PersistentFlags().Lookup("password"))
	viper.BindPFlag("database", Rambler.PersistentFlags().Lookup("database"))
	viper.BindPFlag("migrations", Rambler.PersistentFlags().Lookup("migrations"))
}

var Rambler = &cobra.Command{
	Use:   "rambler",
	Short: "Rambler is a simple and language-independant SQL schema migration tool",
	Long:  "Rambler is a simple and language-independant SQL schema migration tool",
	Run:   nil,
}
