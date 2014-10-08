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
	allegory.BaseGameState
	screenshot *allegro.Bitmap
	overlay    *allegro.Bitmap
}

func (s *PausedState) InitGameState() {
	allegory.Debug("Paused.")
	allegory.NotifyWhere(&allegory.PauseAnimation{}, func(p allegory.Process) bool {
		_, ok := p.(*allegory.AnimationProcess)
		return ok
	})

	w, h := config.DisplaySize()
	s.screenshot, _ = allegory.Display().Backbuffer().Clone()
	s.overlay = allegro.CreateBitmap(w, h).AsTarget(func() {
		allegro.ClearToColor(allegro.MapRGB(255, 255, 255))
	})
	allegory.Display().SetWindowTitle(config.WindowTitle() + PAUSE_MSG)
	allegory.AddView(new(PausedView))
}

func (s *PausedState) RenderGameState(delta float32) {
	s.screenshot.Draw(0, 0, allegro.FLIP_NONE)
	s.overlay.DrawTinted(allegro.MapRGBAf(1, 1, 1, 0.33), 0, 0, allegro.FLIP_NONE)
}

func (s *PausedState) CleanupGameState() {
	allegory.Debug("Resume.")
	allegory.NotifyWhere(&allegory.ResumeAnimation{}, func(p allegory.Process) bool {
		_, ok := p.(*allegory.AnimationProcess)
		return ok
	})

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
			allegory.NewState(_playingState)
			return true
		}
	}

	return false
}
