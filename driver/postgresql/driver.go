package mysql

import (
	"database/sql"
	"fmt"
	"github.com/elwinar/rambler/driver"
	_ "github.com/lib/pq"
)

func init() {
	driver.Register("postgresql", Driver{})
}

type Driver struct{}

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

type Conn struct {
	db     *sql.DB
	schema string
}

func (c *Conn) MigrationTableExists() (bool, error) {
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

func (c *Conn) CreateMigrationTable() error {
	_, err := c.db.Exec(`CREATE TABLE migrations ( version NUMERIC(20) NOT NULL PRIMARY KEY, description VARCHAR(255) NOT NULL, applied_at TIMESTAMP NOT NULL );`)
	return err
}

func (c *Conn) ListAppliedMigrations() ([]uint64, error) {
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

func (c *Conn) SetMigrationApplied(version uint64, description string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (version, description, applied_at) VALUES ($1, $2, NOW())`, version, description)
	return err
}

func (c *Conn) UnsetMigrationApplied(version uint64) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE version = $1`, version)
	return err
}

func (c *Conn) Exec(query string) error {
	_, err := c.db.Exec(query)
	return err
}
