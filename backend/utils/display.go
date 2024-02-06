package utils

import (
	"image"

	"github.com/kbinani/screenshot"
)

func GetDisplayBounds() []image.Rectangle {
	n := screenshot.NumActiveDisplays()

	bounds := make([]image.Rectangle, 0, n)

	for i := 0; i < n; i++ {
		bounds = append(bounds, screenshot.GetDisplayBounds(i))
	}

	return bounds
}
