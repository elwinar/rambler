package main

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"io/ioutil"
)

//go:generate ffjson $GOFILE

// Configuration is the configuration type
type Configuration struct {
	Environment
	Environments map[string]Environment `json:"environments"`
}

// Load open, read and parse the given configuration file
func Load(filename string) (Configuration, error) {
	var c Configuration

	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(raw, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// Env return the requested environment from the configuration
func (c Configuration) Env(name string) (Environment, error) {
	environment := c.Environment
	
	if name == "default" {
		return environment, nil
	}
	
	overrides, found := c.Environments[name]
	if !found {
		return Environment{}, fmt.Errorf("unknown environment %s", name)
	}

	_ = mergo.Merge(&environment, overrides) // No error can possibly occur here
	return environment, nil
}
