package configuration

import (
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
	"io/ioutil"
)

//go:generate ffjson $GOFILE

// Environment is the execution environment of a command. It contains every information
// about the database and migrations to use.
type Environment struct {
	Driver    string `json:"driver"`
	Protocol  string `json:"protocol"`
	Host      string `json:"host"`
	Port      uint64 `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	Directory string `json:"directory"`
}

// Configuration is the configuration type
type Configuration struct {
	Environment
	Environments map[string]Environment `json:"environments"`
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
