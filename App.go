package main

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
	// TODO: put images in app, get pane to be able to see app's data
	loadImages()
}

type RuntimeConfig struct {
}

type App struct {
	support.AppSupport
	led *remote.Matrix
}

func (a *App) Start(cfg *RuntimeConfig) error {
	log.Infof("Making new pane...")
	// The pane must implement the remote.pane interface
	pane := NewDemoPane(a.Conn)

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
