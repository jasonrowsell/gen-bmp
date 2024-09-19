package main

import (
	"encoding/binary"
	"fmt"
	"os"
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

func (img *Image) Export(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	// BMP header 14 B
	// DIB header 40 B
	// pixel data 3 B per pixel (BGR)
	filesize := 14 + 40 + 3*img.Width*img.Height

	// BMP file header (14 B)
	f.Write([]byte("BM")) // Signature
	binary.Write(f, binary.LittleEndian, uint32(filesize))
	binary.Write(f, binary.LittleEndian, uint32(0))  // Reserved
	binary.Write(f, binary.LittleEndian, uint32(54)) // Pixel data offset (14+40)

	// DIB header (40 B)
	binary.Write(f, binary.LittleEndian, uint32(40)) // DIB header size
	binary.Write(f, binary.LittleEndian, uint32(img.Width))
	binary.Write(f, binary.LittleEndian, uint32(img.Height))
	binary.Write(f, binary.LittleEndian, uint16(1))   // Color planes
	binary.Write(f, binary.LittleEndian, uint16(24))  // Bits per pixel
	binary.Write(f, binary.LittleEndian, uint32(0))   // BI_RGB: No compression
	binary.Write(f, binary.LittleEndian, uint32(16))  // Image size
	binary.Write(f, binary.LittleEndian, int64(2835)) // Horizontal resolution (72 dpi)
	binary.Write(f, binary.LittleEndian, int64(2835)) // Vertical resolution (72 dpi)
	binary.Write(f, binary.LittleEndian, uint32(0))   // Colors in color table
	binary.Write(f, binary.LittleEndian, uint32(0))   // Important color count

	return nil
}
