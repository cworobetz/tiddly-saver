package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
)

func main() {
	cfg := getConfig()
	log.Printf("Watching for file \"%s\", will move to \"%s\"", cfg.Watch.Path, cfg.Destination.Path)
	go watch(cfg)
	systray.Run(onReady, onExit)
}

// watch takes the full path of a file and attempts to watch for that file
func watch(cfg Config) {

	duration := time.Duration(cfg.Wait) * time.Second
	timer := time.NewTimer(duration)
	timer.Stop() // We want to initialize the timer so we can reset it later. Stop it here

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	// Goroutine to watch for new files
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// If it's a write event and the name of the events file matches the file to watch
				if event.Op&fsnotify.Write == fsnotify.Write && event.Name == cfg.Watch.Path {
					// TODO - change message if starting or restarting timer
					log.Printf("Change to \"%s\" detected, starting wait period of %d seconds", cfg.Watch.Path, cfg.Wait)
					timer.Reset(duration)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Goroutine to respond to timer events
	go func() {
		for {
			<-timer.C
			timer.Reset(duration)
			timer.Stop()
			log.Printf("%d second wait period has passed, moving \"%s\" to \"%s\"", cfg.Wait, cfg.Watch.Path, cfg.Destination.Path)
			err := os.Rename(cfg.Watch.Path, cfg.Destination.Path)
			if err != nil {
				log.Fatalf("Error moving file: %s", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(cfg.Watch.Path))
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
