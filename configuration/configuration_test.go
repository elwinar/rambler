package configuration

import (
	"reflect"
	"testing"
)

var (
	defaults Environment = Environment{
		Driver: "mysql",
		Protocol: "tcp",
		Host: "localhost",
		Port: 3306,
		User: "root",
		Password: "",
		Database: "rambler_default",
		Directory: ".",
	}
	good Configuration = Configuration{
		Environment: defaults,
		Environments: map[string]Environment{
			"testing": Environment{
				Database: "rambler_testing",
			},
			"development": Environment{
				Database: "rambler_development",
			},
			"production": Environment{
				Database: "rambler_production",
			},
		},
	}
)

func Test_Load_NotFound(t *testing.T) {
	_, err := Load("test/notfound.json")
	if err == nil {
		t.Fail()
	}
}

func Test_Load_InvalidSyntaxd(t *testing.T) {
	_, err := Load("test/bad.json")
	if err == nil {
		t.Fail()
	}
}

func Test_Load_OK(t *testing.T) {
	c, err := Load("test/good.json")
	if err != nil {
		t.Fail()
	}
	
	if !reflect.DeepEqual(c, good) {
		t.Fail()
	}
}

func Test_Env_Unknown(t *testing.T) {
	_, err := good.Env("unknown")
	if err == nil {
		t.Fail()
	}
}

func Test_Env_Override(t *testing.T) {
	e, err := good.Env("default")
	if err != nil {
		t.Fail()
	}
	
	if !reflect.DeepEqual(e, defaults) {
		t.Fail()
	}
}
