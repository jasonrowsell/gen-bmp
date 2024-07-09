package main

import (
	"fmt"
)

type Color struct {
	R, G, B float32
}

type GetColor func(x, y int) (Color, error)
type SetColor func(c Color, x, y int) error

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

	img.GetColor = func(x, y int) (Color, error) {
		if x < 0 || x >= w || y < 0 || y >= h {
			return Color{}, fmt.Errorf("out of bounds (%d, %d)", x, y)
		}
		return img.Colors[y][x], nil
	}

	img.SetColor = func(c Color, x, y int) error {
		if x < 0 || x >= w || y < 0 || y >= h {
			return fmt.Errorf("out of bounds (%d, %d)", x, y)
		}

		img.Colors[y][x] = c
		return nil
	}

	return img
}
