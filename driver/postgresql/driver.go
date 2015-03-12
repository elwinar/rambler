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

func (c *Conn) CreateTable() error {
	_, err := c.db.Exec(`CREATE TABLE migrations ( migration VARCHAR(255) NOT NULL );`)
	return err
}

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

func (c *Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(`INSERT INTO migrations (migration) VALUES ($1)`, migration)
	return err
}

func (c *Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(`DELETE FROM migrations WHERE migration = $1`, migration)
	return err
}

func (c *Conn) Execute(query string) error {
	_, err := c.db.Exec(query)
	return err
}
