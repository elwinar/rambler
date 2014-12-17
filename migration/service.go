package migration

import (
	"fmt"
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration/driver"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Service gather operations to manipulate migrations table and migrations on
// the filesystem.
type Service interface {
	driver.Driver
	ListAvailableMigrations() ([]uint64, error)
}

type service struct {
	driver.Driver
	env configuration.Environment
}

// NewService initialize a new service with the given informations
func NewService(env configuration.Environment) (Service, error) {
	return newService(env, os.Stat, driver.Get)
}

func newService(env configuration.Environment, stat stater, get driverConstructor) (*service, error) {
	if _, err := stat(env.Directory); err != nil {
		return nil, fmt.Errorf(errUnavailableDirectory, env.Directory, err.Error())
	}

	driver, err := get(env)
	if err != nil {
		return nil, fmt.Errorf(errDriverError, env.Driver, err.Error())
	}

	return &service{
		Driver: driver,
		env:    env,
	}, nil
}

// ListAvailableMigrations return the list migrations in the environment's directory
func (s service) ListAvailableMigrations() ([]uint64, error) {
	return listAvailableMigrations(s.env, filepath.Glob)
}

func listAvailableMigrations(env configuration.Environment, glob glober) ([]uint64, error) {
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
