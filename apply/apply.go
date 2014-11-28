package apply

import (
	"database/sql"
	"errors"
	"github.com/elwinar/rambler/lib"
)

var (
	ErrNilMigration   = errors.New("nil migration")
	ErrNilTransaction = errors.New("nil transaction")
)

func Apply(migration *lib.Migration, tx *sql.Tx) (error, error) {
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

func apply(migration migration, tx txer) (err error, sqlerr error) {
	if migration == nil {
		return ErrNilMigration, nil
	}

	if tx == nil {
		return ErrNilTransaction, nil
	}

	for _, statement := range migration.Scan("up") {
		_, sqlerr := tx.Exec(statement)
		if sqlerr != nil {
			err := tx.Rollback()
			return err, sqlerr
		}
	}

	err = tx.Commit()
	return err, nil
}
