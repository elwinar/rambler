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
func (d Driver) New(dsn, schema string) (driver.Conn, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		db:     db,
		schema: schema,
	}

	return c, nil
}

// Connection holds a database connection.
type Conn struct {
	db     *sql.DB
	schema string
}

// HasTable check if the schema has the migration table.
func (c *Conn) HasTable() (bool, error) {
	var name string
	err := c.db.QueryRow(fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_catalog = '%s' AND table_name = 'migrations'`, c.schema)).Scan(&name)
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
	_, err := c.db.Exec(`CREATE TABLE migrations ( migration VARCHAR(255) NOT NULL );`)
	return err
}

// GetApplied returns the list of applied migrations.
func (c *Conn) GetApplied() ([]string, error) {
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

// AddApplied records a migration as applied.
func (c *Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (migration) VALUES ($1)`, migration)
	return err
}

// RemoveApplied records a migration as reversed.
func (c *Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE migration = $1`, migration)
	return err
}

// Execute run a statement on the schema.
func (c *Conn) Execute(query string) error {
	_, err := c.db.Exec(query)
	return err
}
