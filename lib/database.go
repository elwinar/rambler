package lib

import (
	"fmt"
	"github.com/elwinar/viper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// GetDB initialize the database object used by the application
// TODO Implement compatibility with more database vendors
func GetDB() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch viper.GetString("driver") {
	default:
		fallthrough
	case "mysql":
		db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", viper.GetString("user"), viper.GetString("password"), viper.GetString("protocol"), viper.GetString("host"), viper.GetInt("port"), viper.GetString("database")))
	case "postgres":
		db, err = sqlx.Connect("postgres", fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", viper.GetString("protocol"), viper.GetString("user"), viper.GetString("password"), viper.GetString("host"), viper.GetString("database")))
	}

	return db, err
}

// HasMigrationTable check whether the given database has a migration table
// TODO Implement compatibility with more database vendors
func HasMigrationTable(db *sqlx.DB) bool {
	switch viper.GetString("driver") {
	default:
		fallthrough
	case "mysql":
		err := db.Get(new(struct {
			TableName string `db:"table_name"`
		}), fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_schema = '%s' AND table_name = 'migrations'`, viper.GetString("database")))
		return err == nil
	case "postgres":
		err := db.Get(new(struct {
			TableName string `db:"table_name"`
		}), fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_catalog = '%s' AND table_name = 'migrations'`, viper.GetString("database")))
		return err == nil
	}

	return false
}

// CreateMigrationTable create the migration table in the given database
func CreateMigrationTable(db *sqlx.DB) error {
	var err error
	switch viper.GetString("driver") {
	default:
		fallthrough
	case "mysql":
		_, err = db.Exec(`CREATE TABLE migrations ( version BIGINT UNSIGNED NOT NULL PRIMARY KEY, date DATETIME NOT NULL, description CHAR(255) NOT NULL, file CHAR(255) NOT NULL ) DEFAULT CHARSET=utf8;`)
	case "postgres":
		_, err = db.Exec(`CREATE TABLE migrations ( version NUMERIC(20) NOT NULL PRIMARY KEY, date TIMESTAMP NOT NULL, description TEXT NOT NULL, file TEXT NOT NULL );`)
	}
	return err
}
