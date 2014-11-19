package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/elwinar/rambler/lib/driver"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

func init() {
	driver.Register("mysql", NewDriver)
}

var (
	ErrUnknownDatabase      = errors.New("invalid DSN or unreachable database")
	ErrNoSchema             = errors.New("no database schema")
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
	firstQuestion := strings.Index(schema, "?")

	if firstQuestion == 0 {
		return nil, ErrNoSchema
	}

	if firstQuestion != -1 {
		schema = schema[:firstQuestion]
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

	err := d.db.Get(&table, fmt.Sprintf(`
		SELECT table_name as name 
		FROM information_schema.tables
		WHERE table_schema = '%s' 
		AND table_name = 'migrations'
	`, d.schema))

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (d Driver) CreateMigrationTable() error {
	if d.db == nil {
		return ErrNotInitializedDriver
	}

	_, err := d.db.Exec(`
		CREATE TABLE migrations ( 
			version BIGINT UNSIGNED NOT NULL PRIMARY KEY,
			description CHAR(255) NOT NULL,
			applied_at DATETIME NOT NULL
		) DEFAULT CHARSET=utf8
	`)

	return err
}
