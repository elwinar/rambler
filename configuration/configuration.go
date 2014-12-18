package configuration

import (
	"fmt"
	"strconv"
)

// Configuration is the rambler configuration type as loaded from the configuration
// file and extended by the command-line.
type Configuration struct {
	Driver       string
	Protocol     string
	Host         string
	Port         uint64
	User         string
	Password     string
	Database     string
	Directory    string
	Environments map[string]RawEnvironment
}

// GetEnvironment return the requested environment from the configuration, with
// overrides from the given flagset.
func (c Configuration) Env(name string, flags map[string]string) (*Environment, error) {
	override, found := c.Environments[name]
	if !found && name != "default" {
		return nil, fmt.Errorf(errUnknownEnvironment, name)
	}

	var environment Environment

	environment.Driver = c.Driver
	if override.Driver != nil {
		environment.Driver = *override.Driver
	}
	if v, found := flags["driver"]; found {
		environment.Driver = v
	}

	environment.Protocol = c.Protocol
	if override.Protocol != nil {
		environment.Protocol = *override.Protocol
	}
	if v, found := flags["protocol"]; found {
		environment.Protocol = v
	}

	environment.Host = c.Host
	if override.Host != nil {
		environment.Host = *override.Host
	}
	if v, found := flags["host"]; found {
		environment.Host = v
	}

	environment.Port = c.Port
	if override.Port != nil {
		environment.Port = *override.Port
	}
	if v, found := flags["port"]; found {
		environment.Port, _ = strconv.ParseUint(v, 10, 64)
	}

	environment.User = c.User
	if override.User != nil {
		environment.User = *override.User
	}
	if v, found := flags["user"]; found {
		environment.User = v
	}

	environment.Password = c.Password
	if override.Password != nil {
		environment.Password = *override.Password
	}
	if v, found := flags["password"]; found {
		environment.Password = v
	}

	environment.Database = c.Database
	if override.Database != nil {
		environment.Database = *override.Database
	}
	if v, found := flags["database"]; found {
		environment.Database = v
	}

	environment.Directory = c.Directory
	if override.Directory != nil {
		environment.Directory = *override.Directory
	}
	if v, found := flags["directory"]; found {
		environment.Directory = v
	}

	return &environment, nil
}
