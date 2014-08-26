package main

import (
	"errors"
	"github.com/gonuts/flag"
)

const (
	errNoDatabase = `no database provided`
)

type Config struct {
	Protocol string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var (
	config Config
)

func GetConfig(flagset flag.FlagSet) error {
	// Get the database parameters
	config.Protocol = flagset.Lookup("protocol").Value.Get().(string)
	config.Host = flagset.Lookup("host").Value.Get().(string)
	config.Port = flagset.Lookup("port").Value.Get().(int)
	config.User = flagset.Lookup("user").Value.Get().(string)
	config.Password = flagset.Lookup("password").Value.Get().(string)
	config.Database = flagset.Lookup("db").Value.Get().(string)

	// Database is the only mandatory parameters, other parameters being defaulted
	// to the classic value for development platforms.
	if config.Database == "" {
		return errors.New(errNoDatabase)
	}

	return nil
}
