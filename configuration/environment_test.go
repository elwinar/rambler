package configuration

import (
	"reflect"
	"testing"
)

func Test_Env_UnknownEnvironment(t *testing.T) {
	_, err := Env("test", Configuration{}, map[string]string{})

	if err == nil {
		t.Error("didn't failed on unknown environment")
	}

	if err.Error() != "unkown environment test" {
		t.Error("didn't returned expected error:", err)
	}
}

func Test_Env_Override(t *testing.T) {
	var sEnv = "env"
	var iEnv = uint64(1)
	type Case struct {
		Name          string
		Configuration Configuration
		Flags         map[string]string
		Env           *Environment
	}
	var cases = []Case{
		Case{
			Name: "test",
			Configuration: Configuration{
				Driver:    "base",
				Protocol:  "base",
				Host:      "base",
				Port:      0,
				User:      "base",
				Password:  "base",
				Database:  "base",
				Directory: "base",
				Environments: map[string]RawEnvironment{
					"test": RawEnvironment{
						Driver:    &sEnv,
						Protocol:  &sEnv,
						Host:      &sEnv,
						Port:      &iEnv,
						User:      &sEnv,
						Password:  &sEnv,
						Database:  &sEnv,
						Directory: &sEnv,
					},
				},
			},
			Flags: map[string]string{
				"driver":    "flag",
				"protocol":  "flag",
				"host":      "flag",
				"port":      "2",
				"user":      "flag",
				"password":  "flag",
				"database":  "flag",
				"directory": "flag",
			},
			Env: &Environment{
				Driver:    "flag",
				Protocol:  "flag",
				Host:      "flag",
				Port:      2,
				User:      "flag",
				Password:  "flag",
				Database:  "flag",
				Directory: "flag",
			},
		},
		Case{
			Name: "test",
			Configuration: Configuration{
				Driver:    "base",
				Protocol:  "base",
				Host:      "base",
				Port:      0,
				User:      "base",
				Password:  "base",
				Database:  "base",
				Directory: "base",
				Environments: map[string]RawEnvironment{
					"test": RawEnvironment{
						Driver:    &sEnv,
						Protocol:  &sEnv,
						Host:      &sEnv,
						Port:      &iEnv,
						User:      &sEnv,
						Password:  &sEnv,
						Database:  &sEnv,
						Directory: &sEnv,
					},
				},
			},
			Flags: map[string]string{},
			Env: &Environment{
				Driver:    "env",
				Protocol:  "env",
				Host:      "env",
				Port:      1,
				User:      "env",
				Password:  "env",
				Database:  "env",
				Directory: "env",
			},
		},
		Case{
			Name: "test",
			Configuration: Configuration{
				Driver:    "base",
				Protocol:  "base",
				Host:      "base",
				Port:      0,
				User:      "base",
				Password:  "base",
				Database:  "base",
				Directory: "base",
				Environments: map[string]RawEnvironment{
					"test": RawEnvironment{},
				},
			},
			Flags: map[string]string{},
			Env: &Environment{
				Driver:    "base",
				Protocol:  "base",
				Host:      "base",
				Port:      0,
				User:      "base",
				Password:  "base",
				Database:  "base",
				Directory: "base",
			},
		},
	}

	for n, c := range cases {
		e, err := Env(c.Name, c.Configuration, c.Flags)

		if err != nil {
			t.Error("unexpected error")
		}

		if !reflect.DeepEqual(e, c.Env) {
			t.Errorf("uncorrectly overriden environment for case %d", n)
		}
	}
}
