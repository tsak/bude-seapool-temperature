package main

import (
	"bytes"
	"image"
	"image/png"
	"log/slog"
	"sync"
	"time"
)

// ImageGenerator generates images for the set width and height, using its generateImage function
type ImageGenerator struct {
	sync.RWMutex
	width         int
	height        int
	msg           string
	buffer        *bytes.Buffer
	lastUpdate    time.Time
	generateImage func(width, height int, temperature, lastModified, msg string) (image.Image, error)
}

// NewImageGenerator creates a new display with the specified width and height, initializing the display buffer.
func NewImageGenerator(width, height int, msg string, generateImage func(width, height int, temperature, lastModified, msg string) (image.Image, error)) *ImageGenerator {
	return &ImageGenerator{
		width:         width,
		height:        height,
		msg:           msg,
		buffer:        bytes.NewBuffer([]byte{}),
		generateImage: generateImage,
	}
}

// GetImageBytes returns the image data as a byte slice by reading from the display buffer.
// It will be blocked while a call to [ImageGenerator.Refresh] finishes
func (ig *ImageGenerator) GetImageBytes() []byte {
	ig.RLock()
	defer ig.RUnlock()

	return ig.buffer.Bytes()
}

// NeedsUpdate checks if the last update time is before the provided time for comparison.
func (ig *ImageGenerator) NeedsUpdate(check time.Time) bool {
	ig.RLock()
	defer ig.RUnlock()

	return ig.lastUpdate.Before(check)
}

// Refresh generates a new image based on the provided SensorDataMessage
// It updates the display buffer with the new image and sets the last update time
func (ig *ImageGenerator) Refresh(last *SensorDataMessage) error {
	ig.Lock()
	defer ig.Unlock()

	slog.Debug("Refreshing image", "temperature", last.Temperature.String(), "date_time", last.MessageDate.String())

	img, err := ig.generateImage(ig.width, ig.height, last.Temperature.String(), last.MessageDate.String(), ig.msg)
	if err != nil {
		return err
	}

	ig.buffer.Reset()
	err = png.Encode(ig.buffer, img)
	if err != nil {
		return err
	}

	ig.lastUpdate = time.Time(last.MessageDate)

	return nil
}
