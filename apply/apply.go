package apply

import (
	"database/sql"
	"errors"
)

type scanner interface {
	Scan(string) []string
}

type txer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
}

// Apply tries to apply statements of the given migration on the given transaction.
// It will try to roolback in case of error, and will return 2 distinct errors:
// - the SQL error which caused the transaction to fail
// - the SQL error which caused the rollback/commit to fail
func Apply(m scanner, tx txer) (error, error) {
	if m == nil {
		return errors.New("nil migration"), nil
	}

	if tx == nil {
		return errors.New("nil transaction"), nil
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
