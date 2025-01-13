package main

import (
	"fmt"
	"io/ioutil"

	"github.com/client9/xson/hjson"
	"dario.cat/mergo"
)

// Configuration is the configuration type
type Configuration struct {
	Environment
	Environments map[string]Environment `json:"environments"`
}

// Load open, read and parse the given configuration file
func Load(filename string) (Configuration, error) {
	var c Configuration
	c.Table = "migrations"

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	err = hjson.Unmarshal(raw, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}

// Env return the requested environment from the configuration
func (c Configuration) Env(name string) (Environment, error) {
	if name == "default" {
		return c.Environment, nil
	}

	env, found := c.Environments[name]
	if !found {
		return Environment{}, fmt.Errorf("unknown environment %s", name)
	}

	_ = mergo.Merge(&env, c.Environment) // No error can possibly occur here
	return env, nil
}
