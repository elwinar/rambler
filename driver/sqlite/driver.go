package mysql

import (
	"database/sql"
	"fmt"

	"github.com/elwinar/rambler/driver"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	driver.Register("sqlite", Driver{})
}

type Driver struct{}

func (d Driver) New(config driver.Config) (driver.Conn, error) {
	db, err := sql.Open("sqlite3", config.Database)
	if err != nil {
		return nil, err
	}

	return Conn{
		db:    db,
		table: config.Table,
	}, nil
}

type Conn struct {
	db    *sql.DB
	table string
}

func (c Conn) HasTable() (bool, error) {
	var table string
	err := c.db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, c.table).Scan(&table)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func (c Conn) CreateTable() error {
	_, err := c.db.Exec(fmt.Sprintf(`CREATE TABLE %s ( migration VARCHAR(255) NOT NULL );`, c.table))
	return err
}

func (c Conn) GetApplied() ([]string, error) {
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

func (c Conn) AddApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`INSERT INTO %s (migration) VALUES (?)`, c.table), migration)
	return err
}

func (c Conn) RemoveApplied(migration string) error {
	_, err := c.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE migration = ?`, c.table), migration)
	return err
}

func (c Conn) Execute(statement string) error {
	_, err := c.db.Exec(statement)
	return err
}
