package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
)

type Color struct {
	R, G, B float32
}

type Image struct {
	Width  int
	Height int
	Colors [][]Color
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

	return img
}

func (img *Image) getColor(y, x int) (Color, error) {
	if y < 0 || y >= img.Height || x < 0 || x >= img.Width {
		return Color{}, errors.Errorf("out of bounds (%d, %d)", x, y)
	}
	return img.Colors[y][x], nil
}

func (img *Image) setColor(c Color, y, x int) error {
	if y < 0 || y >= img.Height || x < 0 || x >= img.Width {
		return errors.Errorf("out of bounds (%d, %d)", x, y)
	}
	img.Colors[y][x] = c
	return nil
}

func (img *Image) Export(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}
	defer f.Close()

	filesize := 14 + 40 + 3*img.Width*img.Height // BMP header + DIB header + pixel data
	pixelDataOffset := 54                        // Pixel data offset (14+40)

	if err := writeBMPHeader(f, filesize, pixelDataOffset); err != nil {
		return errors.Wrap(err, "failed to write BMP header")
	}

	if err := writeDIBHeader(f, img.Height, img.Width); err != nil {
		return errors.Wrap(err, "failed to write DIB header")
	}

	if err := writePixelData(f, img); err != nil {
		return errors.Wrap(err, "failed to write pixel data")
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
			return errors.Wrap(err, "failed to write BMP header")
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
			return errors.Wrap(err, "failed to write DIB header")
		}
	}
	return nil
}

func writePixelData(f *os.File, img *Image) error {
	padding := (4 - (img.Width*3)%4) % 4
	buffer := make([]byte, img.Width*3+padding) // Add padding to ensure each row is a multiple of 4 bytes

	for y := img.Height - 1; y >= 0; y-- { // BMP stores rows bottom-to-top
		for x := 0; x < img.Width; x++ {
			c, err := img.getColor(y, x)
			if err != nil {
				return errors.Wrapf(err, "failed to get color at (%d, %d)", x, y)
			}
			i := x * 3
			buffer[i] = byte(c.B * 255)
			buffer[i+1] = byte(c.G * 255)
			buffer[i+2] = byte(c.R * 255)
		}

		if _, err := f.Write(buffer); err != nil {
			return errors.Wrap(err, "failed to write pixel data")
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
