package draping

import (
	"testing"
	"image"
	"image/color"
	"path/filepath"
)

func helpTestLoapLevel(t *testing.T, filename string) {
	path := filepath.Join("testdata", filename)
	img, err := LoadLevel(path)
	if err != nil {
		t.Errorf("%s: %s", filename, err)
	} else if img == nil {
		t.Error("shound not be nil")
	}
}

func TestLoapLevel(t *testing.T) {
	helpTestLoapLevel(t, "green_512_512_gray.png")
	helpTestLoapLevel(t, "orange_128_256_gray.png")
}

func helpTestLoapMap(t *testing.T, filename string) {
	path := filepath.Join("testdata", filename)
	img, err := LoadImage(path)
	if err != nil {
		t.Errorf("%s: %s", filename, err)
	} else if img == nil {
		t.Error("shound not be nil")
	}
}

func TestLoapMap(t *testing.T) {
	helpTestLoapMap(t, "green_512_512.png")
	helpTestLoapMap(t, "orange_128_256.png")
}

func TestClear(t *testing.T) {
	dst := image.NewNRGBA(image.Rect(0, 0, 128, 128))
	blue := color.RGBA{0, 0, 255, 255}
	size := dst.Bounds().Size()
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			dst.Set(x, y, blue)
		}
	}
	Clear(dst, BackgroundColor)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, a := dst.At(x, y).RGBA()
			r0, g0, b0, a0 := BackgroundColor.RGBA()
			if r != r0 || g != g0 || b != b0 || a != a0 {
				t.Errorf("invalid color at (%d,%d): %#v instead of %#v",
				x, y, dst.At(x, y), BackgroundColor)
			}
		}
	}
}

func helpTestDrawVerticalLine(t *testing.T, rect image.Rectangle, X, yMax int) {
	dst := image.NewNRGBA(rect)
	size := dst.Bounds().Size()
	red := color.RGBA{255, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	Clear(dst, red)
	DrawVerticalLine(dst, X, yMax, green)
	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			r, g, b, a := dst.At(x, y).RGBA()
			r0, g0, b0, a0 := red.RGBA()
			if x == X && y >= size.Y - yMax {
				r0, g0, b0, a0 = green.RGBA()
			}
			if r != r0 || g != g0 || b != b0 || a != a0 {
				t.Errorf("invalid color at (%d,%d): %#v instead of %#v",
				x, y, dst.At(x, y), BackgroundColor)
			}
		}
	}
}
func TestDrawVerticalLine(t *testing.T) {
	rect := image.Rect(0, 0, 128, 128)
	size := rect.Size()
	for x := 0; x < size.X; x += 5 {
		for y := 0; y < size.Y; y += 5 {
			helpTestDrawVerticalLine(t, rect, x, y)
		}
	}
}

func helpTestRencder(t *testing.T, mapFilename, elevationFilename string) {
	path := filepath.Join("testdata", mapFilename)
	terrain, err := LoadImage(path)
	if err != nil {
		t.Errorf("%s: %s", mapFilename, err)
	}
	path = filepath.Join("testdata", elevationFilename)
	elevation, err := LoadLevel(path)
	if err != nil {
		t.Errorf("%s: %s", elevationFilename, err)
	}
	m := Map {
		Terrain:terrain,
		Elevation: elevation,
	}
	rect := image.Rect(0, 0, 128, 128)
	dst := image.NewNRGBA(rect)
	cam := Camera { image.Point{}, 120, 300, 300 }
	m.Render(dst, cam)
}

func TestRender(t *testing.T) {
	helpTestRencder(t,"green_512_512.png", "green_512_512_gray.png")
}
