package paused

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/go-allegro/allegro"
)

var (
	_screen *allegro.Bitmap
)

func Register() {
	allegory.DefState("playing/paused").
		Init(Init).
		Render(Render).
		HandleEvent(HandleEvent).
		Cleanup(Cleanup)
}

func Init() {
	var err error

	allegory.Debug("Paused.")
	_screen, err = allegory.Display().Backbuffer().Clone()
	if err != nil {
		allegory.Errorf("failed to capture screen image: %s", err.Error())
	}
}

func Render(delta float32) {
	if _screen != nil {
		_screen.Draw(0, 0, 0)
	}
}

func HandleEvent(event interface{}) bool {
	if event, ok := event.(allegro.KeyDownEvent); ok {
		if event.KeyCode() == allegro.KEY_ENTER {
			allegory.PopState()
			return true
		}
	}

	return false
}

func Cleanup() {
	allegory.Debug("Unpaused.")
}
