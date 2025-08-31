package main

import (
	"github.com/fogleman/gg"
	"image"
	"log/slog"
	"strings"
)

// GenerateDisplayImage generates the large display image that is displayed at the seapool
func GenerateDisplayImage(width, height int, temperature, lastModified, msg string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0.3, 0.3, 0.3)
	if err := dc.LoadFontFace("fonts/Roboto-Bold.ttf", 400); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored("POOL TEMP", float64(width)/2, 200, 0.5, 0.5)

	// Temperature display
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/Roboto-Bold.ttf", 800); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored(temperature, float64(width)/2, float64(height)/2+100, 0.5, 0.5)

	// Last updated
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-LightItalic.ttf", 100); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored("Last updated "+lastModified, float64(width)/2, float64(height)-100, 0.5, 0.5)

	return dc.Image(), nil
}

// GenerateMaintenanceDisplayImage generates a large image with the "Annual maintenance" message
func GenerateMaintenanceDisplayImage(width, height int, temperature, lastModified, msg string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Temperature display
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 400); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	parts := strings.Split(msg, "\n")
	switch len(parts) {
	case 2:
		dc.DrawStringAnchored(parts[0], float64(width)/2, float64(height)/2, 0.5, -0.25)
		dc.DrawStringAnchored(parts[1], float64(width)/2, float64(height)/2, 0.5, 1.25)
	case 1:
		dc.DrawStringAnchored(parts[0], float64(width)/2, float64(height)/2, 0.5, 0.5)
	}

	return dc.Image(), nil
}
