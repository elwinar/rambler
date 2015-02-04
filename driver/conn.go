package driver

// Conn is the interface used by the program to manipulate the migration table.
type Conn interface {
	MigrationTableExists() (bool, error)
	CreateMigrationTable() error
	ListAppliedMigrations() ([]uint64, error)
	SetMigrationApplied(version uint64, description string) error
	UnsetMigrationApplied(version uint64) error
	Exec(query string) error
}
