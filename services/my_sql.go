package services

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var mysqlInstance *sqlx.DB

func MySQL() *sqlx.DB {
	if mysqlInstance == nil {
		connectionURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			"root",
			"12345678",
			"127.0.0.1",
			"3306",
			"trabalho_final",
		)

		db, err := sqlx.Connect("mysql", connectionURL)
		if err != nil {
			panic("failed to connect to MySQL: " + err.Error())
		}

		mysqlInstance = db
	}

	return mysqlInstance
}
