package keyboard

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
)

var (
	ctrl_mod  int8
	alt_mod   int8
	shift_mod int8
)

type Mod int

const (
	Ctrl Mod = iota
	Alt
	Shift
)

func update_mods(key allegro.KeyCode, delta int8) (is_mod bool) {
	is_mod = true
	switch {
	case key == allegro.KEY_LCTRL || key == allegro.KEY_RCTRL:
		ctrl_mod += delta
	case key == allegro.KEY_LALT || key == allegro.KEY_RALT:
		alt_mod += delta
	case key == allegro.KEY_LSHIFT || key == allegro.KEY_RSHIFT:
		shift_mod += delta
	default:
		is_mod = false
	}
	return
}

var pressed = make(map[allegro.KeyCode]bool)

func Down(key allegro.KeyCode) {
	if !update_mods(key, 1) {
		pressed[key] = true
	}
}

func Up(key allegro.KeyCode) {
	if !update_mods(key, -1) {
		delete(pressed, key)
	}
}

func matches(mods []Mod, keys []allegro.KeyCode) bool {
	var (
		needCtrl  bool
		needAlt   bool
		needShift bool
	)

	for m := range mods {
		switch m {
		case Ctrl:
			needCtrl = true
		case Alt:
			needAlt = true
		case Shift:
			needShift = true
		}
	}

	if needCtrl != ctrl_mod || needAlt != alt_mod || needShift != shift_mod {
		return
	}

	// we have the modifiers!
}
