package game

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"rover/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Vector2 struct {
	X, Y float64
}
type GameData struct {
	width, height       int
	oldwidth, oldheight int
	jumpTime            int
	tickCounter         int
	gameSpeed           int
	frame               int
	debug               bool
	direction           float64
	position            image.Point
	curVelocity         Vector2
	maxVelocity         Vector2
}

const (
	minSpeed    = 0
	maxSpeed    = 60
	frameCount  = 4
	frameOffset = 1
	gravity     = 2           //units/tick^2
	halfGravity = gravity / 2 //units/tick^2
	groundLevel = 75
)

func NewGame(width, height, gamespeed int, debug bool) (*GameData, error) {
	err := assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	if gamespeed < minSpeed {
		gamespeed = max(1, minSpeed)
	} else if gamespeed > maxSpeed {
		gamespeed = maxSpeed
	}

	ebiten.SetWindowTitle("Rover")
	ebiten.SetWindowSize(width, height)
	position := image.Point{X: width / 2, Y: 100}
	maxVelocity := Vector2{X: 15, Y: -15}
	return &GameData{width: width, height: height, position: position, maxVelocity: maxVelocity, gameSpeed: gamespeed, frame: 1, debug: debug, direction: 1}, nil
}

func (g *GameData) Draw(screen *ebiten.Image) {
	screen.Clear()

	background := assets.GetImage("backgroundH")
	bgScaleX := float64(g.width) / float64(background.Bounds().Dx())
	bgScaleY := float64(g.height) / float64(background.Bounds().Dy())
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(bgScaleX, bgScaleY)
	var drawPosition1 float64 = -float64(g.position.X % int(float64(g.width)))
	if g.position.X < 0 {
		drawPosition1 -= float64(g.width)
	}
	opts.GeoM.Translate(drawPosition1, 0)
	screen.DrawImage(background, opts)
	opts.GeoM.Translate(float64(g.width), 0)
	screen.DrawImage(background, opts)

	rover := assets.GetImage(fmt.Sprintf("rover%d", g.frame%frameCount+frameOffset))
	opts = &ebiten.DrawImageOptions{}
	scale := 0.5
	opts.GeoM.Scale(scale*g.direction, scale)
	opts.GeoM.Translate(float64(g.width/2)-g.direction*float64(rover.Bounds().Dx()/2)*scale, float64(g.position.Y)-float64(rover.Bounds().Dy()/2)*scale)
	screen.DrawImage(rover, opts)

	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("jump time: %d\njumping %v\ngame speed %d", g.jumpTime, g.isJumping(), g.gameSpeed))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("x, y %d, %d\nxscale, yscale %f, %f", g.position.X, g.position.Y, bgScaleX, bgScaleY), g.width/2, 0)

		vector.StrokeCircle(screen, float32(g.width/2), float32(g.position.Y), 2, 2, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeCircle(screen, float32(g.position.X), float32(g.position.Y), 2, 2, color.RGBA{0, 0, 255, 255}, true)
		vector.StrokeRect(screen, float32(g.width/2)-float32(float64(rover.Bounds().Dx()/2)*scale), float32(g.position.Y)-float32(float64(rover.Bounds().Dy()/2)*scale), float32(float64(rover.Bounds().Dx())*scale), float32(float64(rover.Bounds().Dy())*scale), 2, color.RGBA{255, 0, 0, 255}, true)
	}
}

func (g *GameData) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		if !ebiten.IsFullscreen() {
			g.width = g.oldwidth
			g.height = g.oldheight
			ebiten.SetWindowSize(g.width, g.height)
		} else {
			g.oldwidth = g.width
			g.oldheight = g.height
			g.width, g.height = ebiten.Monitor().Size()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		g.gameSpeed = (min(g.gameSpeed+5, maxSpeed))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		g.gameSpeed = (max(g.gameSpeed-5, minSpeed))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startJump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.curVelocity.X = -g.maxVelocity.X
		g.direction = -1

	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.curVelocity.X = g.maxVelocity.X
		g.direction = 1
	} else {
		g.curVelocity.X = 0
	}

	if g.gameSpeed != 0 && float32(g.tickCounter) > float32(ebiten.TPS())/float32(g.gameSpeed) {
		g.tickCounter = 0
		g.position.X += int(g.curVelocity.X)
		if g.isMoving() {
			g.frame++
		}
		if g.isJumping() || g.isAirborne() {
			g.position.Y, g.curVelocity.Y = g.jumpPosition(float64(g.position.Y), g.curVelocity.Y, g.jumpTime)
			g.jumpTime++
			if g.position.Y >= g.groundLevel() {
				g.stopJump()
			}
		} else if g.position.Y > g.groundLevel() {
			g.position.Y = g.groundLevel()
		}
	} else {
		g.tickCounter++
	}

	return nil
}

func (g *GameData) jumpPosition(y0, v0 float64, time int) (int, float64) {
	//y(t) = y0 + v0*t + (1/2)*g*t^2
	//v(t) = v0 + g*t
	yt := y0 + v0*float64(time) + halfGravity*float64(time*time)
	vt := v0 + gravity*float64(time)
	if g.debug {
		fmt.Printf("y0 = %f, v0 = %f, yt = %f, vt = %f, time = %d\n", y0, v0, yt, vt, time)
	}
	return int(math.Round(yt)), vt
}

func (g *GameData) startJump() {
	if g.debug {
		fmt.Println("starting jump")
	}
	g.jumpTime = 1
	g.curVelocity.Y += g.maxVelocity.Y
}

func (g *GameData) stopJump() {
	if g.debug {
		fmt.Println("stopping jump")
	}
	g.jumpTime = 0
	g.curVelocity.Y = 0
	g.position.Y = g.groundLevel()
}
func (g *GameData) isJumping() bool {
	return g.jumpTime > 0
}

func (g *GameData) isAirborne() bool {
	return g.position.Y < g.groundLevel()
}

func (g *GameData) isMoving() bool {
	return g.curVelocity.X != 0 || g.curVelocity.Y != 0
}

func (g *GameData) Layout(width, height int) (int, int) {
	g.width = width
	g.height = height
	return g.width, g.height
}

func (g *GameData) groundLevel() int {
	return g.height - groundLevel
}
