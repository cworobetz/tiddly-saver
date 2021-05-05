package main

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
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

func getConfig() Config {

	// Open config file
	f, err := os.Open("config.yml")
	if err != nil {
		logrus.Fatalf("Error opening config.yml: %s", err)
	}
	defer f.Close()

	// Parse yaml
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		logrus.Fatalf("Error decoding config.yml: %s", err)
	}

	// Normalize watch path
	abs, err := filepath.Abs(cfg.Watch.Path)
	if err != nil {
		logrus.Fatalf("Error getting watch absolute path: %s", err)
	}
	cfg.Watch.Path = abs

	// Normalize destination path
	abs, err = filepath.Abs(cfg.Destination.Path)
	if err != nil {
		logrus.Fatalf("Error getting destination absolute path: %s", err)
	}
	cfg.Destination.Path = abs

	logrus.Print("Configuration loaded successfully")
	return cfg
}

func setupLogging() {
	f, err := os.OpenFile("tiddly-saver.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatalf("Error opening log file: %s", err)
	}
	logrus.SetOutput(f)
}
