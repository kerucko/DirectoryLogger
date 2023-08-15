package main

import (
	"go_directory_logger/internal/config"
	"go_directory_logger/internal/scanner"
	"log"
	"sync"
)

func main() {
	config.ReadConfig()
	//database.Init()
	//defer database.DB.Close()

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
			s.Log(s.RegexpFilter(s.Scan()))
		}()
	}

	wg.Wait()
}
