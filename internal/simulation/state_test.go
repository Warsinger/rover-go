package simulation

import (
	"image"
	"testing"
)

func rover(t *testing.T) *State {
	t.Helper()
	s, err := New(Config{800, 600, 30}, false)
	if err != nil {
		t.Fatal(err)
	}
	return s
}
func step(s *State) { s.TickCounter = 3; s.Update(60) }
func TestInitialStateAndInput(t *testing.T) {
	s := rover(t)
	if s.Position != (image.Point{400, 100}) || s.Frame != 1 || s.Direction != 1 || s.MaxVelocity != (Vector2{15, -15}) {
		t.Fatalf("unexpected initial state: %+v", s)
	}
	tests := []struct {
		left, right         bool
		velocity, direction float64
	}{{false, true, 15, 1}, {true, true, -15, -1}, {false, false, 0, -1}, {true, false, -15, -1}}
	for _, tt := range tests {
		s.Input(tt.left, tt.right, false)
		if s.CurVelocity.X != tt.velocity || s.Direction != tt.direction {
			t.Fatalf("input %+v produced %+v", tt, s)
		}
	}
}
func TestGroundMovementAnimationAndPause(t *testing.T) {
	s := rover(t)
	s.Position.Y = 525
	s.Input(false, true, false)
	for range 3 {
		s.Update(60)
	}
	if s.Position.X != 400 || s.Frame != 1 {
		t.Fatal("advanced before strict threshold")
	}
	s.Update(60)
	if s.Position != (image.Point{415, 525}) || s.Frame != 2 {
		t.Fatalf("movement = %+v", s)
	}
	s.Input(true, false, false)
	step(s)
	if s.Position.X != 400 || s.Frame != 3 {
		t.Fatal("left movement/animation changed")
	}
	s.Input(false, false, false)
	step(s)
	if s.Frame != 3 || s.IsMoving() {
		t.Fatal("stationary rover animated")
	}
	s.GameSpeed = 0
	s.Input(false, true, true)
	position, frame := s.Position, s.Frame
	for range 100 {
		s.Update(60)
	}
	if s.Position != position || s.Frame != frame {
		t.Fatal("pause advanced motion")
	}
	if s.CurVelocity.Y != -15 || !s.IsJumping() {
		t.Fatal("pause no longer accepts queued jump")
	}
	s.GameSpeed = 30
	s.Update(60)
	if s.Position.X != 415 {
		t.Fatal("legacy accumulated pause timing changed")
	}
}
func TestJumpTrajectoryAndLanding(t *testing.T) {
	s := rover(t)
	s.Position.Y = 525
	s.Input(false, false, true)
	trajectory := []struct {
		y    int
		vy   float64
		jump int
	}{{511, -13, 2}, {489, -9, 3}, {471, -3, 4}, {475, 5, 5}, {525, 0, 0}}
	for i, want := range trajectory {
		step(s)
		if s.Position.Y != want.y || s.CurVelocity.Y != want.vy || s.JumpTime != want.jump {
			t.Fatalf("step %d: got y=%d vy=%g jump=%d, want %+v", i, s.Position.Y, s.CurVelocity.Y, s.JumpTime, want)
		}
	}
	if s.IsJumping() || s.IsMoving() {
		t.Fatal("landing retained motion")
	}
}
func TestMidairBoost(t *testing.T) {
	s := rover(t)
	s.Position.Y = 525
	s.Input(false, false, true)
	step(s)
	s.Input(false, false, true)
	if s.JumpTime != 1 || s.CurVelocity.Y != -28 {
		t.Fatal("repeated jump did not reset time and add boost")
	}
	step(s)
	if s.Position.Y != 484 || s.CurVelocity.Y != -26 {
		t.Fatalf("boost trajectory changed: %+v", s)
	}
}
func TestFallAndResizeGroundClamp(t *testing.T) {
	s := rover(t)
	for _, want := range []int{100, 101, 109} {
		step(s)
		if s.Position.Y != want {
			t.Fatalf("fall y=%d want=%d", s.Position.Y, want)
		}
	}
	for range 100 {
		step(s)
	}
	if s.Position.Y != 525 || s.JumpTime != 0 || s.CurVelocity.Y != 0 {
		t.Fatal("fall did not land")
	}
	s.Height = 300
	step(s)
	if s.Position.Y != 225 {
		t.Fatal("resize did not clamp ground")
	}
	s.Height = 600
	for range 100 {
		step(s)
	}
	if s.Position.Y != 525 {
		t.Fatal("resize upward did not fall to new ground")
	}
}
func TestConfigAndIndependentGames(t *testing.T) {
	for _, c := range []Config{{0, 600, 30}, {800, -1, 30}, {8193, 600, 30}, {800, 8193, 30}, {800, 600, -1}, {800, 600, 61}} {
		if _, err := New(c, false); err == nil {
			t.Errorf("accepted %+v", c)
		}
	}
	for _, c := range []Config{{1, 1, 0}, {8192, 8192, 60}} {
		if _, err := New(c, false); err != nil {
			t.Errorf("rejected %+v: %v", c, err)
		}
	}
	a, b := rover(t), rover(t)
	a.Input(true, false, true)
	step(a)
	if b.Position != (image.Point{400, 100}) || b.IsJumping() || b.TickCounter != 0 {
		t.Fatal("games share state")
	}
	// Exercise optional diagnostics without requiring a window or audio device.
	a.Debug = true
	a.Input(false, false, true)
	step(a)
	a.Position.Y = 1000
	step(a)
}
