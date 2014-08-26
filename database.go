package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	dsnMysql = `%s:%s@%s(%s:%d)/%s?parseTime=true`
)

func GetDatabase() (*sqlx.DB, error) {
	// Connect to the database
	// TODO add connection options depending on the SQL database flavour
	db, err := sqlx.Connect("mysql", fmt.Sprintf(dsnMysql, config.User, config.Password, config.Protocol, config.Host, config.Port, config.Database))
	if err != nil {
		return nil, err
	}

	return db, nil
}
