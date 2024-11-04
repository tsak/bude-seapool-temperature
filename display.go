package main

import (
	"github.com/fogleman/gg"
	"image"
	"log/slog"
)

type Display struct {
	width  int
	height int
}

func NewDisplay(width, height int) *Display {
	return &Display{width, height}
}

func (d *Display) Image(temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(d.width, d.height)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Temperature display
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 500); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored(temperature, float64(d.width)/2, float64(d.height)/2, 0.5, 0.5)

	// Last updated
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-LightItalic.ttf", 50); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored("Last updated "+lastModified, float64(d.width)/2, float64(d.height)-50, 0.5, 0.5)

	return dc.Image(), nil
}
