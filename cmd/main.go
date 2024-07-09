package main

type Image struct {
	Height int
	Width  int
}

func NewImage(h, w int) *Image {
	img := &Image{
		Height: h,
		Width:  w,
	}

	return img
}
