package main

import (
	"go_directory_logger/internal/config"
	"go_directory_logger/pkg/database"
)

func main() {
	config.ReadConfig()
	database.Init()
}
