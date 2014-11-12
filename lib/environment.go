package lib

import (
	"errors"
	"github.com/elwinar/cobra"
	"github.com/elwinar/viper"
	"path/filepath"
	"strconv"
)

var (
	Env Environment
)

// RawEnvironment is the potentially-sparse type used in the configuration.
// Nil attributes are to be replaced by the default value, which is defined in
// the default environment
type RawEnvironment struct {
	Driver     *string
	Protocol   *string
	Host       *string
	Port       *int
	User       *string
	Password   *string
	Database   *string
	Migrations *string
}

// Environment is the complete and valid environment issued by the GetEnvironment
// function.
type Environment struct {
	Driver     string
	Protocol   string
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	Migrations string
}

// GetEnvironment parse the configuration to extract a full environment configuration,
// and check whether it is valid or not.
func LoadEnvironment(cmd *cobra.Command) error {
	var env = viper.GetString("environment")
	var err error
	var rawEnvs map[string]RawEnvironment = make(map[string]RawEnvironment)

	// Get the raw environments map
	err = viper.MarshalKey("environments", &rawEnvs)
	if err != nil {
		return err
	}

	// Get the default environment value
	var globalEnv RawEnvironment
	err = viper.Marshal(&globalEnv)
	if err != nil {
		return err
	}

	// Set the default environment value
	rawEnvs["default"] = globalEnv

	// Look for the requested environment
	raw, found := rawEnvs[env]
	if !found {
		return errors.New("unknown environment " + env)
	}

	// Set the driver value
	if cmd.Flag("driver") != nil && cmd.Flag("driver").Changed {
		Env.Driver = cmd.Flag("driver").Value.String()
	} else if raw.Driver != nil {
		Env.Driver = *raw.Driver
	} else {
		Env.Driver = *globalEnv.Driver
	}

	// Set the protocol value
	if cmd.Flag("protocol") != nil && cmd.Flag("protocol").Changed {
		Env.Protocol = cmd.Flag("protocol").Value.String()
	} else if raw.Protocol != nil {
		Env.Protocol = *raw.Protocol
	} else {
		Env.Protocol = *globalEnv.Protocol
	}

	// Set the host value
	if cmd.Flag("host") != nil && cmd.Flag("host").Changed {
		Env.Host = cmd.Flag("host").Value.String()
	} else if raw.Host != nil {
		Env.Host = *raw.Host
	} else {
		Env.Host = *globalEnv.Host
	}

	// Set the port value
	if cmd.Flag("port") != nil && cmd.Flag("port").Changed {
		Env.Port, _ = strconv.Atoi(cmd.Flag("port").Value.String())
	} else if raw.Port != nil {
		Env.Port = *raw.Port
	} else {
		Env.Port = *globalEnv.Port
	}

	// Set the user value
	if cmd.Flag("user") != nil && cmd.Flag("user").Changed {
		Env.User = cmd.Flag("user").Value.String()
	} else if raw.User != nil {
		Env.User = *raw.User
	} else {
		Env.User = *globalEnv.User
	}

	// Set the password value
	if cmd.Flag("password") != nil && cmd.Flag("password").Changed {
		Env.Password = cmd.Flag("password").Value.String()
	} else if raw.Password != nil {
		Env.Password = *raw.Password
	} else {
		Env.Password = *globalEnv.Password
	}

	// Set the database value
	if cmd.Flag("database") != nil && cmd.Flag("database").Changed {
		Env.Database = cmd.Flag("database").Value.String()
	} else if raw.Database != nil {
		Env.Database = *raw.Database
	} else {
		Env.Database = *globalEnv.Database
	}

	// Set the migrations value
	if cmd.Flag("migrations") != nil && cmd.Flag("migrations").Changed {
		Env.Migrations = cmd.Flag("migrations").Value.String()
	} else if raw.Migrations != nil {
		Env.Migrations = *raw.Migrations
	} else {
		Env.Migrations = *globalEnv.Migrations
	}

	return nil
}

func (env Environment) MigrationsDir() string {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), env.Migrations)
}
