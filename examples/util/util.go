package util

import (
	"github.com/dradtke/go-allegro/allegro"
)

func DirToFlags(dir int8) allegro.DrawFlags {
	if dir < 0 {
		return allegro.FLIP_HORIZONTAL
	} else {
		return allegro.FLIP_NONE
	}
}
