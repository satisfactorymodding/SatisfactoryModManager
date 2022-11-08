package utils

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

var UnexpandedMin = Size{Width: 850, Height: 750}
var UnexpandedMax = Size{Width: 0, Height: 0}
var UnexpandedDefault = Size{Width: 950, Height: 950}
var ExpandedMin = Size{Width: 1255, Height: 750}
var ExpandedMax = Size{Width: 0, Height: 0}
var ExpandedDefault = Size{Width: 1300, Height: 950}
