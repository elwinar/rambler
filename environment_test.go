package main

import (
	"reflect"
	"testing"
)

func TestEnvironmentDSN(t *testing.T) {
	var cases = []struct {
		driver string
		output string
	}{
		{
			driver: "unkown",
			output: "",
		},
		{
			driver: "mysql",
			output: "user:password@protocol(host:42)/database",
		},
		{
			driver: "postgresql",
			output: "user=user password=password host=host port=42 dbname=database sslmode=disable",
		},
	}

	for n, c := range cases {
		env := Environment{
			Driver:   c.driver,
			Protocol: "protocol",
			Host:     "host",
			Port:     42,
			User:     "user",
			Password: "password",
			Database: "database",
		}

		dsn := env.DSN()
		if !reflect.DeepEqual(dsn, c.output) {
			t.Error("case", n, "got unexpected output:", dsn)
		}
	}
}
