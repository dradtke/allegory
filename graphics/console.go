package graphics

import (
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
    "github.com/dradtke/gopher/config"
)

const (
	PROMPT = "> "
)

type Line struct {
    Text string
    Color al.Color
}

func RenderConsole(lines []Line, cmd string, is_blunk bool) {
    f := BuiltinFont()
	dw, dh := config.DisplaySize()

    prim.DrawFilledRoundedRectangle(
        prim.Point{X: 5, Y: float32(dh - 32 - ((f.LineHeight()+2) * len(lines)))},
        prim.Point{X: float32(dw - 5), Y: float32(dh - 6)},
        5, 5, al.MapRGBA(0, 0, 0, 120))

	for i, line := range lines {
		font.DrawText(f, line.Color, 10, float32(dh-(i+1)*(f.LineHeight()+2))-24, font.ALIGN_LEFT, line.Text)
	}

	font.DrawText(f, al.MapRGB(255, 255, 255), 10, float32((dh-10)-f.LineHeight()), font.ALIGN_LEFT, PROMPT+cmd)

	if is_blunk {
		x := 10 + f.TextWidth(PROMPT+cmd)
		prim.DrawLine(prim.Point{X: float32(x), Y: float32(dh - 10)}, prim.Point{X: float32(x + 10), Y: float32(dh - 10)}, al.MapRGB(255, 255, 255), 3)
	}
}
