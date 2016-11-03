package main

import (
	"fmt"
	"github.com/urfave/cli"
)

// Bootstrap go the initialization job, and finish by setting the
// `service` global var that will be used by other commands.
func Bootstrap(ctx *cli.Context) error {
	return bootstrap(ctx.GlobalString("configuration"), ctx.GlobalString("environment"))
}

func bootstrap(configuration, environment string) error {
	cfg, err := Load(configuration)
	if err != nil {
		return fmt.Errorf("unable to load configuration file: %s", err)
	}

	env, err := cfg.Env(environment)
	if err != nil {
		return fmt.Errorf("unable to load requested environment: %s", err)
	}

	srv, err := NewService(env)
	if err != nil {
		return fmt.Errorf("unable to initialize the migration service: %s", err)
	}

	service = srv
	return nil
}
