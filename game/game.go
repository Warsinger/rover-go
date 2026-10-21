package game

import (
	"fmt"
	"image/color"
	"rover/assets"
	"rover/internal/platform"
	"rover/internal/simulation"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameData struct {
	*simulation.State
	oldwidth, oldheight int
}

const (
	minSpeed    = 0
	maxSpeed    = 60
	frameCount  = 4
	frameOffset = 1
)

func NewGame(width, height, gamespeed int, debug bool) (*GameData, error) {
	state, err := simulation.New(simulation.Config{Width: width, Height: height, Speed: gamespeed}, debug)
	if err != nil {
		return nil, err
	}
	if err := assets.LoadAssets(); err != nil {
		return nil, err
	}
	ebiten.SetWindowTitle("Rover")
	ebiten.SetWindowSize(width, height)
	return &GameData{State: state}, nil
}

func (g *GameData) Draw(screen *ebiten.Image) {
	screen.Clear()

	background := assets.GetImage("backgroundH")
	bgScaleX := float64(g.Width) / float64(background.Bounds().Dx())
	bgScaleY := float64(g.Height) / float64(background.Bounds().Dy())
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(bgScaleX, bgScaleY)
	var drawPosition1 float64 = -float64(g.Position.X % int(float64(g.Width)))
	if g.Position.X < 0 {
		drawPosition1 -= float64(g.Width)
	}
	opts.GeoM.Translate(drawPosition1, 0)
	screen.DrawImage(background, opts)
	opts.GeoM.Translate(float64(g.Width), 0)
	screen.DrawImage(background, opts)

	rover := assets.GetImage(fmt.Sprintf("rover%d", g.Frame%frameCount+frameOffset))
	opts = &ebiten.DrawImageOptions{}
	scale := 0.5
	opts.GeoM.Scale(scale*g.Direction, scale)
	xmin := float64(g.Width/2) - g.Direction*float64(rover.Bounds().Dx()/2)*scale
	ymin := float64(g.Position.Y) - float64(rover.Bounds().Dy()/2)*scale
	opts.GeoM.Translate(xmin, ymin)
	screen.DrawImage(rover, opts)

	if g.IsMoving() {
		// draw raocket blast out the end of the rover
		var flameLength float64 = 15
		const yoffset = -15
		var yangleOffset float32 = 0
		if g.IsJumping() {
			yangleOffset = 10
		}
		vector.StrokeLine(screen, float32(xmin+g.Direction*(flameLength-4)), float32(ymin+float64(rover.Bounds().Dy()/2)*scale)-4+yoffset, float32(xmin), yangleOffset+float32(ymin+float64(rover.Bounds().Dy()/2)*scale)-4+yoffset, 2, color.RGBA{255, 165, 0, 255}, false)
		vector.StrokeLine(screen, float32(xmin+g.Direction*(flameLength-2)), float32(ymin+float64(rover.Bounds().Dy()/2)*scale)-2+yoffset, float32(xmin), yangleOffset+float32(ymin+float64(rover.Bounds().Dy()/2)*scale)-2+yoffset, 2, color.RGBA{255, 255, 0, 255}, false)
		vector.StrokeLine(screen, float32(xmin+g.Direction*(flameLength)), float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+yoffset, float32(xmin), yangleOffset+float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+yoffset, 2, color.RGBA{255, 0, 0, 255}, false)
		vector.StrokeLine(screen, float32(xmin+g.Direction*(flameLength-2)), float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+2+yoffset, float32(xmin), yangleOffset+float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+2+yoffset, 2, color.RGBA{255, 255, 0, 255}, false)
		vector.StrokeLine(screen, float32(xmin+g.Direction*(flameLength-4)), float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+4+yoffset, float32(xmin), yangleOffset+float32(ymin+float64(rover.Bounds().Dy()/2)*scale)+4+yoffset, 2, color.RGBA{255, 165, 0, 255}, false)
	}

	if g.Debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("jump time: %d\njumping %v\ngame speed %d", g.JumpTime, g.IsJumping(), g.GameSpeed))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("x, y %d, %d\nxscale, yscale %f, %f\nxmin, ymin %f, %f", g.Position.X, g.Position.Y, bgScaleX, bgScaleY, xmin, ymin), g.Width/2, 0)

		vector.StrokeCircle(screen, float32(g.Width/2), float32(g.Position.Y), 2, 2, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeCircle(screen, float32(g.Position.X), float32(g.Position.Y), 2, 2, color.RGBA{0, 0, 255, 255}, true)
		vector.StrokeRect(screen, float32(g.Width/2)-float32(float64(rover.Bounds().Dx()/2)*scale), float32(g.Position.Y)-float32(float64(rover.Bounds().Dy()/2)*scale), float32(float64(rover.Bounds().Dx())*scale), float32(float64(rover.Bounds().Dy())*scale), 2, color.RGBA{255, 0, 0, 255}, true)
	}
}

func (g *GameData) Update() error {
	platform.Ready()
	if !platform.Browser() && inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.Debug = !g.Debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		if !ebiten.IsFullscreen() {
			g.Width = g.oldwidth
			g.Height = g.oldheight
			ebiten.SetWindowSize(g.Width, g.Height)
		} else {
			g.oldwidth = g.Width
			g.oldheight = g.Height
			g.Width, g.Height = ebiten.Monitor().Size()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.GameSpeed = (min(g.GameSpeed+5, maxSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.GameSpeed = (max(g.GameSpeed-5, minSpeed))
	}

	left, right, jump := platform.Input()
	g.State.Input(left || ebiten.IsKeyPressed(ebiten.KeyArrowLeft), right || ebiten.IsKeyPressed(ebiten.KeyArrowRight), jump || inpututil.IsKeyJustPressed(ebiten.KeySpace))
	g.State.Update(ebiten.TPS())
	return nil
}

func (g *GameData) Layout(width, height int) (int, int) {
	if platform.Browser() {
		return g.Width, g.Height
	}
	g.Width = width
	g.Height = height
	return g.Width, g.Height
}
