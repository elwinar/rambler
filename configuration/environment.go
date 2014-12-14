package configuration

import (
	"fmt"
	"strconv"
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
func Env(name string, configuration Configuration, flags map[string]string) (*Environment, error) {
	override, found := configuration.Environments[name]
	if !found {
		return nil, fmt.Errorf(errUnknownEnvironment, name)
	}

	var environment Environment

	environment.Driver = configuration.Driver
	if override.Driver != nil {
		environment.Driver = *override.Driver
	}
	if v, found := flags["driver"]; found {
		environment.Driver = v
	}

	environment.Protocol = configuration.Protocol
	if override.Protocol != nil {
		environment.Protocol = *override.Protocol
	}
	if v, found := flags["protocol"]; found {
		environment.Protocol = v
	}

	environment.Host = configuration.Host
	if override.Host != nil {
		environment.Host = *override.Host
	}
	if v, found := flags["host"]; found {
		environment.Host = v
	}

	environment.Port = configuration.Port
	if override.Port != nil {
		environment.Port = *override.Port
	}
	if v, found := flags["port"]; found {
		environment.Port, _ = strconv.ParseUint(v, 10, 64)
	}

	environment.User = configuration.User
	if override.User != nil {
		environment.User = *override.User
	}
	if v, found := flags["user"]; found {
		environment.User = v
	}

	environment.Password = configuration.Password
	if override.Password != nil {
		environment.Password = *override.Password
	}
	if v, found := flags["password"]; found {
		environment.Password = v
	}

	environment.Database = configuration.Database
	if override.Database != nil {
		environment.Database = *override.Database
	}
	if v, found := flags["database"]; found {
		environment.Database = v
	}

	environment.Directory = configuration.Directory
	if override.Directory != nil {
		environment.Directory = *override.Directory
	}
	if v, found := flags["directory"]; found {
		environment.Directory = v
	}

	return &environment, nil
}
