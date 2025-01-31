package main

import (
	"flag"
	"log"
	"rover/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	width := flag.Int("width", 800, "Board width in pixels")
	height := flag.Int("height", 600, "Board height in pixels")
	speed := flag.Int("speed", 30, "Ticks per second, min 0 max 60, + or - to adjust in game")
	debug := flag.Bool("debug", false, "Show debug info, D to toggle in game")

	flag.Parse()

	g, err := game.NewGame(*width, *height, *speed, *debug)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
