package loading

import (
    al "github.com/dradtke/go-allegro/allegro"
    "github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/gopher"
)

type LoadingState struct{
    dots string
}

func (s *LoadingState) Enter() {
    font.Install()
    gopher.RunProcess(&LoadingDotAnimation{DotDelay: 30})
}

func (s *LoadingState) Render() {
    font.DrawText(gopher.Font(), al.MapRGB(255, 255, 255), 10, 10, font.ALIGN_LEFT, "Loading" + s.dots)
}

func (s *LoadingState) Leave() {}
