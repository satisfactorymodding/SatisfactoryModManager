package utils

import (
	"image"
	"runtime"

	"github.com/kbinani/screenshot"
)

func GetDisplayBounds() []image.Rectangle {
	n := screenshot.NumActiveDisplays()

	bounds := make([]image.Rectangle, 0, n)

	for i := 0; i < n; i++ {
		bounds = append(bounds, screenshot.GetDisplayBounds(i))
	}

	if runtime.GOOS == "linux" {
		// gdk_monitor_get_geometry considers 0,0 to be the corner of the bounding box of all the monitors,
		// not the 0,0 of the main monitor
		boundingBox := bounds[0]
		for _, b := range bounds {
			boundingBox = boundingBox.Union(b)
		}
		for i := range bounds {
			bounds[i] = bounds[i].Sub(boundingBox.Min)
		}
	}

	return bounds
}

func GetDisplayBoundsAt(x, y int) image.Rectangle {
	point := image.Pt(x, y)

	displays := GetDisplayBounds()

	curDisplay := displays[0] // use main display as fallback

	for _, d := range displays {
		if point.In(d) {
			curDisplay = d
			break
		}
	}

	return curDisplay
}
