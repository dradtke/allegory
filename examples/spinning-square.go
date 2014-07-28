// +build ignore

package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/go-allegro/allegro"
	"math"
)

var (
	FULL_ROTATION = float32(2 * math.Pi)
)

type SpinningState struct {
	allegory.BaseState
	Square *allegro.Bitmap
	Angle  float32
}

func (s *SpinningState) InitState() {
	s.Square = allegro.CreateBitmap(100, 100).AsTarget(func() {
		allegro.ClearToColor(allegro.MapRGB(255, 0, 0))
	})
}

func (s *SpinningState) UpdateState() {
	s.Angle += (FULL_ROTATION / 240) // at 60fps, this takes about 4 seconds to do a full rotation
	for s.Angle > FULL_ROTATION {
		s.Angle -= FULL_ROTATION
	}
}

func (s *SpinningState) RenderState(delta float32) {
	dw, dh := config.DisplaySize()
	s.Square.DrawRotated(50, 50, float32(dw/2), float32(dh/2), s.Angle, allegro.FLIP_NONE)
}

func main() {
	allegory.Init(new(SpinningState))
	defer allegory.Cleanup()
	allegory.Loop()
}
