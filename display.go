package main

import (
	"bytes"
	"image"
	"image/png"
	"log/slog"
	"sync"
	"time"
)

// Display represents a display with width, height, buffer, and last update time.
type Display struct {
	sync.RWMutex
	width         int
	height        int
	buffer        *bytes.Buffer
	lastUpdate    time.Time
	generateImage func(width, height int, temperature, lastModified string) (image.Image, error)
}

// NewDisplay creates a new display with the specified width and height, initializing the display buffer.
func NewDisplay(width, height int, generateImage func(width, height int, temperature, lastModified string) (image.Image, error)) *Display {
	return &Display{
		width:         width,
		height:        height,
		buffer:        bytes.NewBuffer([]byte{}),
		generateImage: generateImage,
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

	img, err := d.generateImage(d.width, d.height, last.Temperature.String(), last.MessageDate.String())
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
