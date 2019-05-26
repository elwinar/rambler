package main

import (
	"github.com/elwinar/rambler/driver"
)

// Environment is the execution environment of a command. It contains every information
// about the database and migrations to use.
type Environment struct {
	Driver string `json:"driver"`
	driver.Config
	Directory string `json:"directory"`
}
