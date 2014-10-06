package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/go-allegro/allegro"
)

const (
	PAUSE_MSG = " [Paused]"
)

type PausedState struct {
	allegory.BaseState
	screenshot *allegro.Bitmap
}

func (s *PausedState) InitState() {
	allegory.Display().SetWindowTitle(config.WindowTitle() + PAUSE_MSG)
	allegory.AddView(new(PausedView))
}

func (s *PausedState) RenderState(delta float32) {
	s.screenshot.Draw(0, 0, allegro.FLIP_NONE)
}

func (s *PausedState) CleanupState() {
	allegory.Display().SetWindowTitle(config.WindowTitle())
	s.screenshot.Destroy()
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
