package main

type Color struct {
	R, G, B float32
}

type GetColor func(x, y int) Color
type SetColor func(c Color, x, y int)

type Image struct {
	Height   int
	Width    int
	Colors   [][]Color
	GetColor GetColor
	SetColor SetColor
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

	img.GetColor = func(x, y int) Color {
		return img.Colors[y][x]
	}

	img.SetColor = func(c Color, x, y int) {
		img.Colors[y][x] = c
	}

	return img
}
