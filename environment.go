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
	Directory string `json:"directory"`
}

// DSN return the connection string for the current environment
func (e Environment) DSN() string {
	switch e.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", e.User, e.Password, e.Protocol, e.Host, e.Port, e.Database)
	case "postgresql":
		return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", e.User, e.Password, e.Host, e.Port, e.Database)
	default:
		return ""
	}
}
