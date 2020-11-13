package mysql

import (
	"database/sql"
	"fmt"

	"github.com/elwinar/rambler/driver"
	_ "github.com/lib/pq" // Working with the lib/pq PostgreSQL driver here.
)

func init() {
	driver.Register("postgresql", Driver{})
}

// Driver initialize new connections to a PostgreSQL database schema.
type Driver struct{}

// New returns a new connection.
func (d Driver) New(config driver.Config) (driver.Conn, error) {
	if config.Schema == "" {
		config.Schema = "public"
	}

	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}

	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.Schema, config.SSLMode,
	))
	if err != nil {
		return nil, err
	}

	c := &Conn{
		db:      db,
		catalog: config.Database,
		schema:  config.Schema,
		table:   config.Table,
	}

	return c, nil
}

// Connection holds a database connection.
type Conn struct {
	db      *sql.DB
	catalog string
	schema  string
	table   string
}

// HasTable check if the schema has the migration table.
func (c *Conn) HasTable() (bool, error) {
	var name string
	err := c.db.QueryRow(`SELECT table_name FROM information_schema.tables WHERE table_catalog = $1 AND table_schema = $2 AND table_name = $3`, c.catalog, c.schema, c.table).Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err != nil {
		return false, nil
	}
	return true, nil
}

// CreateTable create the migration table using a PostgreSQL-compatible syntax.
func (c *Conn) CreateTable() error {
	_, err := c.db.Exec(fmt.Sprintf(`CREATE TABLE %s ( migration VARCHAR(255) NOT NULL );`, c.table))
	return err
}

// GetApplied returns the list of applied migrations.
func (c *Conn) GetApplied() ([]string, error) {
	rows, err := c.db.Query(fmt.Sprintf(`SELECT migration FROM %s ORDER BY migration ASC`, c.table))
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

// AddApplied records a migration as applied.
func (c *Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`INSERT INTO %s (migration) VALUES ($1)`, c.table), migration)
	return err
}

// RemoveApplied records a migration as reversed.
func (c *Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE migration = $1`, c.table), migration)
	return err
}

// Execute run a statement on the schema.
func (c *Conn) Execute(query string) error {
	_, err := c.db.Exec(query)
	return err
}
