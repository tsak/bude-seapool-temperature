package main

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/template/html/v2"
	"image/png"
)

func FiberApp(cfg *Config, sm *StateManager, monnit *Monnit) *fiber.App {
	display := NewDisplay(cfg.ImageWidth, cfg.ImageHeight)

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
		img, err := display.Image(last.Temperature.String(), last.MessageDate.String())
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		var b []byte
		buf := bytes.NewBuffer(b)
		if err := png.Encode(buf, img); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		sm.IncrementImageRequests()
		c.Set("Cache-Control", "no-cache")
		c.Set("Content-Type", "image/png")
		return c.Send(buf.Bytes())

	})

	return app
}
