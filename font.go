package gopher

import (
    "github.com/dradtke/go-allegro/allegro/font"
)

var builtinFont *font.Font

func Font() *font.Font {
    if builtinFont == nil {
        var err error
        builtinFont, err = font.Builtin()
        if err != nil {
            panic(err)
        }
    }
    return builtinFont
}
