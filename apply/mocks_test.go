package apply

import (
	"database/sql"
)

type MockMigration struct {
	statements map[string][]string
}

func (m MockMigration) Scan(section string) []string {
	return m.statements[section]
}

type MockTxer struct {
	exec     func(query string, args ...interface{}) (sql.Result, error)
	commit   func() error
	rollback func() error
}

func (tx MockTxer) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.exec(query, args)
}

func (tx MockTxer) Commit() error {
	return tx.commit()
}

func (tx MockTxer) Rollback() error {
	return tx.rollback()
}

type MockResult struct{}

func (res MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (res MockResult) RowsAffected() (int64, error) {
	return 0, nil
}
