package graphics

import (
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/primitives"
)

const (
	PROMPT = "> "
)

type Line struct {
	Text  string
	Color allegro.Color
}

func RenderConsole(lines []Line, cmd string, is_blunk bool) {
	f := BuiltinFont()
	dw, dh := config.DisplaySize()

	primitives.DrawFilledRoundedRectangle(
		primitives.Point{X: 5, Y: float32(dh - 32 - ((f.LineHeight() + 2) * len(lines)))},
		primitives.Point{X: float32(dw - 5), Y: float32(dh - 6)},
		5, 5, allegro.MapRGBA(0, 0, 0, 120))

	for i, line := range lines {
		font.DrawText(f, line.Color, 10, float32(dh-(i+1)*(f.LineHeight()+2))-24, font.ALIGN_LEFT, line.Text)
	}

	font.DrawText(f, allegro.MapRGB(255, 255, 255), 10, float32((dh-10)-f.LineHeight()), font.ALIGN_LEFT, PROMPT+cmd)

	if is_blunk {
		x := 10 + f.TextWidth(PROMPT+cmd)
		primitives.DrawLine(
			primitives.Point{X: float32(x), Y: float32(dh - 10)},
			primitives.Point{X: float32(x + 10), Y: float32(dh - 10)},
			allegro.MapRGB(255, 255, 255),
			3,
		)
	}
}
