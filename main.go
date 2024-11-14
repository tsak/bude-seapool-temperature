package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/log"
	"log/slog"
	"os"
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

	if cfg.SensorId == "" || cfg.ApiKeyId == "" || cfg.ApiSecretKey == "" || cfg.ApiUrl == "" {
		slog.Error("missing configuration", "config", cfg)
		os.Exit(1)
	}

	// Initiate sensor reader
	monnit := NewMonnit(cfg.SensorId, cfg.ApiKeyId, cfg.ApiSecretKey, cfg.ApiUrl, cfg.RefreshInterval)

	// Initiate state
	sm, err := NewStateManager(cfg.StateFile, cfg.StateAutosaveInterval)
	if err != nil {
		slog.Error("unable to load or create state", "error", err)
	}
	slog.Debug("loaded application state", "state", sm.state, "filename", sm.filename)

	// Initialise Telegram bot
	if cfg.TelegramToken != "" {
		bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
		if err != nil {
			slog.Error("unable to initialise Telegram bot", "error", err)
			os.Exit(1)
		}
		//bot.Debug = cfg.Debug
		slog.Info("became Telegram bot", "user", bot.Self.UserName)
		go TelegramBot(bot, sm, monnit)
	}
	// Set up Fiber app
	app := FiberApp(cfg, sm, monnit)

	// Start app server
	log.Fatal(app.Listen(cfg.Address))
}
