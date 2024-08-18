package game

import (
	"fmt"
	"image"
	"math"
	"rover/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Vector2 struct {
	X, Y float64
}
type GameData struct {
	width, height int
	position      image.Point
	curVelocity   Vector2
	maxVelocity   Vector2
	jumpTime      int
	groundY       int
	tickCounter   int
	gameSpeed     int
	frame         int
	debug         bool
}

const (
	minSpeed    = 0
	maxSpeed    = 60
	frameCount  = 4
	frameOffset = 1
	gravity     = 2           //units/tick^2
	halfGravity = gravity / 2 //units/tick^2
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
	groundY := height - 50
	position := image.Point{X: width / 2, Y: groundY}
	maxVelocity := Vector2{X: 15, Y: -15}
	return &GameData{width: width, height: height, position: position, maxVelocity: maxVelocity, groundY: groundY, gameSpeed: gamespeed, frame: 1, debug: debug}, nil
}

func (g *GameData) Draw(screen *ebiten.Image) {
	screen.Clear()

	background := assets.GetImage("backgroundH")
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(-float64(g.position.X%g.width), 0)
	screen.DrawImage(background, opts)
	opts.GeoM.Translate(float64(g.width), 0)
	screen.DrawImage(background, opts)

	rover := assets.GetImage(fmt.Sprintf("rover%d", g.frame%frameCount+frameOffset))
	opts = &ebiten.DrawImageOptions{}
	scale := 0.5
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Translate(float64(g.width/2)-float64(rover.Bounds().Dx()/2)*scale, float64(g.position.Y)-float64(rover.Bounds().Dy()/2)*scale)
	screen.DrawImage(rover, opts)

	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("jump time: %d\njumping %v", g.jumpTime, g.isJumping()))
	}
}

func (g *GameData) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startJump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.curVelocity.X = -g.maxVelocity.X

	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.curVelocity.X = g.maxVelocity.X
	} else {
		g.curVelocity.X = 0
	}

	if g.gameSpeed != 0 && float32(g.tickCounter) > float32(ebiten.TPS())/float32(g.gameSpeed) {
		g.tickCounter = 0
		g.position.X += int(g.curVelocity.X)
		if g.isMoving() {
			g.frame++
		}
		if g.isJumping() {
			g.position.Y, g.curVelocity.Y = g.jumpPosition(float64(g.position.Y), g.curVelocity.Y, g.jumpTime)
			g.jumpTime++
			if g.position.Y >= g.groundY {
				g.stopJump()
			}
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
	g.position.Y = g.groundY
}
func (g *GameData) isJumping() bool {
	return g.jumpTime > 0
}

func (g *GameData) isMoving() bool {
	return g.curVelocity.X != 0 || g.curVelocity.Y != 0
}

func (g *GameData) Layout(width, height int) (int, int) {
	return g.width, g.height
}
