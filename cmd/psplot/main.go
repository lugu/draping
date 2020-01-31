package main

import (
	"math"
	"time"
	"fmt"
	"image"
	"image/draw"

	"golang.org/x/mobile/event/key"
	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"

	"github.com/lugu/draping"
	"github.com/lugu/draping/ps"
)

var (
	camera draping.Camera
	stats = ps.NewStats()
	world = stats.Map()

	hasChanged = false
	showMap = false
	showLevel = false
	screen = image.NewRGBA(image.Rect(0, 0, 0, 0))

	wnd nucular.MasterWindow = nil
)

func resetCamera(screenSize, mapSize image.Point) {

	camera.Pos = image.Point{
		mapSize.X/2,
		mapSize.Y/2 - 30,
	}

	camera.Height = 50
	// camera.Phi-= math.Pi
	camera.Distance = screenSize.X/2
	camera.Horizon = screenSize.Y/3*2
	camera.ScaleHeight = 300
}

func render(w *nucular.Window) {
	size := w.WidgetBounds()
	screen = image.NewRGBA(image.Rect(0, 0, size.W, size.H))
	world.Render(screen, camera)
	if showMap {
		draw.Draw(screen, world.Terrain.Bounds(), world.Terrain,
			image.Point{}, draw.Src)
	}
	if showLevel {
		draw.Draw(screen, world.Elevation.Bounds(), world.Elevation,
			image.Point{}, draw.Src)
	}
}

func updatefn(w *nucular.Window) {
	w.Row(0).Dynamic(1)

	if w.WidgetBounds().W != screen.Bounds().Dx() ||
		w.WidgetBounds().H != screen.Bounds().Dy() {
		screenRect := image.Rect(0, 0,
			w.WidgetBounds().W, w.WidgetBounds().H)
		screen = image.NewRGBA(screenRect)
		mapSize := world.Terrain.Bounds().Size()
		resetCamera(screenRect.Size(), mapSize)
		render(w)
	}

	input := w.Input()
	for _, e := range input.Keyboard.Keys {
		switch e.Rune {
		case 'm':
			showLevel = false
			showMap = !showMap
		case 'M':
			showMap = false
			showLevel = !showLevel
		case 'd':
			camera.Distance+=10
		case 'D':
			camera.Distance-=10
		case 's':
			camera.ScaleHeight+=1.0
		case 'S':
			camera.ScaleHeight-=1.0
		case 'h':
			camera.Horizon+=10
		case 'H':
			camera.Horizon-=10
		case 'e':
			camera.Height+=+10
		case 'E':
			camera.Height-=10
		case 'r':
			camera.Phi+= math.Pi/6.0
		case 'R':
			camera.Phi-= math.Pi/6.0
		case 'p':
			fmt.Printf("%#v\n", camera)
		}
		switch e.Code {
		case key.CodeQ, key.CodeEscape: // quit
			wnd.Close()
			return
		case key.CodeLeftArrow:
			camera.Pos.X += 20
		case key.CodeRightArrow:
			camera.Pos.X -= 20
		case key.CodeUpArrow:
			camera.Pos.Y += 20
		case key.CodeDownArrow:
			camera.Pos.Y -= 20
		}
	}

	if hasChanged || len(input.Keyboard.Keys) != 0 {
		hasChanged = false
		world = stats.Map()
		render(w)
	}
	w.Image(screen)
}

func main() {

	tick := time.NewTicker(200 * time.Millisecond)
	go func() {
		for {
			<-tick.C
			wnd.Changed()
			err := stats.Update()
			if err != nil {
				panic(err)
			}
			hasChanged = true
			wnd.Changed()
		}
	}()

	wnd = nucular.NewMasterWindow(0, "Render", updatefn)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 2.0))
	wnd.Main()
}
