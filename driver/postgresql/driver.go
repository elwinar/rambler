package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/elwinar/rambler/log"

	"github.com/elwinar/rambler/driver"
	_ "github.com/lib/pq" // Working with the lib/pq PostgreSQL driver here.
)

var logger  *log.Logger

func init() {
	driver.Register("postgresql", Driver{})
	logger = log.NewLogger(func(l *log.Logger) {
		l.PrintDebug = true
	})
}



// Driver initialize new connections to a PostgreSQL database schema.
type Driver struct{}

// New returns a new connection.
func (d Driver) New(dsn, database, schema, table string) (driver.Conn, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Use the default schema.
	if schema == "" {
		schema = "public"
	}
	c := &Conn{
		db:     db,
		database: database,
		schema: schema,
		table:  table,
	}

	return c, nil
}

func (c*  Conn) qualifiedTable() string {
	if c.schema == "" {
		return c.table
	}
	return c.schema + "." + c.table
}

// Connection holds a database connection.
type Conn struct {
	db     *sql.DB
	database string
	schema string
	table  string
}

// HasTable check if the schema has the migration table.
func (c *Conn) HasTable() (bool, error) {
	logger.Debug(fmt.Sprintf("Loading against: %s", c))

	var name string

	err := c.db.QueryRow(
		`SELECT table_name FROM information_schema.tables WHERE table_catalog = $1 AND table_schema = $2 AND table_name = $3`,
		c.database, c.schema, c.table).Scan(&name)
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
	_, err := c.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( migration VARCHAR(255) NOT NULL );`, c.qualifiedTable()))
	if err != nil {
		err = fmt.Errorf("failed to create table '%s'. %s", c.table, err)
	}
	return err
}

// GetApplied returns the list of applied migrations.
func (c *Conn) GetApplied() ([]string, error) {

	rows, err := c.db.Query(fmt.Sprintf(`SELECT migration FROM %s ORDER BY migration ASC`, c.qualifiedTable()))
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
	_, err := c.db.Exec(fmt.Sprintf(`INSERT INTO %s (migration) VALUES ($1)`, c.qualifiedTable()), migration)
	return err
}

// RemoveApplied records a migration as reversed.
func (c *Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE migration = $1`, c.qualifiedTable()), migration)
	return err
}

// Execute run a statement on the schema.
func (c *Conn) Execute(query string) error {
	_, err := c.db.Exec(query)
	return err
}

func (c *Conn) String() string {
	return fmt.Sprintf("database: '%s', schema: '%s', table: '%s'\n", c.database, c.schema, c.table)
}