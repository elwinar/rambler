package lib

import (
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

// Migrations is a wrapper for a migration slice
type Migrations []Migration

// Len return the number of migrations in the slice
func (m Migrations) Len () int {
	return len(m)
}

// Less return if the migration at index i comes before the one at
// index j
func (m Migrations) Less (i, j int) bool {
	return m[i].Version < m[j].Version
}

// Swap change the position of two migrations in the slice
func (m Migrations) Swap (i, j int) {
	m[i], m[j] = m[j], m[i]
}

// GetAppliedMigrations return all migrations listed in the migration
// table
func GetAppliedMigrations (db *sqlx.DB) (Migrations, error) {
	var migrations Migrations
	err := db.Select(&migrations, `SELECT * FROM migrations`)
	return migrations, err
}

func GetAvailableMigrations () (Migrations, error) {
	var migrations Migrations
	
	files, err := filepath.Glob(filepath.Join(GetMigrationsDir(), "*.sql"))
	if err != nil {
		return migrations, err
	}
	
	for _, file := range files {
		migration, err := NewMigration(filepath.Base(file))
		if err != nil {
			continue
		}
		migrations = append(migrations, migration)
	}
	
	return migrations, nil
}
