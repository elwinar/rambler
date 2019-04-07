package main

import (
	"fmt"
)

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
	Schema    string `json:"schema"`
	Directory string `json:"directory"`
	Table     string `json:"table"`
}

// DSN return the connection string for the current environment
func (e Environment) DSN() string {
	switch e.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", e.User, e.Password, e.Protocol, e.Host, e.Port, e.Database)
	case "postgresql":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", e.Host, e.Port, e.User, e.Password, e.Database)
	case "sqlite":
		return e.Database
	default:
		return ""
	}
}
