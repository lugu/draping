package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/style"
	"github.com/lugu/draping"
	"golang.org/x/mobile/event/key"
)

var (
	camera draping.Camera
	world  draping.Map

	showMap   = false
	showLevel = false
	screen    = image.NewRGBA(image.Rect(0, 0, 0, 0))

	wnd nucular.MasterWindow = nil
)

func resetCamera(screenSize, mapSize image.Point) {

	camera.Pos = image.Point{
		mapSize.X / 2,
		mapSize.Y / 2,
	}

	camera.Height = 50
	camera.Distance = screenSize.X / 2
	camera.Horizon = screenSize.Y / 3 * 2
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
			camera.Distance += 10
		case 'D':
			camera.Distance -= 10
		case 's':
			camera.ScaleHeight += 1.0
		case 'S':
			camera.ScaleHeight -= 1.0
		case 'h':
			camera.Horizon += 10
		case 'H':
			camera.Horizon -= 10
		case 'e':
			camera.Height += +10
		case 'E':
			camera.Height -= 10
		case 'r':
			camera.Phi += math.Pi / 6.0
		case 'R':
			camera.Phi -= math.Pi / 6.0
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

	if len(input.Keyboard.Keys) != 0 {
		render(w)
	}
	w.Image(screen)
}

func main() {
	var mapFilename = flag.String("m", "", "Map file")
	var elevationFilename = flag.String("e", "", "Elevation map file")
	flag.Parse()

	if *mapFilename == "" || *elevationFilename == "" {
		log.Fatalf("Missing map files")
	}
	terrain, err := draping.LoadImage(*mapFilename)
	if err != nil {
		log.Fatalf("cannot load map %s: %s", *mapFilename, err)
	}
	elevation, err := draping.LoadLevel(*elevationFilename)
	if err != nil {
		log.Fatalf("cannot load elevation map %s: %s",
			*elevationFilename, err)
	}

	world.Terrain = terrain
	world.Elevation = elevation

	wnd = nucular.NewMasterWindow(0, "Render", updatefn)
	wnd.SetStyle(style.FromTheme(style.DarkTheme, 2.0))
	wnd.Main()
}
