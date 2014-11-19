package lib

import (
	"github.com/spf13/pflag"
	"strings"
	"strconv"
)

type Environment struct {
	Driver string
	Protocol string
	Host string
	Port uint64
	User string
	Password string
	Database string
	Migrations string
}

func GetEnvironment(name string, configuration Configuration, flags *pflag.FlagSet) (*Environment, error) {
	return getEnvironment(name, configuration, flags)
}

func getEnvironment(name string, configuration Configuration, flags *pflag.FlagSet) (*Environment, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrUnknownEnvironment
	}
	
	override, found := configuration.Environments[name]
	if !found {
		return nil, ErrUnknownEnvironment
	}
	
	var environment Environment
	
	environment.Driver = configuration.Driver
	if override.Driver != nil {
		environment.Driver = *override.Driver
	}
	if flags.Lookup("driver").Changed {
		environment.Driver = flags.Lookup("driver").Value.String()
	}
	
	environment.Protocol = configuration.Protocol
	if override.Protocol != nil {
		environment.Protocol = *override.Protocol
	}
	if flags.Lookup("protocol").Changed {
		environment.Protocol = flags.Lookup("protocol").Value.String()
	}
	
	environment.Host = configuration.Host
	if override.Host != nil {
		environment.Host = *override.Host
	}
	if flags.Lookup("host").Changed {
		environment.Host = flags.Lookup("host").Value.String()
	}
	
	environment.Port = configuration.Port
	if override.Port != nil {
		environment.Port = *override.Port
	}
	if flags.Lookup("port").Changed {
		environment.Port, _ = strconv.ParseUint(flags.Lookup("port").Value.String(), 10, 64)
	}
	
	environment.User = configuration.User
	if override.User != nil {
		environment.User = *override.User
	}
	if flags.Lookup("user").Changed {
		environment.User = flags.Lookup("user").Value.String()
	}
	
	environment.Password = configuration.Password
	if override.Password != nil {
		environment.Password = *override.Password
	}
	if flags.Lookup("password").Changed {
		environment.Password = flags.Lookup("password").Value.String()
	}
	
	environment.Database = configuration.Database
	if override.Database != nil {
		environment.Database = *override.Database
	}
	if flags.Lookup("database").Changed {
		environment.Database = flags.Lookup("database").Value.String()
	}
	
	environment.Migrations = configuration.Migrations
	if override.Migrations != nil {
		environment.Migrations = *override.Migrations
	}
	if flags.Lookup("migrations").Changed {
		environment.Migrations = flags.Lookup("migrations").Value.String()
	}
	
	return &environment, nil
}
