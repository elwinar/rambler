package mysql

import (
	"database/sql"

	"github.com/elwinar/rambler/driver"
	_ "github.com/go-sql-driver/mysql" // Where are working with the go-sql-driver/mysql driver for database/sql.
)

func init() {
	driver.Register("mysql", Driver{})
}

// Driver is the type that initialize new connections.
type Driver struct{}

func (d Driver) New(dsn, schema string) (driver.Conn, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return Conn{
		db:     db,
		schema: schema,
	}, nil
}

// Conn holds a connection to a MySQL database schema.
type Conn struct {
	db     *sql.DB
	schema string
}

// HasTable check if the schema has the migration table needed for Rambler to operate on it.
func (c Conn) HasTable() (bool, error) {
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

// CreateTable create the migration table using a MySQL-compatible syntax.
func (c Conn) CreateTable() error {
	_, err := c.db.Exec(`CREATE TABLE migrations ( migration VARCHAR(255) NOT NULL ) DEFAULT CHARSET=utf8`)
	return err
}

// GetApplied returns the list of already applied migrations.
func (c Conn) GetApplied() ([]string, error) {
	rows, err := c.db.Query(`SELECT migration FROM migrations ORDER BY migration ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	return migrations, nil
}

// AddApplied record that a migration was applied.
func (c Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (migration) VALUES (?)`, migration)
	return err
}

// RemoveApplied record that a migration was reversed.
func (c Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE migration = ?`, migration)
	return err
}

// Execute run a statement on the schema.
func (c Conn) Execute(statement string) error {
	_, err := c.db.Exec(statement)
	return err
}
