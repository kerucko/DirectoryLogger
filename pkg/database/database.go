package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_directory_logger/internal/config"
	"log"
)

var DB *sql.DB

func Init() {
	var err error
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.C.Storage.User, config.C.Storage.Password, config.C.Storage.Host, config.C.Storage.Port, config.C.Storage.Database)
	log.Println("connect to DB:", connection)
	DB, err = sql.Open("mysql", connection)
	if err != nil {
		log.Println("error open DB")
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Println("error ping DB")
		panic(err)
	}
}
