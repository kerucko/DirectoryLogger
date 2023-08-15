package scanner

import (
	"github.com/fsnotify/fsnotify"
	"go_directory_logger/internal/config"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Scanner struct {
	Watcher       *fsnotify.Watcher
	IncludeRegexp []string
	ExcludeRegexp []string
}

func NewScanner(config config.DirConfig) (*Scanner, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(config.Path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Printf("dir '%s' added in watcher\n", path)
			if err := w.Add(path); err != nil {
				return err
			}
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
		for {
			select {
			case event, ok := <-s.Watcher.Events:
				if !ok {
					log.Println("error watcher.Events; event")
					close(events)
					return
				}
				log.Println("event:", event)
				events <- event
				//if event.Has(fsnotify.Write) {
				//	log.Println("modified file:", event.Name)
				//}
			case err, ok := <-s.Watcher.Errors:
				if !ok {
					log.Println("error watcher.Errors:", err)
					close(events)
					return
				}
				log.Println("error:", err)
			}
		}
		close(events)
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
					log.Printf("path: '%s' is %s\n", event.Name, includeReg)
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
					log.Printf("path: '%s' isnt %s\n", event.Name, excludeReg)
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

func (s *Scanner) Log(events chan fsnotify.Event) {
	for event := range events {
		log.Println("LOG:", event)
	}
}
