package commands

import (
	"github.com/elwinar/cobra"
	"github.com/elwinar/rambler/lib"
	"github.com/elwinar/viper"
)

func init() {
	// Setup the configuration file lookup
	viper.SetConfigName("rambler")
	viper.AddConfigPath("/etc/rambler")
	viper.AddConfigPath("$HOME/.rambler")
	viper.AddConfigPath(".")

	// Add configuration flags to the command-line
	Rambler.PersistentFlags().StringP("configuration", "c", "", "read the configuration from the given file")
	
	Rambler.PersistentFlags().StringP("environment", "e", "default", "environment to run on")
	viper.BindPFlag("environment", Rambler.PersistentFlags().Lookup("environment"))
	
	Rambler.PersistentFlags().BoolP("quiet", "q", false, "supress all output")
	viper.BindPFlag("quiet", Rambler.PersistentFlags().Lookup("quiet"))
	
	Rambler.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", Rambler.PersistentFlags().Lookup("verbose"))

	// Set the default configuration
	viper.SetDefault("environments", map[string]lib.RawEnvironment{})
	
	// Set the default environment
	viper.SetDefault("database", "")
	viper.SetDefault("driver", "mysql")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("migrations", ".")
	viper.SetDefault("password", "")
	viper.SetDefault("port", 3306)
	viper.SetDefault("protocol", "tcp")
	viper.SetDefault("user", "root")
	
	// Set the environment overrides
	Rambler.PersistentFlags().StringP("database", "d", "", "database to use")
	Rambler.PersistentFlags().String("driver", "mysql", "database driver to use")
	Rambler.PersistentFlags().StringP("host", "h", "localhost", "host to connect to")
	Rambler.PersistentFlags().StringP("migrations", "m", ".", "migrations directory")
	Rambler.PersistentFlags().StringP("password", "p", "", "password to connect with")
	Rambler.PersistentFlags().Int("port", 3306, "host to connect to")
	Rambler.PersistentFlags().String("protocol", "tcp", "host to connect to")
	Rambler.PersistentFlags().StringP("user", "u", "root", "user to connect as")
}

var Rambler = &cobra.Command{
	Use:   "rambler",
	Short: "Rambler is a simple and language-independant SQL schema migration tool",
	Long:  "Rambler is a simple and language-independant SQL schema migration tool",
	Run:   nil,
}
