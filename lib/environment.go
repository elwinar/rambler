package lib

import (
	"errors"
	"github.com/elwinar/viper"
	"path/filepath"
)

var (
	Env Environment
)

// RawEnvironment is the potentially-sparse type used in the configuration.
// Nil attributes are to be replaced by the default value, which is defined in
// the default environment
type RawEnvironment struct {
	Driver *string
	Protocol *string
	Host *string
	Port *int
	User *string
	Password *string
	Database *string
	Migrations *string
	Seeds *string
}

// Environment is the complete and valid environment issued by the GetEnvironment
// function.
type Environment struct {
	Driver string
	Protocol string
	Host string
	Port int
	User string
	Password string
	Database string
	Migrations string
	Seeds string
}

// GetEnvironment parse the configuration to extract a full environment configuration,
// and check whether it is valid or not.
func LoadEnvironment() (error) {
	var env = viper.GetString("environment")
	var err error
	var rawEnvs map[string]RawEnvironment = make(map[string]RawEnvironment)
	
	// Get the raw environments map
	err = viper.MarshalKey("environments", &rawEnvs)
	if err != nil {
		return err
	}
	
	// Fill the environment with the default values (taken from the configuration)
	Env.Driver = viper.GetString("driver")
	Env.Protocol = viper.GetString("protocol")
	Env.Host = viper.GetString("host")
	Env.Port = viper.GetInt("port")
	Env.User = viper.GetString("user")
	Env.Password = viper.GetString("password")
	Env.Database = viper.GetString("database")
	Env.Migrations = viper.GetString("migrations")
	Env.Seeds = viper.GetString("seeds")
	
	// If requesting the default environment, return now
	if env == "default" {
		return nil
	}
	
	// Check if the requested environment is in the map
	raw, found := rawEnvs[env]
	if !found  {
		return errors.New("unknown environment " + env)
	}
	
	// Override default environment with non-null values from the raw environment
	// TODO Find or write a lib to do it for me, preventing this ugly if-medley
	if raw.Driver != nil {
		Env.Driver = *raw.Driver
	}
	if raw.Protocol != nil {
		Env.Protocol = *raw.Protocol
	}
	if raw.Host != nil {
		Env.Host = *raw.Host
	}
	if raw.Port != nil {
		Env.Port = *raw.Port
	}
	if raw.User != nil {
		Env.User = *raw.User
	}
	if raw.Password != nil {
		Env.Password = *raw.Password
	}
	if raw.Database != nil {
		Env.Database = *raw.Database
	}
	if raw.Migrations != nil {
		Env.Migrations = *raw.Migrations
	}
	if raw.Seeds != nil {
		Env.Seeds = *raw.Seeds
	}
	
	return nil
}

func (env Environment) MigrationsDir() string {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), env.Migrations)
}
