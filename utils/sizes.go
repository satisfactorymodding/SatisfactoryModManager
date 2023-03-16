package utils

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

var (
	UnexpandedMin     = Size{Width: 850, Height: 750}
	UnexpandedMax     = Size{Width: 0, Height: 0}
	UnexpandedDefault = Size{Width: 950, Height: 950}
	ExpandedMin       = Size{Width: 1255, Height: 750}
	ExpandedMax       = Size{Width: 0, Height: 0}
	ExpandedDefault   = Size{Width: 1300, Height: 950}
)
