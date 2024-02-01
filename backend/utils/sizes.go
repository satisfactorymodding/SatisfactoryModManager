package utils

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var (
	UnexpandedMin     = Size{Width: 660, Height: 350}
	UnexpandedMax     = Size{Width: 0, Height: 0}
	UnexpandedDefault = Size{Width: 950, Height: 950}
	ExpandedMin       = Size{Width: 720, Height: 350}
	ExpandedMax       = Size{Width: 0, Height: 0}
	ExpandedDefault   = Size{Width: 1300, Height: 950}
)
