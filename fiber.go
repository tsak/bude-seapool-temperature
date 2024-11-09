package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"log/slog"
	"time"
)

func FiberApp(cfg *Config, sm *StateManager, monnit *Monnit) *fiber.App {
	displayGenerator := NewImageGenerator(cfg.ImageWidth, cfg.ImageHeight, GenerateDisplayImage)
	websiteImageGenerator := NewImageGenerator(300, 125, GenerateWebsiteImage)
	tinyImageGenerator := NewImageGenerator(100, 50, GenerateTinyImage)

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		AppName:               "Bude Seapool Temperature Display",
		DisableStartupMessage: true,
		Views:                 engine,
	})

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.png",
		URL:  "/favicon.ico",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-cache")
		return c.Render("index", fiber.Map{
			"width":  cfg.ImageWidth,
			"height": cfg.ImageHeight,
		})
	})

	app.Get("/temperature.png", func(c *fiber.Ctx) error {
		last := monnit.LastReading()

		// Refresh image if stale
		if displayGenerator.NeedsUpdate(time.Time(last.MessageDate)) {
			slog.Debug("Stale display image", "update", displayGenerator.lastUpdate, "current", time.Time(last.MessageDate))
			err := displayGenerator.Refresh(last)
			if err != nil {
				slog.Warn("Unable to refresh display image", "error", err)
				return c.Status(500).SendString(err.Error())
			}
		}

		sm.IncrementImageRequests()
		c.Set("Cache-Control", "no-cache")
		c.Set("Content-Type", "image/png")
		return c.Send(displayGenerator.GetImageBytes())

	})

	app.Get("/website.png", func(c *fiber.Ctx) error {
		last := monnit.LastReading()

		// Refresh image if stale
		if websiteImageGenerator.NeedsUpdate(time.Time(last.MessageDate)) {
			slog.Debug("Stale website image", "update", websiteImageGenerator.lastUpdate, "current", time.Time(last.MessageDate))
			err := websiteImageGenerator.Refresh(last)
			if err != nil {
				slog.Warn("Unable to refresh website image", "error", err)
				return c.Status(500).SendString(err.Error())
			}
		}

		sm.IncrementImageRequests()
		c.Set("Cache-Control", "no-cache")
		c.Set("Content-Type", "image/png")
		return c.Send(websiteImageGenerator.GetImageBytes())

	})

	app.Get("/tiny.png", func(c *fiber.Ctx) error {
		last := monnit.LastReading()

		// Refresh image if stale
		if tinyImageGenerator.NeedsUpdate(time.Time(last.MessageDate)) {
			slog.Debug("Stale tiny image", "update", tinyImageGenerator.lastUpdate, "current", time.Time(last.MessageDate))
			err := tinyImageGenerator.Refresh(last)
			if err != nil {
				slog.Warn("Unable to refresh tiny image", "error", err)
				return c.Status(500).SendString(err.Error())
			}
		}

		sm.IncrementImageRequests()
		c.Set("Cache-Control", "no-cache")
		c.Set("Content-Type", "image/png")
		return c.Send(tinyImageGenerator.GetImageBytes())

	})

	return app
}
