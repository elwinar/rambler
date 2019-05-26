package main

import (
	"reflect"
	"testing"

	"github.com/elwinar/rambler/driver"
	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	var cases = map[string]struct {
		input  string
		err    bool
		output Configuration
	}{
		"not_found": {
			input:  "testdata/notfound.json",
			err:    true,
			output: Configuration{},
		},
		"invalid": {
			input:  "testdata/invalid.json",
			err:    true,
			output: Configuration{},
		},
		"valid_json": {
			input: "testdata/valid.json",
			err:   false,
			output: Configuration{
				Environment: Environment{
					Driver: "mysql",
					Config: driver.Config{
						Protocol: "tcp",
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "",
						Database: "rambler_default",
						Table:    "migrations",
					},
					Directory: ".",
				},
				Environments: map[string]Environment{
					"testing": {
						Config: driver.Config{
							Database: "rambler_testing",
						},
					},
					"development": {
						Config: driver.Config{
							Database: "rambler_development",
						},
					},
					"production": {
						Config: driver.Config{
							Database: "rambler_production",
						},
					},
				},
			},
		},
		"valid_hjson": {
			input: "testdata/valid.hjson",
			err:   false,
			output: Configuration{
				Environment: Environment{
					Driver: "mysql",
					Config: driver.Config{
						Protocol: "tcp",
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "",
						Database: "rambler_default",
						Table:    "migrations",
					},
					Directory: ".",
				},
				Environments: map[string]Environment{
					"testing": {
						Config: driver.Config{
							Database: "rambler_testing",
						},
					},
					"development": {
						Config: driver.Config{
							Database: "rambler_development",
						},
					},
					"production": {
						Config: driver.Config{
							Database: "rambler_production",
						},
					},
				},
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			cfg, err := Load(c.input)
			if (err != nil) != c.err {
				t.Errorf("%s: unexpected error: %s", n, err)
				return
			}

			if !reflect.DeepEqual(cfg, c.output) {
				t.Errorf("%s: unexpected output: %s\n", n, cmp.Diff(c.output, cfg))
			}
		})
	}
}

func TestConfigurationEnv(t *testing.T) {
	var cases = map[string]struct {
		input  string
		err    bool
		output Environment
	}{
		"unknown": {
			input: "unknown",
			err:   true,
		},
		"default": {
			input: "default",
			output: Environment{
				Driver: "mysql",
				Config: driver.Config{
					Protocol: "tcp",
					Host:     "localhost",
					Port:     3306,
					User:     "root",
					Password: "",
					Database: "rambler_default",
					Table:    "migrations",
				},
				Directory: ".",
			},
		},
		"testing": {
			input: "testing",
			output: Environment{
				Driver: "mysql",
				Config: driver.Config{
					Protocol: "tcp",
					Host:     "localhost",
					Port:     3306,
					User:     "root",
					Password: "",
					Database: "rambler_testing",
					Table:    "migrations",
				},
				Directory: ".",
			},
		},
	}

	for n, c := range cases {
		t.Run(n, func(t *testing.T) {
			cfg := Configuration{
				Environment: Environment{
					Driver: "mysql",
					Config: driver.Config{
						Protocol: "tcp",
						Host:     "localhost",
						Port:     3306,
						User:     "root",
						Password: "",
						Database: "rambler_default",
						Table:    "migrations",
					},
					Directory: ".",
				},
				Environments: map[string]Environment{
					"testing": {
						Config: driver.Config{
							Database: "rambler_testing",
						},
					},
					"development": {
						Config: driver.Config{
							Database: "rambler_development",
						},
					},
					"production": {
						Config: driver.Config{
							Database: "rambler_production",
						},
					},
				},
			}

			env, err := cfg.Env(c.input)
			if (err != nil) != c.err {
				t.Errorf("%s: unexpected error: %s", n, err)
				return
			}

			if !reflect.DeepEqual(env, c.output) {
				t.Errorf("%s: unexpected output: %s", n, cmp.Diff(c.output, env))
			}
		})
	}
}
