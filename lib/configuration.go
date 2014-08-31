package lib

import (
	"errors"
	"github.com/elwinar/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	pflags map[string]*pflag.Flag
)

func init() {
	pflags = make(map[string]*pflag.Flag)
}

func Init(configuration *pflag.Flag) error {
	if configuration != nil {
		viper.SetConfigFile(configuration.Value.String())
	}
	
	viper.ReadInConfig()
	
	for key, flag := range pflags {
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
	
	SetQuiet(viper.GetBool("quiet"))
	
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

func BindPFlag(key string, flag *pflag.Flag) {
	pflags[key] = flag
}
