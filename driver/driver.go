package driver

import (
	"fmt"
)

// Driver is the interface used by the program to initialize the database
// connection.
type Driver interface {
	New(Config) (Conn, error)
}

var drivers = make(map[string]Driver)

func Register(name string, driver Driver) error {
	if _, found := drivers[name]; found {
		return fmt.Errorf(`driver "%s" already registered`, name)
	}

	if driver == nil {
		return fmt.Errorf(`not a valid driver`)
	}

	drivers[name] = driver
	return nil
}

// Get returns the requested driver.
func Get(driver string) (Driver, error) {
	d, found := drivers[driver]
	if !found {
		return nil, fmt.Errorf(`driver "%s" not registered`, driver)
	}

	return d, nil
}

type Config struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     uint64 `json:"port"`
	User     string `json:"user"`
	Role     string `json:"role"` // For PostgreSQL
	Password string `json:"password"`
	Database string `json:"database"`
	Schema   string `json:"schema"` // For PostgreSQL
	Table    string `json:"table"`
	SSLMode  string `json:"sslmode"` // For PostgreSQL
}

// Conn is the interface used by the program to manipulate the database.
type Conn interface {
	HasTable() (bool, error)
	CreateTable() error
	GetApplied() ([]string, error)
	AddApplied(string) error
	RemoveApplied(string) error
	Execute(string) error
}
