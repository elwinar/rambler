package apply

import (
	"database/sql"
)

type MockMigration struct {
	scan func(string) []string
}

func (m MockMigration) Scan(section string) []string {
	return m.scan(section)
}

type MockTransaction struct {
	exec     func(query string, args ...interface{}) (sql.Result, error)
	commit   func() error
	rollback func() error
}

func (tx MockTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.exec(query, args)
}

func (tx MockTransaction) Commit() error {
	return tx.commit()
}

func (tx MockTransaction) Rollback() error {
	return tx.rollback()
}

type MockResult struct{}

func (res MockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (res MockResult) RowsAffected() (int64, error) {
	return 0, nil
}

type MockService struct {
	migrationTableExists func() (bool, error)
	createMigrationTable func() error
	listAppliedMigrations func() ([]uint64, error)
	listAvailableMigrations func() ([]uint64, error)
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
