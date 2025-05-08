package main

import (
	"bytes"
	_ "embed"
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"log/slog"
)

//go:embed thermometer.png
var thermometerPng []byte
var thermometerImage image.Image

func init() {
	var err error
	thermometerImage, err = png.Decode(bytes.NewReader(thermometerPng))
	if err != nil {
		panic(err)
	}
}

// GenerateTinyImage generates a small image that is displayed on https://www.budeseapool.org/
func GenerateTinyImage(width, height int, temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// Transparent background
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	dc.DrawImage(thermometerImage, 12, 10)

	// Temperature display
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("fonts/Roboto-Medium.ttf", 16); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}

	dc.DrawStringAnchored(temperature, float64(width)/2+8, float64(height)/2, 0.5, 0.5)

	return dc.Image(), nil
}

// GenerateMaintenanceTinyImage generates a transparent image during the annual cleanup
func GenerateMaintenanceTinyImage(width, height int, temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// Transparent background
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	return dc.Image(), nil
}
