package apply

import (
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
	"github.com/elwinar/rambler/migration/driver"
)

var (
	errNilTransaction = "nil transaction"
	errOutOfOrder = "out of order migration %d"
	errMissingMigration = "missing migration %d"
)

type filterer func([]uint64, []uint64) ([]uint64, error)
type applyer func([]string, driver.Tx) (error, error)
type serviceConstructor func(configuration.Environment) (migration.Service, error)
type migrationConstructor func(string, uint64) (*migration.Migration, error)
