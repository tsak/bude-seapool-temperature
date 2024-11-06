package main

import (
	"bytes"
	"github.com/fogleman/gg"
	"image"
	"image/png"
	"log/slog"
	"sync"
	"time"
)

// Display represents a display with width, height, buffer, and last update time.
type Display struct {
	sync.RWMutex
	width      int
	height     int
	buffer     *bytes.Buffer
	lastUpdate time.Time
}

// NewDisplay creates a new display with the specified width and height, initializing the display buffer.
func NewDisplay(width, height int) *Display {
	return &Display{
		width:  width,
		height: height,
		buffer: bytes.NewBuffer([]byte{}),
	}
}

// GetImageBytes returns the image data as a byte slice by reading from the display buffer.
// It will be blocked while a call to [Display.Refresh] finishes
func (d *Display) GetImageBytes() []byte {
	d.RLock()
	defer d.RUnlock()

	return d.buffer.Bytes()
}

// NeedsUpdate checks if the last update time is before the provided time for comparison.
func (d *Display) NeedsUpdate(check time.Time) bool {
	d.RLock()
	defer d.RUnlock()

	return d.lastUpdate.Before(check)
}

// Refresh generates a new image based on the provided SensorDataMessage
// It updates the display buffer with the new image and sets the last update time
func (d *Display) Refresh(last *SensorDataMessage) error {
	d.Lock()
	defer d.Unlock()

	slog.Debug("Refreshing image", "temperatur", last.Temperature.String(), "date_time", last.MessageDate.String())

	img, err := d.generateImage(last.Temperature.String(), last.MessageDate.String())
	if err != nil {
		return err
	}

	d.buffer.Reset()
	err = png.Encode(d.buffer, img)
	if err != nil {
		return err
	}

	d.lastUpdate = time.Time(last.MessageDate)

	return nil
}

// generateImage creates an image with temperature and last modified text anchored in the center.
// It requires temperature and lastModified strings as input, and returns the generated image.
func (d *Display) generateImage(temperature, lastModified string) (image.Image, error) {
	dc := gg.NewContext(d.width, d.height)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Temperature display
	dc.SetRGB(0, 0, 0)
	if err := dc.LoadFontFace("fonts/Roboto-Regular.ttf", 800); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored(temperature, float64(d.width)/2, float64(d.height)/2, 0.5, 0.5)

	// Last updated
	dc.SetRGB(0.5, 0.5, 0.5)
	if err := dc.LoadFontFace("fonts/Roboto-LightItalic.ttf", 100); err != nil {
		slog.Error("unable to load font: ", "error", err)
		return nil, err
	}
	dc.DrawStringAnchored("Last updated "+lastModified, float64(d.width)/2, float64(d.height)-100, 0.5, 0.5)

	return dc.Image(), nil
}
