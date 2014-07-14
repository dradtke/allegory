package gopher

import (
	"runtime"
)

// GameState is an interface to the game's current state. Only one game
// state is active at any point in time, and states can be changed
// by using either NewState() or NewStateNow().
type GameState interface {
	// Perform initialization; this method is called once, when the
	// state becomes the game state.
	InitState()

	// Render; this is called (ideally) once per frame, with a delta
	// value calculated based on lag.
	RenderState(delta float32)

	// Perform cleanup; this method is called once, when the state
	// has been replaced by another one.
	CleanupState()
}

// NewState() waits for all processes to finish without
// blocking the current goroutine, then changes the game state.
func NewState(state GameState, views ...View) {
	go func() {
		for _processes.Len() > 0 {
			runtime.Gosched()
		}
		setState(state, views...)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(state GameState, views ...View) {
	NotifyAllProcesses(quit{})
	for _processes.Len() > 0 {
		runtime.Gosched()
	}
	setState(state, views...)
}

type BaseState struct{}

func (s *BaseState) InitState()                {}
func (s *BaseState) RenderState(delta float32) {}
func (s *BaseState) CleanupState()             {}

var _ GameState = (*BaseState)(nil)

func setState(state GameState, views ...View) {
	if _state != nil {
		_state.CleanupState()
	}
	_state = state
	_state.InitState()
	_views.Init()
	if views != nil {
		for _, v := range views {
			_views.PushBack(v)
		}
	}
}
