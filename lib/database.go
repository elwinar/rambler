package lib

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"log"
)

const Schema = `CREATE TABLE migrations (
	version BIGINT UNSIGNED NOT NULL PRIMARY KEY,
	date DATETIME NOT NULL DEFAULT current_timestamp,
	description CHAR(255) NOT NULL,
	file CHAR(255) NOT NULL
) DEFAULT CHARSET=utf8;`

var (
	db *sqlx.DB
)

func Connect () error {
	var err error
	var dsn string
	
	switch viper.GetString("driver") {
		case "mysql":
			dsn = fmt.Sprintf(`%s:%s@%s(%s:%d)/%s?parseTime=true`, viper.GetString("user"), viper.GetString("password"), viper.GetString("protocol"), viper.GetString("host"), viper.GetInt("port"), viper.GetString("database"))
	}
	
	log.Println("Connecting to", viper.GetString("driver"), ":", dsn)
	db, err = sqlx.Connect(viper.GetString("driver"), dsn)
	
	return err
}

func Connected () bool {
	return db != nil
}

func HasMigrationTable () bool {
	log.Println("Looking for migration table")
	err := db.Get(new(struct{ TableName string `db:"table_name"` }), fmt.Sprintf(`SELECT table_name
FROM information_schema.tables 
WHERE table_schema = '%s' 
AND table_name = 'migrations'`, viper.GetString("database")))
	return err == nil
}

func CreateMigrationTable () error {
	log.Println("Creating migration table")
	_, err := db.Exec(Schema)
	return err
}
