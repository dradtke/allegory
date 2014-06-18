package gopher

import (
	"runtime"
)

type GameState interface {
	Enter()
	Render()
	Leave()
}

func newState(s GameState) {
	if state != nil {
		state.Leave()
	}
	state = s
	state.Enter()
}

// NewState() waits for all processes
// to finish without blocking the current goroutine,
// then changes the game state.
func NewState(s GameState) {
	go func() {
		for processes.Len() > 0 {
			runtime.Gosched()
		}
		NewState(s)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(s GameState) {
    Broadcast(quit{})
    for processes.Len() > 0 {
        runtime.Gosched()
    }
    newState(s)
}

type BlankState struct{}

func (s *BlankState) Enter()  {}
func (s *BlankState) Render() {}
func (s *BlankState) Leave()  {}
