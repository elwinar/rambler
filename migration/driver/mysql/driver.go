package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/elwinar/rambler/migration/driver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

func init() {
	driver.Register("mysql", NewDriver)
}

var (
	ErrUnknownDatabase      = errors.New("invalid DSN or unreachable database")
	ErrNotInitializedDriver = errors.New("driver not initialized")
)

type Driver struct {
	db     *sqlx.DB
	schema string
}

// NewDriver return a mysql driver that implements the driver interface
func NewDriver(dsn string) (driver.Driver, error) {
	return newDriver(dsn, sqlx.Connect)
}

type connecter func(string, string) (*sqlx.DB, error)

func newDriver(dsn string, connect connecter) (*Driver, error) {
	db, err := connect("mysql", dsn)
	if err != nil {
		return nil, ErrUnknownDatabase
	}

	lastSlash := strings.LastIndex(dsn, "/")
	schema := dsn[lastSlash+1:]
	firstQuestion := strings.Index(dsn, "?")
	if firstQuestion != -1 {
		schema = schema[:firstQuestion-1]
	}

	return &Driver{
		db:     db,
		schema: schema,
	}, nil
}

// MigrationTableExists return wheter or not a tbale named migration exists in the
// target schema.
func (d Driver) MigrationTableExists() (bool, error) {
	if d.db == nil {
		return false, ErrNotInitializedDriver
	}

	var table struct {
		Name string `db:'name'`
	}
	
	err := d.db.Get(&table, fmt.Sprintf(`SELECT table_name as name FROM information_schema.tables WHERE table_schema = '%s' AND table_name = 'migrations'`, d.schema))

	if err == sql.ErrNoRows {
		return false, nil
	}
	
	if err != nil {
		return false, err
	}
	
	return true, nil
}

func (d Driver) CreateMigrationTable() error {
	return nil
}
