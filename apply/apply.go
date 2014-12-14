package apply

import (
	"fmt"
	"github.com/elwinar/rambler/migration/driver"
)

// Apply tries to apply statements on the given transaction.
// It will try to roolback in case of error, and will return 2 distinct errors:
// - the SQL error which caused the transaction to fail
// - the SQL error which caused the rollback/commit to fail
func Apply(statements []string, tx driver.Tx) (error, error) {
	if tx == nil {
		return nil, fmt.Errorf(errNilTransaction)
	}

	for _, statement := range statements {
		_, sqlerr := tx.Exec(statement)
		if sqlerr != nil {
			txerr := tx.Rollback()
			return sqlerr, txerr
		}
	}

	txerr := tx.Commit()
	return nil, txerr
}
