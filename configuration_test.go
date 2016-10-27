package main

import (
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	var cases = []struct {
		input  string
		err    bool
		output Configuration
	}{
		{
			input:  "test/notfound.json",
			err:    true,
			output: Configuration{},
		},
		{
			input:  "test/invalid.json",
			err:    true,
			output: Configuration{},
		},
		{
			input: "test/valid.json",
			err:   false,
			output: Configuration{
				Environment: Environment{
					Driver:    "mysql",
					Protocol:  "tcp",
					Host:      "localhost",
					Port:      3306,
					User:      "root",
					Password:  "",
					Database:  "rambler_default",
					Directory: ".",
				},
				Environments: map[string]Environment{
					"testing": {
						Database: "rambler_testing",
					},
					"development": {
						Database: "rambler_development",
					},
					"production": {
						Database: "rambler_production",
					},
				},
			},
		},
		{
			input: "test/valid.hjson",
			err:   false,
			output: Configuration{
				Environment: Environment{
					Driver:    "mysql",
					Protocol:  "tcp",
					Host:      "localhost",
					Port:      3306,
					User:      "root",
					Password:  "",
					Database:  "rambler_default",
					Directory: ".",
				},
				Environments: map[string]Environment{
					"testing": {
						Database: "rambler_testing",
					},
					"development": {
						Database: "rambler_development",
					},
					"production": {
						Database: "rambler_production",
					},
				},
			},
		},
	}

	for n, c := range cases {
		cfg, err := Load(c.input)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(cfg, c.output) {
			t.Error("case", n, "got unexpected output:", cfg)
		}
	}
}

func TestConfigurationEnv(t *testing.T) {
	var cases = []struct {
		input  string
		err    bool
		output Environment
	}{
		{
			input:  "unknown",
			err:    true,
			output: Environment{},
		},
		{
			input: "default",
			err:   false,
			output: Environment{
				Driver:    "mysql",
				Protocol:  "tcp",
				Host:      "localhost",
				Port:      3306,
				User:      "root",
				Password:  "",
				Database:  "rambler_default",
				Directory: ".",
			},
		},
		{
			input: "testing",
			err:   false,
			output: Environment{
				Driver:    "mysql",
				Protocol:  "tcp",
				Host:      "localhost",
				Port:      3306,
				User:      "root",
				Password:  "",
				Database:  "rambler_testing",
				Directory: ".",
			},
		},
	}

	for n, c := range cases {
		cfg := Configuration{
			Environment: Environment{
				Driver:    "mysql",
				Protocol:  "tcp",
				Host:      "localhost",
				Port:      3306,
				User:      "root",
				Password:  "",
				Database:  "rambler_default",
				Directory: ".",
			},
			Environments: map[string]Environment{
				"testing": {
					Database: "rambler_testing",
				},
				"development": {
					Database: "rambler_development",
				},
				"production": {
					Database: "rambler_production",
				},
			},
		}

		env, err := cfg.Env(c.input)
		if (err != nil) != c.err {
			t.Error("case", n, "got unexpected error:", err)
			continue
		}

		if !reflect.DeepEqual(env, c.output) {
			t.Error("case", n, "got unexpected output:", cfg)
		}
	}
}
