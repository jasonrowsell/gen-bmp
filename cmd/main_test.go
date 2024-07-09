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
}
