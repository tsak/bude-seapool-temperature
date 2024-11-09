package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"log/slog"
	"time"
)

func FiberApp(cfg *Config, sm *StateManager, monnit *Monnit) *fiber.App {
	generators := make(map[string]*ImageGenerator)
	generators["temperature"] = NewImageGenerator(cfg.ImageWidth, cfg.ImageHeight, GenerateDisplayImage)
	generators["website"] = NewImageGenerator(300, 125, GenerateWebsiteImage)
	generators["tiny"] = NewImageGenerator(100, 50, GenerateTinyImage)

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

	app.Get(`/:type<regex((temperature|website|tiny))>.png`, func(c *fiber.Ctx) error {
		imageType := c.Params("type")
		generator := generators[imageType]
		last := monnit.LastReading()

		// Refresh image if stale
		if generator.NeedsUpdate(time.Time(last.MessageDate)) {
			slog.Debug("Stale image", "update", generator.lastUpdate, "current", time.Time(last.MessageDate), "image_type", imageType)
			err := generator.Refresh(last)
			if err != nil {
				slog.Warn("Unable to refresh image", "error", err, "image_type", imageType)
				return c.Status(500).SendString(err.Error())
			}
		}

		sm.IncrementImageRequests()
		c.Set("Cache-Control", "no-cache")
		c.Set("Content-Type", "image/png")
		return c.Send(generator.GetImageBytes())
	})

	return app
}
