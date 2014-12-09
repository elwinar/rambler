package apply

import (
	"database/sql"
	"errors"
)

var (
	ErrNilMigration   = errors.New("nil migration")
	ErrNilTransaction = errors.New("nil transaction")
)

type scanner interface {
	Scan(string) []string
}

type txer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

func Apply(m scanner, tx txer) (error, error) {
	if m == nil {
		return ErrNilMigration, nil
	}

	if tx == nil {
		return ErrNilTransaction, nil
	}

	for _, statement := range m.Scan("up") {
		_, sqlerr := tx.Exec(statement)
		if sqlerr != nil {
			err := tx.Rollback()
			return err, sqlerr
		}
	}

	err := tx.Commit()
	return err, nil
}
