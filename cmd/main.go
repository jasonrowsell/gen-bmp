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

	filesize := 14 + 40 + 3*img.Width*img.Height // BMP header + DIB header + pixel data
	pixelDataOffset := 54                        // Pixel data offset (14+40)

	if err := writeBMPHeader(f, filesize, pixelDataOffset); err != nil {
		return fmt.Errorf("failed to write BMP header: %w", err)
	}

	if err := writeDIBHeader(f, img.Height, img.Width); err != nil {
		return fmt.Errorf("failed to write DIB header: %w", err)
	}

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

func writeBMPHeader(f *os.File, filesize, pixelDataOffset int) error {
	// BMP file header (14 B)
	header := []interface{}{
		[2]byte{'B', 'M'},
		uint32(filesize),
		uint32(0), // Reserved
		uint32(pixelDataOffset),
	}

	for _, v := range header {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	return nil
}

func writeDIBHeader(f *os.File, width, height int) error {
	// DIB file header (40 B)
	header := []interface{}{
		uint32(40), // DIB header size
		int32(width),
		int32(height),
		uint16(1),   // Color planes
		uint16(24),  // Bits per pixel
		uint32(0),   // BI_RGB: No compression
		uint32(0),   // Image size
		int32(2835), // Horizontal resolution (72 dpi)
		int32(2835), // Vertical resolution (72 dpi)
		uint32(0),   // Colors in color table
		uint32(0),   // Important color count
	}

	for _, v := range header {
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			return err
		}
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
