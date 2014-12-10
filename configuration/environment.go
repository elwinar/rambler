package configuration

import (
	"errors"
	"github.com/spf13/pflag"
	"strconv"
	"strings"
)

// The various errors returned by the package
var (
	ErrUnknownEnvironment = errors.New("unknwon environment")
)

// Environment is the execution environment of a command. It contains every information
// about the database and migrations to use.
type Environment struct {
	Driver    string
	Protocol  string
	Host      string
	Port      uint64
	User      string
	Password  string
	Database  string
	Directory string
}

// GetEnvironment return the requested environment from the configuration, with
// overrides from the given flagset.
func GetEnvironment(name string, configuration Configuration, flags *pflag.FlagSet) (*Environment, error) {
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

	environment.Directory = configuration.Directory
	if override.Directory != nil {
		environment.Directory = *override.Directory
	}
	if flags.Lookup("directory").Changed {
		environment.Directory = flags.Lookup("directory").Value.String()
	}

	return &environment, nil
}
