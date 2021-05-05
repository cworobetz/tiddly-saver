package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

func main() {

	setupLogging()
	cfg := getConfig()
	logrus.Printf("Watching for file %s, will move to %s", cfg.Watch.Path, cfg.Destination.Path)
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
		logrus.Fatal(err)
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
					if timer.Reset(duration) {
						logrus.Printf("Change to \"%s\" detected, restarting wait period of %d seconds", cfg.Watch.Path, cfg.Wait)
					} else {
						logrus.Printf("Change to \"%s\" detected, starting wait period of %d seconds", cfg.Watch.Path, cfg.Wait)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Println("error:", err)
			}
		}
	}()

	// Goroutine to respond to timer events
	go func() {
		for {
			<-timer.C
			timer.Reset(duration)
			timer.Stop()
			logrus.Printf("%d second wait period has passed, moving \"%s\" to \"%s\"", cfg.Wait, cfg.Watch.Path, cfg.Destination.Path)
			err := os.Rename(cfg.Watch.Path, cfg.Destination.Path)
			if err != nil {
				logrus.Fatalf("Error moving file: %s", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(cfg.Watch.Path))
	if err != nil {
		logrus.Fatal(err)
	}
	<-done
}
