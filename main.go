package main

import (
	"github.com/gofiber/fiber/v2/log"
	"log/slog"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	// Load configuration from environment
	cfg := LoadConfig()
	if cfg.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("config", "config", cfg)
	}

	// Initiate sensor reader
	monnit := NewMonnit(cfg.SensorId, cfg.ApiKeyId, cfg.ApiSecretKey, cfg.ApiUrl, cfg.RefreshInterval)

	// Initiate state
	sm, err := NewStateManager(cfg.StateFile, cfg.StateAutosaveInterval)
	if err != nil {
		slog.Error("unable to load or create state", "error", err)
	}
	slog.Debug("loaded application state", "state", sm.state, "filename", sm.filename)

	// Set up Fiber app
	app := FiberApp(cfg, sm, monnit)

	// Start app server
	log.Fatal(app.Listen(cfg.Address))
}
