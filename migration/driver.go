package migration

import (
	"errors"
	"github.com/elwinar/rambler/configuration"
)

// Driver is the interface used by the program to interact with the migration
// table in database
type Driver interface {
	MigrationTableExists() (bool, error)
	CreateMigrationTable() error
}

// Constructor is the function type used to create drivers
type Constructor func(configuration.Environment) (Driver, error)

// The various errors returned by the package
var (
	ErrDriverAlreadyRegistered = errors.New("driver already registered")
	ErrDriverNotRegistered     = errors.New("driver not registered")
)

var (
	constructors map[string]Constructor
)

func init() {
	constructors = make(map[string]Constructor)
}

// Register register a constructor for a driver
func RegisterDriver(name string, constructor Constructor) error {
	return registerDriver(name, constructor, constructors)
}

func registerDriver(name string, constructor Constructor, constructors map[string]Constructor) error {
	if _, found := constructors[name]; found {
		return ErrDriverAlreadyRegistered
	}

	constructors[name] = constructor
	return nil
}

// Get initialize a driver from the given name and options
func GetDriver(env configuration.Environment) (Driver, error) {
	return getDriver(env, constructors)
}

func getDriver(env configuration.Environment, constructors map[string]Constructor) (Driver, error) {
	constructor, found := constructors[env.Driver]
	if !found {
		return nil, ErrDriverNotRegistered
	}

	return constructor(env)
}
