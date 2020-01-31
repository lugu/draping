package ps

import (
	"time"
	"math"
	"fmt"
	"image/color"
	"image"

	"github.com/lugu/draping"
)

const (
	ProcessWidth = 10
	ElevationMax = 200
)

var (
	Tenth = []color.RGBA{
		color.RGBA{ 0x55, 0xff, 0x00, 0xff }, //   0
		color.RGBA{ 0x80, 0xff, 0x00, 0xff }, //  10
		color.RGBA{ 0xaa, 0xff, 0x00, 0xff }, //  20
		color.RGBA{ 0xd4, 0xff, 0x00, 0xff }, //  30
		color.RGBA{ 0xff, 0xff, 0x00, 0xff }, //  40
		color.RGBA{ 0xff, 0xd5, 0x00, 0xff }, //  50
		color.RGBA{ 0xff, 0xaa, 0x00, 0xff }, //  60
		color.RGBA{ 0xff, 0x80, 0x00, 0xff }, //  60
		color.RGBA{ 0xff, 0x55, 0x00, 0xff }, //  80
		color.RGBA{ 0xff, 0x2a, 0x00, 0xff }, //  90
		color.RGBA{ 0xff, 0x00, 0x00, 0xff }, // 100
	}
)

type Metric struct {
	Percent float64
	TS time.Time
}

type Serie struct {
	Label string
	Values []Metric
}

// convention: split on the y axis
func DrawSerie(m *draping.Map, s Serie, r image.Rectangle) {
	dy := float64(r.Dy()) / (float64(len(s.Values)))
	ri := image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y + int(dy))
	dp := image.Point{ X: 0, Y: int(dy) }
	for _, metric := range s.Values {
		DrawMetric(m, metric, ri)
		ri = ri.Add(dp)
	}
}

func Color(percent float64) color.RGBA {
	i := int(math.Floor(percent/10.0))
	return Tenth[i]
}

func Elevation(percent float64) color.Gray {
	return color.Gray{ Y: uint8(percent/100*ElevationMax) }
}

func levels(percent float64, width int) []float64 {
	levels := make([]float64, width)
	epsilon := math.Pi / float64(width)
	for i := 0; i < width; i++ {
		theta := float64(i) * epsilon
		coef := math.Abs(math.Sin(theta))
		levels[i] = coef * percent
		if math.IsNaN(levels[i]) || levels[i] < 0 || levels[i] > 100 {
			msg := fmt.Sprintf("invalid level: %f, %d, %d, %f\n",
				percent, width, i, levels[i])
			panic(msg)
		}
	}
	return levels
}

func DrawMetric(m *draping.Map, metric Metric, rect image.Rectangle) {
	width := rect.Max.X - rect.Min.X
	l := levels(metric.Percent, width)
	for y := rect.Min.Y; y < rect.Max.Y ; y ++ {
		for w := 0; w < width; w++ {
			p := l[w]
			c := Color(p)
			e := Elevation(p)
			x := w + rect.Min.X
			m.Elevation.Set(x, y, e)
			m.Terrain.Set(x, y, c)
		}
	}
}
