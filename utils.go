package main

import (
	"encoding/json"
	"github.com/codegangsta/cli"
	"github.com/elwinar/rambler/configuration"
	"io/ioutil"
	"log"
	"os"
)

func bootstrap(c *cli.Context) (*configuration.Environment, *log.Logger, *log.Logger, *log.Logger, error) {
	var Debug, Info, Error *log.Logger
	var C configuration.Configuration = configuration.Configuration{
		Environment: configuration.Environment{
			Driver:    "mysql",
			Protocol:  "tcp",
			Host:      "localhost",
			Port:      3306,
			User:      "root",
			Password:  "",
			Database:  "",
			Directory: ".",
		},
	}

	var flags int
	if c.GlobalBool("verbose") {
		flags = log.Ltime | log.Lshortfile
	} else {
		flags = log.Ltime
	}

	Error = log.New(os.Stdout, "error ", flags)

	if c.GlobalBool("debug") {
		Debug = log.New(os.Stdout, "debug ", flags)
	} else {
		Debug = log.New(ioutil.Discard, "", flags)
	}

	if c.GlobalBool("quiet") {
		Info = log.New(ioutil.Discard, "", flags)
	} else {
		Info = log.New(os.Stdout, "info ", flags)
	}

	raw, err := ioutil.ReadFile(c.GlobalString("configuration"))
	if err != nil {
		return nil, Debug, Info, Error, err
	}

	err = json.Unmarshal(raw, &C)
	if err != nil {
		return nil, Debug, Info, Error, err
	}

	Env, err := C.Env(c.GlobalString("environment"))
	if err != nil {
		return nil, Debug, Info, Error, err
	}

	return Env, Debug, Info, Error, err
}
