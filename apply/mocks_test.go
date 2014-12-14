package apply

import (
	"database/sql"
)

type MockTx struct {
	exec     func(string, ...interface{}) (sql.Result, error)
	commit   func() error
	rollback func() error
}

func (tx MockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.exec(query, args...)
}

func (tx MockTx) Commit() error {
	return tx.commit()
}

func (tx MockTx) Rollback() error {
	return tx.rollback()
}

type MockResult struct{}

func (res MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (res MockResult) RowsAffected() (int64, error) {
	return 0, nil
}

/*
import (
	"database/sql"
	"github.com/elwinar/rambler/migration"
)

type MockMigration struct {
	scan func(string) []string
}

func (m MockMigration) Scan(section string) []string {
	return m.scan(section)
}

type MockService struct {
	migrationTableExists func() (bool, error)
	createMigrationTable func() error
	listAppliedMigrations func() ([]uint64, error)
	listAvailableMigrations func() ([]uint64, error)
	startTransaction func() (migration.Transaction, error)
}

func (s MockService) MigrationTableExists() (bool, error) {
	return s.migrationTableExists()
}

func (s MockService) CreateMigrationTable() error {
	return s.createMigrationTable()
}

func (s MockService) ListAppliedMigrations() ([]uint64, error) {
	return s.listAppliedMigrations()
}

func (s MockService) ListAvailableMigrations() ([]uint64, error) {
	return s.listAvailableMigrations()
}

func (s MockService) StartTransaction() (migration.Transaction, error) {
	return s.startTransaction()
}
*/
