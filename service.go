package main

import (
	"fmt"
	"github.com/elwinar/rambler/driver"
	_ "github.com/elwinar/rambler/driver/mysql"
	_ "github.com/elwinar/rambler/driver/postgresql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Service is the interface that gather operations to manipulate migrations table
// and migrations on the filesystem.
type Service interface {
	driver.Conn
	ListAvailableMigrations() ([]uint64, error)
}

// CoreService is the basic implementation of the Service interface
type CoreService struct {
	driver.Conn
	env Environment
}

// NewService initialize a new service with the given informations
func NewService(env Environment) (Service, error) {
	return newService(env, os.Stat, driver.Get)
}

func newService(env Environment, stat stater, get connConstructor) (*CoreService, error) {
	if _, err := stat(env.Directory); err != nil {
		return nil, fmt.Errorf(errUnavailableDirectory, env.Directory, err.Error())
	}

	conn, err := get(env.Driver, env.DSN(), env.Database)
	if err != nil {
		return nil, fmt.Errorf(errDriverError, env.Driver, err.Error())
	}

	return &CoreService{
		Conn: conn,
		env:  env,
	}, nil
}

// ListAvailableMigrations return the list migrations in the environment's directory
func (s CoreService) ListAvailableMigrations() ([]uint64, error) {
	return listAvailableMigrations(s.env, filepath.Glob)
}

func listAvailableMigrations(env Environment, glob glober) ([]uint64, error) {
	raw, err := glob(filepath.Join(env.Directory, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf(errUnavailableDirectory, env.Directory, err.Error())
	}

	versions := make([]uint64, 0)
	for _, r := range raw {
		file := filepath.Base(r)

		chunks := strings.SplitN(file, "_", 2)

		if len(chunks) != 2 {
			continue
		}

		version, err := strconv.ParseUint(chunks[0], 10, 64)
		if err != nil {
			continue
		}

		versions = append(versions, version)
	}

	return versions, nil
}
