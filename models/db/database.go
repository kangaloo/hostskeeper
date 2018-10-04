package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

var DB *sqlx.DB

func init() {

	var err error

	DB, err = sqlx.Open(
		"mysql",
		"hostkeeper:123456@tcp(127.0.0.1:3306)/dbtest1?charset=utf8&parseTime=True&loc=Local",
	)

	if err != nil {
		log.Fatalln(err)
	}
}
