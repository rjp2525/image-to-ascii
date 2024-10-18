package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestLoadImage(t *testing.T) {
	testFile := "test_input.png"
	createTestImage(testFile, 100, 100)

	defer os.Remove(testFile)

	img := LoadImage(testFile)
	if img == nil {
		t.Fatalf("Failed to load image: %v", testFile)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Loaded image dimensions are incorrect, got: %dx%d, want: 100x100", bounds.Dx(), bounds.Dy())
	}
}

func TestResizeImage(t *testing.T) {
	img := createTestGrayImage(100, 100)

	resized := ResizeImage(img, 50)
	if resized == nil {
		t.Fatal("ResizeImage returned nil")
	}

	bounds := resized.Bounds()
	if bounds.Dx() != 50 {
		t.Errorf("Resized image width is incorrect, got: %d, want: 50", bounds.Dx())
	}
}

func TestConvGrayScale(t *testing.T) {
	img := createTestColorImage(100, 100)

	gray := ConvGrayScale(img)
	if gray == nil {
		t.Fatal("ConvGrayScale returned nil")
	}

	for x := 0; x < 100; x++ {
		for y := 0; y < 100; y++ {
			r, g, b, _ := gray.At(x, y).RGBA()
			if r != g || g != b {
				t.Errorf("Pixel at (%d, %d) is not grayscale", x, y)
			}
		}
	}
}

func TestMapAscii(t *testing.T) {
	img := createTestGrayImage(100, 50)

	asciiArt := MapAscii(img)
	if len(asciiArt) == 0 {
		t.Fatal("MapAscii returned empty ASCII art")
	}

	expectedWidth := 100
	for _, line := range asciiArt {
		if len(line) != expectedWidth {
			t.Errorf("Line length incorrect, got: %d, want: %d", len(line), expectedWidth)
		}
	}
}

func createTestGrayImage(width, height int) image.Image {
	img := image.NewGray(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grayValue := uint8((x + y) % 256)
			img.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}
	return img
}

func createTestColorImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			col := color.RGBA{R: uint8(x), G: uint8(y), B: 0, A: 255}
			img.Set(x, y, col)
		}
	}
	return img
}

func createTestImage(fileName string, width, height int) {
	img := createTestGrayImage(width, height)

	file, err := os.Create(fileName)
	if err != nil {
		panic("Failed to create test image")
	}
	defer file.Close()

	png.Encode(file, img)
}
