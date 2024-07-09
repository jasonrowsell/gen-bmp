package main

import (
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

	c := Color{R: 0.4, G: 0.5, B: 0.6}
	err := img.SetColor(c, 9, 4)
	if err != nil {
		t.Errorf("Unexpected error setting color: %v", err)
	}
	retrieved, err := img.GetColor(9, 4)
	if err != nil {
		t.Errorf("Unexpected error getting color: %v", err)
	}
	if retrieved != c {
		t.Errorf("Expected color %v, got %v", c, retrieved)
	}

	err = img.SetColor(c, 10, 5)
	if err == nil {
		t.Errorf("Expected error for out-of-bounds set, got nil")
	}
	_, err = img.GetColor(10, 5)
	if err == nil {
		t.Errorf("Expected error for out-of-bounds get, got nil")
	}
}
