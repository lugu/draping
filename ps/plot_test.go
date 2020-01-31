package ps

import (
	"testing"
	"image"
	"image/color"
	"time"

	"github.com/lugu/draping"
)

func TestDrawMetric(t *testing.T) {

	size :=image.Rect(0, 0, 30, 120)
	mp := &draping.Map {
		Terrain: image.NewRGBA(size),
		Elevation: image.NewGray(size),
	}

	serie := Serie{
		Label: "test",
		Values: make([]Metric, 100),
	}
	for i := 0; i < len(serie.Values); i++ {
		serie.Values[i] = Metric {
			Percent: float64(i),
			TS: time.Now(),
		}
	}

	DrawSerie(mp, serie, image.Rect(10, 10, 20, 110))


	for i := 0; i < len(serie.Values); i++ {
		y := 10 + i
		for x := 0; x < 30; x++ {
			e := mp.Elevation.GrayAt(x, y)
			c := mp.Terrain.At(x, y)
			r, g, b, a := c.RGBA()
			if x < 10 || x >= 20 {
				if r != 0 || g != 0 || b != 0 || a != 0 || e.Y != 0 {
					t.Errorf("(x:%d, y:%d) e:%d, r:%d, g:%d, b:%d, a:%d", x, y, e.Y, r, g, b, a)
				}
			} else {
				e0 := uint8(float64(i)/100*ElevationMax)
				r0, g0, b0, a0 := Tenth[i/10].RGBA()
				if r != r0 || g != g0 || b != b0 || a != a0 || e.Y != e0 {
					t.Errorf("(x:%d, y:%d) e:%d, r:%d, g:%d, b:%d, a:%d instead of e: %d, r:%d, g:%d, b:%d, a:%d", x, y, e.Y, r, g, b, a, e0, r0, g0, b0, a0)
				}
			}
		}
	}
}

func helpTestElevation(t *testing.T, percent float64, heigh uint8) {
	e := Elevation(percent)
	if e.Y != heigh {
		t.Errorf("%f => %d instead of %d", percent, e.Y, heigh)
	}
}

func TestElevation(t *testing.T) {
	helpTestElevation(t, 100, ElevationMax)
	helpTestElevation(t, 0, 0)
	helpTestElevation(t, 50, ElevationMax/2)
}

func helpTestColor(t *testing.T, percent float64, c0 color.RGBA) {
	c := Color(percent)
	r, g, b, a := c.RGBA()
	r0, g0, b0, a0 := c0.RGBA()
	if r != r0 || g != g0 || b != b0 || a != a0 {
		t.Errorf("%f => r:%d, g:%d, b:%d, a:%d instead of r:%d, g:%d, b:%d, a:%d", percent, r, g, b, a, r0, g0, b0, a0)
	}
}

func TestColor(t *testing.T) {
	helpTestColor(t, 0, Tenth[0])
	helpTestColor(t, 9, Tenth[0])
	helpTestColor(t, 10, Tenth[1])
	helpTestColor(t, 50, Tenth[5])
	helpTestColor(t, 100, Tenth[10])
}
