// Package simulation contains Rover motion independently of graphics and keyboard input.
package simulation

import (
	"fmt"
	"image"
	"math"
)

type Vector2 struct{ X, Y float64 }

type Config struct{ Width, Height, Speed int }

func (c Config) Validate() error {
	if c.Width < 1 || c.Height < 1 || c.Width > 8192 || c.Height > 8192 {
		return fmt.Errorf("viewport dimensions must be 1..8192 pixels")
	}
	if c.Speed < 0 || c.Speed > 60 {
		return fmt.Errorf("speed must be 0..60")
	}
	return nil
}

// State is owned by one game. The renderer reads it; input and physics mutate it.
type State struct {
	Width, Height, JumpTime, TickCounter, GameSpeed, Frame int
	Debug                                                  bool
	Direction                                              float64
	Position                                               image.Point
	CurVelocity, MaxVelocity                               Vector2
}

const (
	gravity     = 2
	halfGravity = gravity / 2
	groundLevel = 75
)

func New(c Config, debug bool) (*State, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &State{Width: c.Width, Height: c.Height, Position: image.Point{X: c.Width / 2, Y: 100}, MaxVelocity: Vector2{X: 15, Y: -15}, GameSpeed: c.Speed, Frame: 1, Direction: 1, Debug: debug}, nil
}

// Input preserves left-key priority and the original repeated midair jump boosts.
func (g *State) Input(left, right, jump bool) {
	if jump {
		g.startJump()
	}
	if left {
		g.CurVelocity.X = -g.MaxVelocity.X
		g.Direction = -1

	} else if right {
		g.CurVelocity.X = g.MaxVelocity.X
		g.Direction = 1
	} else {
		g.CurVelocity.X = 0
	}

}

// Update preserves the original strict tick gate and accumulated-time jump formula.
func (g *State) Update(tps int) {
	if g.GameSpeed != 0 && float32(g.TickCounter) > float32(tps)/float32(g.GameSpeed) {
		g.TickCounter = 0
		g.Position.X += int(g.CurVelocity.X)
		if g.isMoving() {
			g.Frame++
		}
		if g.isJumping() || g.isAirborne() {
			g.Position.Y, g.CurVelocity.Y = g.jumpPosition(float64(g.Position.Y), g.CurVelocity.Y, g.JumpTime)
			g.JumpTime++
			if g.Position.Y >= g.groundLevel() {
				g.stopJump()
			}
		} else if g.Position.Y > g.groundLevel() {
			g.Position.Y = g.groundLevel()
		}
	} else {
		g.TickCounter++
	}

}

func (g *State) jumpPosition(y0, v0 float64, time int) (int, float64) {
	//y(t) = y0 + v0*t + (1/2)*g*t^2
	//v(t) = v0 + g*t
	yt := y0 + v0*float64(time) + halfGravity*float64(time*time)
	vt := v0 + gravity*float64(time)
	if g.Debug {
		fmt.Printf("y0 = %f, v0 = %f, yt = %f, vt = %f, time = %d\n", y0, v0, yt, vt, time)
	}
	return int(math.Round(yt)), vt
}

func (g *State) startJump() {
	if g.Debug {
		fmt.Println("starting jump")
	}
	g.JumpTime = 1
	g.CurVelocity.Y += g.MaxVelocity.Y
}

func (g *State) stopJump() {
	if g.Debug {
		fmt.Println("stopping jump")
	}
	g.JumpTime = 0
	g.CurVelocity.Y = 0
	g.Position.Y = g.groundLevel()
}
func (g *State) isJumping() bool {
	return g.JumpTime > 0
}

func (g *State) isAirborne() bool {
	return g.Position.Y < g.groundLevel()
}

func (g *State) isMoving() bool {
	return g.CurVelocity.X != 0 || g.CurVelocity.Y != 0
}

func (g *State) groundLevel() int {
	return g.Height - groundLevel
}

func (g *State) IsMoving() bool  { return g.isMoving() }
func (g *State) IsJumping() bool { return g.isJumping() }
