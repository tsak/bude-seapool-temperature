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

func GenerateTinyImage(width, height int, temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(width, height)

	// White background
	dc.SetRGBA(1, 1, 1, 0)
	dc.Clear()

	dc.DrawImage(thermometerImage, 8, 10)

	// Temperature display
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 20); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}

	dc.DrawStringAnchored(temperature, float64(width)/2+8, float64(height)/2, 0.5, 0.5)

	return dc.Image(), nil
}
