package paused

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
)

var (
	_inited bool

	_screen    *allegro.Bitmap
	_overlay   *allegro.Bitmap

	_font      *font.Font
	_fontColor allegro.Color

	_dw int
	_dh int
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

	if _inited {
		return
	}

	d := allegory.Display()
	_dw, _dh = d.Width(), d.Height()
	_overlay = allegro.CreateBitmap(_dw, _dh).AsTarget(func() {
		prim.DrawFilledRectangle(prim.Point{0, 0}, prim.Point{float32(_dw), float32(_dh)},
			allegro.MapRGBA(0, 0, 0, 120))
	})

	_font, err = font.Builtin()
	if err != nil {
		allegory.Errorf("failed to load builtin font: %s", err.Error())
	}

	_fontColor = allegro.MapRGB(0xFF, 0xFF, 0xFF)

	_inited = true
}

func Render(delta float32) {
	if _screen != nil {
		_screen.Draw(0, 0, 0)
	}
	_overlay.Draw(0, 0, 0)
	if _font != nil {
		font.DrawText(_font, _fontColor, float32(_dw/2), 100, font.ALIGN_CENTRE, "Paused.")
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
