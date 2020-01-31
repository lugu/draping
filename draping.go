package draping

import (
	"os"
	"fmt"
	"math"
	"errors"
	"image"
	"image/draw"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

var (
	ErrInvalidFormat = errors.New("invalid format")
)

var (
	BackgroundColor = color.RGBA{ 128, 128, 128, 255 }
)

// LoadImage retrieves and decodes the image file representing the
// color of the ground.
func LoadImage(input string) (draw.Image, error) {
	file, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	dst := image.NewRGBA(img.Bounds())
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Src)
	return dst, nil
}

// LoadLevel return a gray image representing a level elevation.
func LoadLevel(input string) (*image.Gray, error) {
	src, err := LoadImage(input)
	if err != nil {
		return nil, err
	}
	size := src.Bounds().Size()
	gray := image.NewGray(src.Bounds())
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			gray.Set(x, y, src.At(x, y))
		}
	}
	return gray, nil
}

// Clear paint an image with a color.
func Clear(d draw.Image, c color.Color) {
	draw.Draw(d, d.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

type Map struct {
	Terrain draw.Image
	Elevation *image.Gray
}

type Camera struct {
	Pos image.Point
	Height int // camera z position
	Horizon int // heigh of the horizon line
	Distance int // max distance to render
	ScaleHeight float64
	Phi float64 // view angle
}

func DrawVerticalLine(d draw.Image, x, yMax, yMin int, c color.Color) {
	screenHeight := d.Bounds().Dy()
	drect := image.Rect(x, screenHeight - yMax, x+1, screenHeight - yMin)
	draw.Draw(d,  drect, &image.Uniform{c}, image.Point{}, draw.Src)
}

func (m *Map) Render(d draw.Image, c Camera) {

	Clear(d, BackgroundColor)

	/*
	posZ := int(m.Elevation.GrayAt(c.Pos.X, c.Pos.Y).Y)
	if posZ > c.Height {
		c.Height = posZ
	}
	*/

	screenWidth := d.Bounds().Dx()

	yBuffer := make([]int, screenWidth)
	for y := 0; y < screenWidth; y++ {
		yBuffer[y] = 0
	}

	sinphi := math.Sin(c.Phi)
	cosphi := math.Cos(c.Phi)

	Z := 1.0
	dz := 1.0

	for Z < float64(c.Distance) {

		pleftX := ( cosphi*Z + sinphi*Z) + float64(c.Pos.X)
		pleftY := (-sinphi*Z + cosphi*Z) + float64(c.Pos.Y)

		prightX := (-cosphi*Z + sinphi*Z) + float64(c.Pos.X)
		prightY := ( sinphi*Z + cosphi*Z) + float64(c.Pos.Y)

		dx := (prightX - pleftX) / float64(screenWidth)
		dy := (prightY - pleftY) / float64(screenWidth)

		for i := 0; i < screenWidth; i++ {
			X := int(pleftX + dx*float64(i))
			Y := int(pleftY + dy*float64(i))

			mapSize := m.Elevation.Bounds().Size()
			if X < - 2 || X > mapSize.X + 1 || Y < -2 || Y > mapSize.Y + 1 {
				continue
			}

			elevation := int(m.Elevation.GrayAt(X, Y).Y)
			col := m.Terrain.At(X, Y)
			if X < 0 || X > mapSize.X - 1 || Y < 0 || Y > mapSize.Y - 1 {
				col = BackgroundColor
			}

			heighOnScreen := int((float64(elevation - c.Height) / Z) * c.ScaleHeight) + c.Horizon
			if heighOnScreen >  yBuffer[i] {
				DrawVerticalLine(d, i, heighOnScreen, yBuffer[i], col)
				yBuffer[i] = heighOnScreen
			}
		}
		Z = Z + dz
		// dz = dz + 0.2
	}
}
