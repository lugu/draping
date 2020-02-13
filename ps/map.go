package ps

import (
	"fmt"
	"image"

	"github.com/lugu/draping"
)

const (
	OuterBorder = 2
	InnerBorder = 6
)

type Stats struct {
	labels  []string
	rawData map[string][]ProcessStatus
}

func NewStats() *Stats {
	return &Stats{
		labels:  []string{},
		rawData: make(map[string][]ProcessStatus),
	}
}

func (s *Stats) Series() []Serie {
	series := make([]Serie, len(s.labels))
	for i, label := range s.labels {
		stats := s.rawData[label]
		serie := Serie{
			Label:  label,
			Values: make([]Metric, len(stats)),
		}
		for i, stat := range stats {
			serie.Values[i] = Metric{
				Percent: float64(stat.CPU),
			}
		}
		series[i] = serie
	}
	return series
}

func (s *Stats) bound() image.Rectangle {
	width := (len(s.labels)-1)*(InnerBorder+ProcessWidth) + ProcessWidth + 2*OuterBorder
	height := 0
	for _, stat := range s.rawData {
		if height < len(stat) {
			height = len(stat)
		}
	}
	height += 2 * OuterBorder
	return image.Rect(0, 0, width, height)
}

func (s *Stats) Map() *draping.Map {
	r := s.bound()
	terrain := image.NewRGBA(r)
	// black = color.RGBA{ 0, 0, 0, 0xff }
	// draping.Clear(terrain, Black)
	m := &draping.Map{
		Terrain:   terrain,
		Elevation: image.NewGray(r),
	}
	series := s.Series()

	dp := image.Point{X: ProcessWidth + InnerBorder, Y: 0}
	ri := image.Rect(OuterBorder, OuterBorder, OuterBorder+ProcessWidth, r.Max.Y-OuterBorder)

	for _, serie := range series {
		DrawSerie(m, serie, ri)
		ri = ri.Add(dp)
	}
	return m
}

func (s *Stats) Update() error {
	status, err := PS()
	if err != nil {
		return err
	}
	for _, stat := range status {
		if stat.Command == "ps" {
			continue
		}
		label := fmt.Sprintf("%s [%d]", stat.Command, stat.PID)
		raw, ok := s.rawData[label]
		if !ok {
			s.labels = append(s.labels, label)
			s.rawData[label] = []ProcessStatus{stat}
		} else {
			s.rawData[label] = append(raw, stat)
		}
	}
	return nil
}
