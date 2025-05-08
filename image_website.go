package main

import (
	"github.com/fogleman/gg"
	"image"
	"log/slog"
)

// GenerateWebsiteImage generates a smaller image for websites.
func GenerateWebsiteImage(width, height int, temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	// Temperature display
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 100); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored(temperature, float64(width)/2, float64(height)/2, 0.5, 0.25)

	// Last updated
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-LightItalic.ttf", 15); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}

	dc.DrawStringAnchored("Last updated "+lastModified, float64(width)/2, float64(height)-20, 0.5, 0.5)

	return dc.Image(), nil
}

// GenerateMaintenanceWebsiteImage generates a smaller image with an "Annual maintenance" message.
func GenerateMaintenanceWebsiteImage(width, height int, temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	// Temperature display
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 30); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored("Annual maintenance", float64(width)/2, float64(height)/2, 0.5, 0.5)

	return dc.Image(), nil
}
