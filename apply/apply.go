package apply

import (
	"database/sql"
	"errors"
	"github.com/elwinar/rambler/lib"
)

var (
	ErrNilMigration = errors.New("nil migration")
)

func Apply(migration *lib.Migration, db *sql.DB) (error, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	
	return apply(migration, tx)
}

type migration interface {
	Scan(string) []string
}

type txer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

func apply(migration migration, tx txer) (error, error) {
	if migration == nil {
		return ErrNilMigration, nil
	}
	
	for _, statement := range migration.Scan("up") {
		_, err := tx.Exec(statement)
		if err != nil {
			rollbackErr := tx.Rollback()
			return err, rollbackErr
		}
	}
	
	commitErr := tx.Commit()
	return nil, commitErr
}
