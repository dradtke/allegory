package actors

import (
	"github.com/dradtke/go-allegro/allegro"
)

func dirToFlags(dir int8) allegro.DrawFlags {
	if dir < 0 {
		return allegro.FLIP_HORIZONTAL
	}
	return allegro.FLIP_NONE
}
