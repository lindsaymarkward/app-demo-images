package main

// A DEMO of displaying images on the Spheramid LED matrix, and tap gestures

import (
	"fmt"
	"net"
	"os"

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

	// Connect to the LED controller remote pane interface via TCP
	log.Infof("Connecting to LED controller...")
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	// This creates a TCP connection, conn
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("DialTCP failed:", err.Error())
		os.Exit(1)
	}

	log.Infof("Connected. Now making new matrix...")

	// Export our pane over the TCP connection we just made
	a.led = remote.NewMatrix(pane, conn)

	// TODO: try a second NewMatrix - see if we can swipe between them

	return nil
}

// Stop the app.
func (a *App) Stop() error {
	a.led.Close()
	a.led = nil
	return nil
}
