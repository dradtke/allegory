package gopher

import (
    "github.com/dradtke/go-allegro/allegro/font"
)

var builtinFont *font.Font

func BuiltinFont() *font.Font {
    if builtinFont == nil {
        if f, err := font.Builtin(); err != nil {
            panic(err)
        } else {
            builtinFont = f
        }
    }
    return builtinFont
}
