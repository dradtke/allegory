package loading

import (
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/gopher"
)

type LoadingState struct {
	dots string
}

func (s *LoadingState) InitState() {
	font.Install()
	gopher.RunProcess(&LoadingDotAnimation{DotDelay: 30})
}

func (s *LoadingState) RenderState(delta float32) {
	font.DrawText(gopher.Font(), allegro.MapRGB(255, 255, 255), 10, 10, font.ALIGN_LEFT, "Loading"+s.dots)
}

func (s *LoadingState) CleanupState() {}
