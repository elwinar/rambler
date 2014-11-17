package lib

import (
	"errors"
	"fmt"

	"github.com/elwinar/viper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// GetDB initialize the database object used by the application
// TODO Implement compatibility with more database vendors
func GetDB() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch Env.Driver {
	case "mysql":
		db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", Env.User, Env.Password, Env.Protocol, Env.Host, Env.Port, Env.Database))
	case "postgres":
		db, err = sqlx.Connect("postgres", fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable", viper.GetString("protocol"), viper.GetString("user"), viper.GetString("password"), viper.GetString("host"), viper.GetString("database")))
	default:
		return nil, errors.New("unsupported driver " + Env.Driver)
	}

	return db, err
}

// HasMigrationTable check whether the given database has a migration table
// TODO Implement compatibility with more database vendors
func HasMigrationTable(db *sqlx.DB) (bool, error) {
	type Table struct {
		Name string `db:'name'`
	}

	switch Env.Driver {
	case "mysql":
		err := db.Get(new(Table), fmt.Sprintf(`SELECT table_name as name FROM information_schema.tables WHERE table_schema = '%s' AND table_name = 'migrations'`, Env.Database))
		return err == nil, nil
	case "postgres":
		err := db.Get(new(struct {
			TableName string `db:"table_name"`
		}), fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_catalog = '%s' AND table_name = 'migrations'`, viper.GetString("database")))
		return err == nil
	default:
		return false, errors.New("unsupported driver " + Env.Driver)
	}
}

// CreateMigrationTable create the migration table in the given database
// TODO Implement compatibility with more database vendors
func CreateMigrationTable(db *sqlx.DB) error {
	switch Env.Driver {
	case "mysql":
		_, err := db.Exec(`CREATE TABLE migrations ( version BIGINT UNSIGNED NOT NULL PRIMARY KEY, date DATETIME NOT NULL, description CHAR(255) NOT NULL, file CHAR(255) NOT NULL ) DEFAULT CHARSET=utf8;`)
		return err
	case "postgres":
		_, err = db.Exec(`CREATE TABLE migrations ( version NUMERIC(20) NOT NULL PRIMARY KEY, date TIMESTAMP NOT NULL, description TEXT NOT NULL, file TEXT NOT NULL );`)
	default:
		return errors.New("unsupported driver " + Env.Driver)
	}
}
