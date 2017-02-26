package main

import (
	"fmt"

	"github.com/elwinar/rambler/log"
	"github.com/kelseyhightower/envconfig"
	"github.com/urfave/cli"
)

// Bootstrap do the initialization job, and finish by setting the
// `service` global var that will be used by other commands.
func Bootstrap(ctx *cli.Context) error {
	return bootstrap(ctx.GlobalString("configuration"), ctx.GlobalString("environment"), ctx.GlobalBool("debug"))
}

func bootstrap(configuration, environment string, debug bool) error {
	logger = log.NewLogger(func(l *log.Logger) {
		l.PrintDebug = debug
	})

	logger.Debug("loading configuration from %s", configuration)
	cfg, err := Load(configuration)
	if err != nil {
		return fmt.Errorf("unable to load configuration from file: %s", err)
	}

	err = envconfig.Process("rambler", &cfg)
	if err != nil {
		return fmt.Errorf("unable to load configuration from env: %s", err)
	}

	logger.Debug("loading environment", environment)
	env, err := cfg.Env(environment)
	if err != nil {
		return fmt.Errorf("unable to load requested environment: %s", err)
	}

	logger.Debug("initializing service")
	service, err = NewService(env)
	if err != nil {
		return fmt.Errorf("unable to initialize the migration service: %s", err)
	}

	return nil
}
