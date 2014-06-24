package gopher

import (
	"runtime"
)

type GameState interface {
	Enter()
	Render()
	Leave()
}

func newState(state GameState, views ...View) {
	if _state != nil {
		_state.Leave()
	}
	_state = state
	_state.Enter()
    _views.Init()
    if views != nil {
        for v := range views {
            _views.PushBack(v)
        }
    }
}

// NewState() waits for all processes to finish without
// blocking the current goroutine, then changes the game state.
func NewState(state GameState, views ...View) {
	go func() {
		for _processes.Len() > 0 {
			runtime.Gosched()
		}
		NewState(state, views...)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(state GameState, views ...View) {
    NotifyAll(quit{})
    for _processes.Len() > 0 {
        runtime.Gosched()
    }
    newState(state, views...)
}

type BlankState struct{}

func (s *BlankState) Enter()  {}
func (s *BlankState) Render() {}
func (s *BlankState) Leave()  {}
