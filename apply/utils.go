package apply

import (
	"github.com/elwinar/rambler/configuration"
	"github.com/elwinar/rambler/migration"
	"github.com/elwinar/rambler/migration/driver"
)

var (
	errNilTransaction         = "nil transaction"
	errOutOfOrder             = "out of order migration %d"
	errMissingMigration       = "missing migration %d"
	errInitializeNewService   = "unable to initialize migration service: %s"
	errMigrationTableCheck    = "unable to check for existence of the migration table: %s"
	errCreateMigrationTable   = "unable to create migration table: %s"
	errListApplied            = "error while listing applied migrations: %s"
	errListAvailable          = "error while listing available migrations: %s"
	errFilter                 = "error while filtering mgirations to apply: %s"
	errNewMigration           = "unable to read migration %d: %s"
	errStartTransactionFailed = "unable to start transaction: %s"
	errRollbackFailed         = "migration failed: %s; rollback failed: %s"
	errMigrationError         = "migration failed: %s"
	errCommitFailed           = "transaction commit failed: %s"
	errSetMigrationApplied    = "unable to set migration %d as applied: %s"
)

type filterer func([]uint64, []uint64) ([]uint64, error)
type applyer func([]string, driver.Tx) (error, error)
type serviceConstructor func(configuration.Environment) (migration.Service, error)
type migrationConstructor func(string, uint64) (*migration.Migration, error)
