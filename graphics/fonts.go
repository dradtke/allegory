package graphics

import (
	"github.com/dradtke/go-allegro/allegro/font"
)

var builtin *font.Font

func BuiltinFont() *font.Font {
    if builtin == nil {
        var err error
        builtin, err = font.Builtin()
        if err != nil {
            panic(err)
        }
    }
    return builtin
}
