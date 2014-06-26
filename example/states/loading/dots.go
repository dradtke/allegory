package loading

import (
	"github.com/dradtke/gopher"
	"github.com/dradtke/gopher/bus"
	"github.com/dradtke/gopher/example/events"
)

type LoadingDotAnimation struct{
    timer int
    DotDelay int
}

func (p *LoadingDotAnimation) InitProcess() {}

func (p *LoadingDotAnimation) HandleMessage(msg interface{}) error {
    return nil
}

func (p *LoadingDotAnimation) Tick() (bool, error) {
    p.timer++
    if p.timer >= p.DotDelay {
        state := gopher.State().(*LoadingState)
        if state.dots == "..." {
            state.dots = ""
        } else {
            state.dots += "."
        }
        bus.Signal(events.DotNotifyEvent, state.dots)
        p.timer = 0
    }
    return true, nil
}

func (p *LoadingDotAnimation) CleanupProcess() {}
