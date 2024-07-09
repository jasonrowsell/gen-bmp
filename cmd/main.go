package main

type Color struct {
	R, G, B float32
}

type Image struct {
	Height int
	Width  int
	Colors [][]Color
}

func NewImage(h, w int) *Image {
	img := &Image{
		Height: h,
		Width:  w,
		Colors: make([][]Color, h),
	}
	for i := range img.Colors {
		img.Colors[i] = make([]Color, w)
	}

	return img
}
