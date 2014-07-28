// +build ignore

package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/bus"
	"github.com/dradtke/go-allegro/allegro"
	"math"
)

var (
	FULL_ROTATION = float32(2 * math.Pi)
)

const (
	_ uint = iota
	SPIN_SQUARE
)

/* -- Square -- */

type Square struct {
	allegory.BaseActor
	Bitmap *allegro.Bitmap
	Color  allegro.Color
	Angle  float32
}

func (s *Square) InitActor() {
	bus.AddListener(SPIN_SQUARE, s.Spin)
	s.Bitmap = allegro.CreateBitmap(100, 100).AsTarget(func() {
		allegro.ClearToColor(s.Color)
	})
}

func (s *Square) Spin(id allegory.ActorId, theta float32) {
	if s.Id != id {
		return
	}
	s.Angle += theta
}

func (s *Square) RenderActor(delta float32) {
	bw, bh := s.Bitmap.Width(), s.Bitmap.Height()
	s.Bitmap.DrawRotated(float32(bw/2), float32(bh/2), s.X, s.Y, s.Angle, allegro.FLIP_NONE)
}

/* -- State -- */

type MainState struct {
	allegory.BaseState
	Square *Square
}

func (s *MainState) InitState() {
	s.Square = new(Square)
	s.Square.Color = allegro.MapRGB(255, 0, 0)
	s.Square.X, s.Square.Y = 200, 100
	allegory.AddActor(0, s.Square)

	blueSquare := new(Square)
	blueSquare.Color = allegro.MapRGB(0, 0, 255)
	blueSquare.X, blueSquare.Y = 400, 200
	allegory.AddActor(1, blueSquare)
}

/* -- View -- */

type HumanView struct {
	allegory.BaseView
	SpinSpeed           float32
	LeftDown, RightDown bool
}

func (v *HumanView) InitView() {
	v.SpinSpeed = (FULL_ROTATION / 240)
}

func (v *HumanView) OnKeyEvent(code allegro.KeyCode, down bool) bool {
	switch code {
	case allegro.KEY_LEFT:
		v.LeftDown = down
		return true
	case allegro.KEY_RIGHT:
		v.RightDown = down
		return true
	}
	return false
}

func (v *HumanView) HandleEvent(event interface{}) bool {
	switch e := event.(type) {
	case allegro.KeyDownEvent:
		return v.OnKeyEvent(e.KeyCode(), true)
	case allegro.KeyUpEvent:
		return v.OnKeyEvent(e.KeyCode(), false)
	}
	return false
}

func (v *HumanView) UpdateView() {
	if v.LeftDown && !v.RightDown {
		bus.Signal(SPIN_SQUARE, allegory.ActorId(1), -v.SpinSpeed)
	} else if v.RightDown && !v.LeftDown {
		bus.Signal(SPIN_SQUARE, allegory.ActorId(1), v.SpinSpeed)
	}
}

/* -- Main -- */

func main() {
	allegory.Init(new(MainState), new(HumanView))
	defer allegory.Cleanup()
	allegory.Loop()
}
