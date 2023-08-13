package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_directory_logger/internal/config"
)

var DB *sql.DB

func Init() {
	var err error
	//log.Println(c)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.C.Storage.User, config.C.Storage.Password, config.C.Storage.Host, config.C.Storage.Port, config.C.Storage.Database)
	//log.Println(connection)
	DB, err = sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}
	//defer DB.Close()

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
