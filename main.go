package main

import (
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/support"
)

var log = logger.GetLogger(info.ID)

// these values come from the command line parameters when run (see README) or from the Sphere system
var host = config.String("localhost", "led.host")
var port = config.Int(3115, "led.remote.port")

func main() {

	app := &App{}
	err := app.Init(info)
	if err != nil {
		app.Log.Fatalf("failed to initialize app: %v", err)
	}

	err = app.Export(app)
	if err != nil {
		app.Log.Fatalf("failed to export app: %v", err)
	}

	support.WaitUntilSignal()
}
