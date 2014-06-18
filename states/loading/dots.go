package loading

import (
	"github.com/dradtke/gopher"
)

type LoadingDotAnimation struct{
    timer int
    DotDelay int
}

func (p *LoadingDotAnimation) HandleMessage(msg interface{}) {
}

func (p *LoadingDotAnimation) HandleEvent(ev interface{}) bool {
    return false
}

func (p *LoadingDotAnimation) Tick(delta float32) (bool, error) {
    p.timer++
    if p.timer >= p.DotDelay {
        state := gopher.State().(*LoadingState)
        if state.dots == "..." {
            state.dots = ""
        } else {
            state.dots += "."
        }
        p.timer = 0
    }
    return true, nil
}
