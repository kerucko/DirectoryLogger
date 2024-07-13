package main

import (
	"log"
	"sync"
	
	"go_directory_logger/internal/config"
	"go_directory_logger/internal/scanner"
	"go_directory_logger/pkg/database"
)

func main() {
	config.ReadConfig()
	database.Init()
	defer database.DB.Close()

	var wg sync.WaitGroup
	wg.Add(len(config.C.Directories))
	for _, dir := range config.C.Directories {
		s, err := scanner.NewScanner(dir)
		if err != nil {
			log.Printf("NewScanner error")
			panic(err)
		}
		go func() {
			defer wg.Done()
			err := s.Log(s.RegexpFilter(s.Scan()))
			if err != nil {
				log.Println("error in scanner.Log")
				panic(err)
			}
		}()
	}

	wg.Wait()
}
