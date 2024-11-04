package main

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	// Monnit sensor ID
	SensorId string `env:"MONNIT_SENSOR_ID"`

	// Monnit API key
	ApiKeyId string `env:"MONNIT_API_KEY_ID"`

	// Monnit API secret key
	ApiSecretKey string `env:"MONNIT_API_SECRET_KEY"`

	// Monnit API URL
	ApiUrl string `env:"MONNIT_API_URL"`

	// Monnit refresh interval
	RefreshInterval time.Duration `env:"MONNIT_REFRESH_INTERVAL" envDefault:"10m"`

	// Image width
	ImageWidth int `env:"IMAGE_WIDTH" envDefault:"2560"`

	// Image height
	ImageHeight int `env:"IMAGE_HEIGHT" envDefault:"1440"`

	// Address the webserver will listen on
	Address string `env:"ADDRESS"`

	// State file
	StateFile string `env:"STATE_FILE" envDefault:"state.gob"`

	// State autosave interval
	StateAutosaveInterval time.Duration `env:"STATE_AUTOSAVE_INTERVAL" envDefault:"1m"`

	// Debug mode
	Debug bool `env:"DEBUG" envDefault:"false"`
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("sensor_id", c.SensorId),
		slog.String("api_key_id", c.ApiKeyId),
		slog.String("api_url", c.ApiUrl),
		slog.Duration("refresh_interval", c.RefreshInterval),
		slog.Int("image_width", c.ImageWidth),
		slog.Int("image_height", c.ImageHeight),
		slog.String("address", c.Address),
		slog.String("state_file", c.StateFile),
		slog.Duration("state_autosave_interval", c.StateAutosaveInterval),
	)
}

func LoadConfig() *Config {
	// Check if .env file exists and Load into environment if so
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			slog.Warn("error loading .env file", "error", err)
		}
	}

	// Load config from environment
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		slog.Error("unable to parse config", "error", err)
		os.Exit(1)
	}
	return &cfg
}
