package utils

import (
	"maps"
	"time"
)

type ProgressTracker struct {
	windowSize time.Duration
	data       map[time.Time]int64
	Total      int64
}

func NewProgressTracker(windowSize time.Duration) *ProgressTracker {
	return &ProgressTracker{
		windowSize: windowSize,
		data:       make(map[time.Time]int64),
	}
}

func (pt *ProgressTracker) Add(value int64) {
	pt.Total += value
	pt.data[time.Now()] = value
	pt.evict()
}

func (pt *ProgressTracker) Speed() float64 {
	pt.evict()
	var first, last time.Time
	for t := range pt.data {
		if first.IsZero() || t.Before(first) {
			first = t
		}
		if last.IsZero() || t.After(last) {
			last = t
		}
	}
	firstData := pt.data[first]
	lastData := pt.data[last]
	return float64(lastData-firstData) / pt.windowSize.Seconds()
}

func (pt *ProgressTracker) ETA() time.Duration {
	speed := pt.Speed()
	if speed == 0 {
		return 0
	}
	var latest time.Time
	for t := range pt.data {
		if latest.IsZero() || t.After(latest) {
			latest = t
		}
	}
	latestData := pt.data[latest]
	return time.Duration(float64(pt.Total-latestData)/speed) * time.Second
}

func (pt *ProgressTracker) evict() {
	maps.DeleteFunc(pt.data, func(t time.Time, _ int64) bool {
		return time.Since(t) > pt.windowSize
	})
}
