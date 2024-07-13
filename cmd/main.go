package main

import (
	"log"
	"os"
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
	for _, dir := range config.C.Directories {
		s, err := scanner.NewScanner(dir)
		if err != nil {
			if os.IsNotExist(err) {
				log.Printf("dir '%s' not found\n", dir.Path)
			} else {
				log.Printf("NewScanner error: %v\n", err)
			}
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.Log(s.RegexpFilter(s.Scan()))
			if err != nil {
				log.Printf("error in scanner.Log: %v\n", err)
			}
		}()
	}

	wg.Wait()
}
