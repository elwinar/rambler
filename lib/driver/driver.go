package driver

import (
	"errors"
)

// Driver is the interface used by the program to interact with the migration
// table in database
type Driver interface {
	MigrationTableExists() (bool, error)
	CreateMigrationTable() error
}

// Constructor is the function type used to create drivers
type Constructor func(string) (Driver, error)

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
func Register(name string, constructor Constructor) error {
	return register(name, constructor, constructors)
}

func register(name string, constructor Constructor, constructors map[string]Constructor) error {
	if _, found := constructors[name]; found {
		return ErrDriverAlreadyRegistered
	}

	constructors[name] = constructor
	return nil
}

// Get initialize a driver from the given name and options
func Get(name, options string) (Driver, error) {
	return get(name, options, constructors)
}

func get(name, options string, constructors map[string]Constructor) (Driver, error) {
	constructor, found := constructors[name]
	if !found {
		return nil, ErrDriverNotRegistered
	}

	return constructor(options)
}
