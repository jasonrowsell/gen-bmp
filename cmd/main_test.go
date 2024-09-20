package main

import (
	"os"
	"sync"
	"testing"
)

func TestNewImage(t *testing.T) {
	img := NewImage(5, 10)

	if img.Height != 5 || img.Width != 10 {
		t.Errorf("Expected dimensions 5x10, got %dx%d", img.Height, img.Width)
	}

	if len(img.Colors) != 5 {
		t.Errorf("Expected total Colors rows %d, got %d", 5, len(img.Colors))
	}

	for i := range img.Height {
		if len(img.Colors[i]) != 10 {
			t.Errorf("Expected Colors columns %d, got %d", 10, len(img.Colors[i]))
		}
	}
}

func TestBoundaryConditions(t *testing.T) {
	img := NewImage(5, 5)
	c := Color{R: 1, G: 0, B: 0}

	err := img.setColor(c, 0, 0)
	if err != nil {
		t.Errorf("Unexpected error setting color at (0, 0): %v", err)
	}

	_, err = img.getColor(4, 4)
	if err != nil {
		t.Errorf("Unexpected error getting color at (4, 4): %v", err)
	}

	err = img.setColor(c, -1, 0)
	if err == nil {
		t.Error("Expected error for out-of-bounds, but got nil")
	}

	_, err = img.getColor(5, 5)
	if err == nil {
		t.Error("Expected error for out-of-bounds, but got nil")
	}
}

func TestConcurrentAccess(t *testing.T) {
	img := NewImage(10, 5)
	var wg sync.WaitGroup
	c := Color{R: 0.5, G: 0.6, B: 0.7}

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			wg.Add(1)
			go func(x, y int) {
				defer wg.Done()
				err := img.setColor(c, y, x)
				if err != nil {
					t.Errorf("Unexpected error setting color: %v", err)
				}
				retrieved, err := img.getColor(y, x)
				if err != nil {
					t.Errorf("Unexpected error getting color: %v", err)
				}
				if retrieved != c {
					t.Errorf("Expected color %v, got %v", c, retrieved)
				}
			}(x, y)
		}
	}
	wg.Wait()
}

func TestExport(t *testing.T) {
	img := NewImage(10, 5)
	color := Color{R: 0.5, G: 0.5, B: 0.5}
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			img.setColor(color, y, x)
		}
	}

	tmpFile := "test_image.bmp"
	err := img.Export(tmpFile)
	if err != nil {
		t.Errorf("Unexpected error exporting image: %v", err)
	}
	defer os.Remove(tmpFile)

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it does not", tmpFile)
	}
}
