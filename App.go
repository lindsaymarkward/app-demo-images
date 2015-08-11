package main

// A DEMO of displaying images on the Spheramid LED matrix, and tap gestures

import (
	"fmt"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
	"github.com/ninjasphere/sphere-go-led-controller/remote"
)

var info = ninja.LoadModuleInfo("./package.json")

// init is a Go standard that runs first
func init() {
	loadImages()
}

type RuntimeConfig struct {
}

type App struct {
	support.AppSupport
	led      *remote.Matrix
	textData string
}

// Start the app (called by the system)
func (a *App) Start(cfg *RuntimeConfig) error {
	// save some text to display (to show we can access app data from LED pane)
	a.textData = "OK!"
	log.Infof("Making new pane...")
	// The pane must implement the remote.pane interface
	pane := NewDemoPane(a)

	// Export our pane over TCP
	a.led = remote.NewTCPMatrix(pane, fmt.Sprintf("%s:%d", host, port))

	// TODO: try a second NewMatrix - see if we can swipe between them

	return nil
}

// Stop the app.
func (a *App) Stop() error {
	a.led.Close()
	a.led = nil
	return nil
}
