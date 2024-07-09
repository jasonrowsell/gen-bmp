package main

import "testing"

func TestNewImage(t *testing.T) {
	img := NewImage(5, 10)

	if img.Height != 5 || img.Width != 10 {
		t.Errorf("Expected dimensions 5x10, got %dx%d", img.Height, img.Width)
	}
}
