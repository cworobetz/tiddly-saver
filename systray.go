package main

import (
	"io/ioutil"
	"os"

	"github.com/getlantern/systray"
	"github.com/sirupsen/logrus"
)

func getIcon(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatalf("Error opening icon file: %s", err)
	}
	return b
}

func onReady() {
	// Main systray icon
	systray.SetIcon(getIcon("assets/pencil.ico"))
	systray.SetTitle("Tiddly Saver")
	systray.SetTooltip("Tiddly Saver - Download your Tiddlywiki somewhere else")

	// Tooltips
	setSystrayMenuItem("Exit", "Shutdown the app")
}

func setSystrayMenuItem(title string, tooltip string) {
	mExit := systray.AddMenuItem(title, tooltip)

	go func() {
		<-mExit.ClickedCh
		onExit()
	}()
}

func onExit() {
	logrus.Printf("Exit signal received, stopping program.")
	os.Exit(0)
}
