package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/go-allegro/allegro"
)

type PausedState struct {
	allegory.BaseState
	snapshot *allegro.Bitmap
}

func (s *PausedState) InitState() {
	allegory.AddView(new(PausedView))
}

func (s *PausedState) RenderState(delta float32) {
	s.snapshot.Draw(0, 0, allegro.FLIP_NONE)
}

type PausedView struct {
	allegory.BaseView
}

func (s *PausedView) HandleEvent(event interface{}) bool {
	_playingState.heroView.BaseKeyView.HandleEvent(event)

	switch e := event.(type) {
	case allegro.KeyDownEvent:
		if e.KeyCode() == allegro.KEY_ENTER {
			allegory.NewStateNow(_playingState)
			return true
		}
	}

	return false
}
