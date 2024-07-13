package scanner

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"go_directory_logger/internal/config"
	"go_directory_logger/pkg/database"

	"github.com/fsnotify/fsnotify"
)

type Scanner struct {
	Watcher       *fsnotify.Watcher
	IncludeRegexp []string
	ExcludeRegexp []string
}

func NewScanner(config config.DirConfig) (*Scanner, error) {
	if _, err := os.Stat(config.Path); err != nil {
		return nil, err
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if err := w.Add(path); err != nil {
				return err
			}
			log.Printf("dir '%s' added in watcher\n", path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Scanner{
		Watcher:       w,
		IncludeRegexp: config.IncludeRegexp,
		ExcludeRegexp: config.ExcludeRegexp,
	}, nil
}

func (s *Scanner) Scan() chan fsnotify.Event {
	log.Println("start listening")
	events := make(chan fsnotify.Event, 1)

	go func() {
		defer s.Watcher.Close()
		defer close(events)
		for {
			select {
			case event, ok := <-s.Watcher.Events:
				if !ok {
					log.Println("error watcher.Events; event")
					return
				}
				//log.Println("event:", event)
				events <- event
			case err, ok := <-s.Watcher.Errors:
				if !ok {
					log.Println("error watcher.Errors:", err)
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	return events
}

func (s *Scanner) RegexpFilter(events chan fsnotify.Event) chan fsnotify.Event {
	out := make(chan fsnotify.Event, 1)

	go func() {
		for event := range events {
			flag := false

			for _, includeReg := range s.IncludeRegexp {
				ok, err := regexp.Match(includeReg, []byte(event.Name))
				if err != nil {
					log.Println("error regexp.Match includeRegexp:", err)
					panic(err)
				}
				if ok {
					flag = true
					break
				}
			}
			if s.IncludeRegexp == nil {
				flag = true
			}
			for _, excludeReg := range s.ExcludeRegexp {
				ok, err := regexp.Match(excludeReg, []byte(event.Name))
				if err != nil {
					log.Println("error regexp.Match excludeRegexp:", err)
					panic(err)
				}
				if ok {
					flag = false
					break
				}
			}

			if flag {
				log.Println("filter access:", event.Name, event.Op)
				out <- event
			} else {
				log.Println("filter ban:", event.Name, event.Op)
			}
		}

		close(out)
	}()

	return out
}

func (s *Scanner) Log(events chan fsnotify.Event) error {
	for event := range events {
		operations := map[int]string{
			1:  "CREATE",
			2:  "WRITE",
			4:  "REMOVE",
			8:  "RENAME",
			16: "CHMOD",
		}

		stmt, err := database.DB.Prepare("INSERT INTO directory_logger.files (dirPath, filename, operation, date) VALUES (?, ?, ?, NOW())")
		if err != nil {
			return err
		}
		defer stmt.Close()

		dirPath, filename := filepath.Split(event.Name)
		res, err := stmt.Exec(dirPath, filename, operations[int(event.Op)])
		if err != nil {
			return err
		}
		log.Printf("add in DB: %s %s %s\n", dirPath, filename, operations[int(event.Op)])
		
		if r, _ := res.RowsAffected(); r != 1 {
			return errors.New("res.RowsAffected != 1")
		}
	}

	return nil
}
