package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// Config holds the parsed config from config.yml
type Config struct {
	Watch struct {
		Path string `yaml:"path"` // Absolute path to the file to watch
	} `yaml:"watch"`
	Destination struct {
		Path string `yaml:"path"` // Absolute path to the file to watch
	} `yaml:"destination"`
	Wait int `yaml:"wait"` // Number in seconds of how long to wait after the last write to copy the file
}

func main() {

	cfg := setup()

	log.Printf("Watching for file \"%s\", will move to \"%s\"", cfg.Watch.Path, cfg.Destination.Path)
	watch(cfg)
}

func setup() Config {

	// Open config file
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Error opening config.yml: %s", err)
	}
	defer f.Close()

	// Parse yaml
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding config.yml: %s", err)
	}

	// Normalize watch path
	abs, err := filepath.Abs(cfg.Watch.Path)
	if err != nil {
		log.Fatalf("Error getting watch absolute path: %s", err)
	}
	cfg.Watch.Path = abs

	// Normalize destination path
	abs, err = filepath.Abs(cfg.Destination.Path)
	if err != nil {
		log.Fatalf("Error getting destination absolute path: %s", err)
	}
	cfg.Destination.Path = abs

	return cfg
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
					log.Printf("Resetting move timer")
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
