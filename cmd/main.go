package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

type Color struct {
	R, G, B float32
}

type Image struct {
	Width    int
	Height   int
	Colors   [][]Color
	getColor func(y, x int) (Color, error)
	setColor func(c Color, y, x int) error
}

func NewImage(height, width int) *Image {
	img := &Image{
		Width:  width,
		Height: height,
		Colors: make([][]Color, height),
	}
	for i := range img.Colors {
		img.Colors[i] = make([]Color, width)
	}

	img.getColor = func(y, x int) (Color, error) {
		if y < 0 || y >= height || x < 0 || x >= width {
			return Color{}, fmt.Errorf("out of bounds (%d, %d)", y, x)
		}
		return img.Colors[y][x], nil
	}

	img.setColor = func(c Color, y, x int) error {
		if y < 0 || y >= height || x < 0 || x >= width {
			return fmt.Errorf("out of bounds (%d, %d)", y, x)
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

	// BMP header 14 B + DIB header 40 B + Pixel data 3 B per pixel (BGR)
	filesize := 14 + 40 + 3*img.Width*img.Height

	// BMP file header (14 B)
	f.Write([]byte("BM")) // Signature
	binary.Write(f, binary.LittleEndian, uint32(filesize))
	binary.Write(f, binary.LittleEndian, uint32(0))  // Reserved
	binary.Write(f, binary.LittleEndian, uint32(54)) // Pixel data offset (14+40)

	// DIB header (40 B)
	binary.Write(f, binary.LittleEndian, uint32(40)) // DIB header size
	binary.Write(f, binary.LittleEndian, int32(img.Width))
	binary.Write(f, binary.LittleEndian, int32(img.Height))
	binary.Write(f, binary.LittleEndian, uint16(1))   // Color planes
	binary.Write(f, binary.LittleEndian, uint16(24))  // Bits per pixel
	binary.Write(f, binary.LittleEndian, uint32(0))   // BI_RGB: No compression
	binary.Write(f, binary.LittleEndian, uint32(0))   // Image size
	binary.Write(f, binary.LittleEndian, int32(2835)) // Horizontal resolution (72 dpi)
	binary.Write(f, binary.LittleEndian, int32(2835)) // Vertical resolution (72 dpi)
	binary.Write(f, binary.LittleEndian, uint32(0))   // Colors in color table
	binary.Write(f, binary.LittleEndian, uint32(0))   // Important color count

	// Pixel data
	padding := (4 - (img.Width*3)%4) % 4
	for y := img.Height - 1; y >= 0; y-- { // BMP stores rows bottom-to-top
		for x := 0; x < img.Width; x++ {
			c, err := img.getColor(y, x)
			if err == nil {
				_, err := f.Write([]byte{
					byte(c.B * 255),
					byte(c.G * 255),
					byte(c.R * 255),
				})
				if err != nil {
					return fmt.Errorf("failed to write pixel data: %w", err)
				}
			}
		}
		// Add padding to ensure each row is a multiple of 4 bytes
		f.Write(make([]byte, padding))
	}
	return nil
}

func main() {
	img := NewImage(1000, 1000)

	// Fill image with gradient
	var wg sync.WaitGroup
	wg.Add(img.Height)
	for y := 0; y < img.Height; y++ {
		go func(y int) {
			defer wg.Done()
			for x := 0; x < img.Width; x++ {
				img.setColor(Color{
					R: float32(y) / float32(img.Height),
					G: float32(x) / float32(img.Width),
					B: 0.5,
				}, y, x)
			}
		}(y)
	}
	wg.Wait()

	if err := img.Export("gradient.bmp"); err != nil {
		fmt.Printf("Failed to export image: %v\n", err)
		return
	}

	fmt.Println("Image exported successfully")
}
