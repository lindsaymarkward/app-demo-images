package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/support"
	"github.com/ninjasphere/go-uber"
	"github.com/ninjasphere/sphere-go-led-controller/remote"
)

var info = ninja.LoadModuleInfo("./package.json")

var client *uber.Client

type UberConfig struct {
	ClientID    string `json:"clientId"`
	ServerToken string `json:"serverToken"`
	Secret      string `json:"secret"`
}

var states []string

// init is a Go standard that runs first
func init() {

	loadImages()
}

type RuntimeConfig struct {
}

type App struct {
	support.AppSupport
	led *remote.Matrix
}

func (a *App) Start(cfg *RuntimeConfig) error {
	log.Infof("making new pane...")
	pane := NewDemoPane(a.Conn)

	// Connect to the led controller remote pane interface
	log.Infof("Connecting to led controller")
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	log.Infof("Connected. Now making new matrix")

	// Export our pane over this interface
	a.led = remote.NewMatrix(pane, conn)

	return nil
}

// Stop the app.
func (a *App) Stop() error {
	a.led.Close()
	a.led = nil
	return nil
}
