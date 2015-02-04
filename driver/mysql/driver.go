package mysql

import (
	"database/sql"
	"github.com/elwinar/rambler/driver"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	driver.Register("mysql", Driver{})
}

type Driver struct {}

func (d Driver) New(dsn, schema string) (driver.Conn, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return Conn{
		db: db,
		schema: schema,
	}, nil
}

type Conn struct{
	db *sql.DB
	schema string
}

func (c Conn) MigrationTableExists() (bool, error) {
	var table string
	err := c.db.QueryRow(`select table_name from information_schema.tables where table_schema = ? and table_name = ?`, c.schema, "migrations").Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	
	if err == sql.ErrNoRows {
		return false, nil
	}
	
	return true, nil
}

func (c Conn) CreateMigrationTable() error {
	_, err := c.db.Exec(`CREATE TABLE migrations ( version BIGINT UNSIGNED NOT NULL PRIMARY KEY, description VARCHAR(255) NOT NULL, applied_at DATETIME NOT NULL ) DEFAULT CHARSET=utf8`)
	return err
}

func (c Conn) ListAppliedMigrations() ([]uint64, error) {
	rows, err := c.db.Query(`SELECT version FROM migrations ORDER BY version ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []uint64
	for rows.Next() {
		var version uint64
		err := rows.Scan(&version)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	return versions, nil
}

func (c Conn) SetMigrationApplied(version uint64, description string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (version, description, applied_at) VALUES (?, ?, NOW())`, version, description)
	return err
}

func (c Conn) UnsetMigrationApplied(version uint64) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE version = ?`, version)
	return err
}

func (c Conn) Exec(query string) error {
	_, err := c.db.Exec(query)
	return err
}
