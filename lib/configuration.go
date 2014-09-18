package lib

import (
	"errors"
	"github.com/elwinar/cast"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	flags map[string]*pflag.Flag
	Flags  pflag.FlagSet
)

func init() {
	flags = make(map[string]*pflag.Flag)
	
	viper.SetConfigName("rambler")
	viper.AddConfigPath(".")
	
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
	
	Flags.StringP("configuration", "c", "rambler.json", "configuration file")
	Flags.BoolP("quiet", "q", false, "supress all output")
	Flags.BoolP("verbose", "v", false, "verbose output")
	Flags.String("driver", "mysql", "database driver to use")
	Flags.String("protocol", "tcp", "host to connect to")
	Flags.StringP("host", "h", "localhost", "host to connect to")
	Flags.Int("port", 3306, "host to connect to")
	Flags.StringP("user", "u", "root", "user to connect as")
	Flags.StringP("password", "p", "", "password to connect with")
	Flags.StringP("database", "d", "", "database to use")
	Flags.StringP("migrations", "m", ".", "migrations directory")
	
	flags["quiet"] = Flags.Lookup("quiet")
	flags["verbose"] = Flags.Lookup("verbose")
	flags["driver"] = Flags.Lookup("driver")
	flags["protocol"] = Flags.Lookup("protocol")
	flags["host"] = Flags.Lookup("host")
	flags["host"] = Flags.Lookup("host")
	flags["user"] = Flags.Lookup("user")
	flags["password"] = Flags.Lookup("password")
	flags["database"] = Flags.Lookup("database")
	flags["migrations"] = Flags.Lookup("migrations")
	
}

func Init(configuration *pflag.Flag) error {
	if configuration != nil {
		viper.SetConfigFile(configuration.Value.String())
	}
	
	viper.ReadInConfig()
	
	for key, flag := range flags {
		if flag != nil && flag.Changed {
			switch flag.Value.Type() {
				case "int", "int8", "int16", "int32", "int64":
					viper.Set(key, cast.ToInt(flag.Value.String()))
				case "bool":
					viper.Set(key, cast.ToBool(flag.Value.String()))
				default:
					viper.Set(key, flag.Value.String())
			}
		}
	}
	
	jww.SetStdoutThreshold(jww.LevelTrace)
	jww.SetLogThreshold(jww.LevelTrace)
	if viper.GetBool("quiet") {
		jww.SetStdoutThreshold(jww.LevelError)
		jww.DiscardLogging()
	} else if viper.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelInfo)
		jww.SetLogThreshold(jww.LevelInfo)
	}
	
	if !viper.IsSet("database") || viper.GetString("database") == "" {
		return errors.New("No database selected")
	}
	
	err := Connect()
	if err != nil {
		return err
	}
	
	if !HasMigrationTable() {
		err := CreateMigrationTable()
		if err != nil {
			return err
		}
	}
	
	return nil
}
