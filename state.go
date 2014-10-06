package allegory

import (
	"runtime"
)

// State is an interface to the game's current state. Only one game
// state is active at any point in time, and states can be changed
// by using either NewState() or NewStateNow().
type State interface {
	// Perform initialization; this method is called once, when the
	// state becomes the game state.
	InitState()

	// Called once per frame to perform any necessary updates.
	UpdateState()

	// Perform cleanup; this method is called once, when the state
	// has been replaced by another one.
	CleanupState()
}

type RenderableState interface {
	State

	// Render; this is called (ideally) once per frame, with a delta
	// value calculated based on lag.
	RenderState(delta float32)
}

// NewState() waits for all processes to finish without
// blocking the current goroutine, then changes the game state.
func NewState(state State) {
	go func() {
		for len(_processes) > 0 {
			runtime.Gosched()
		}
		setState(state)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(state State) {
	NotifyAllProcesses(&quit{})
	for len(_processes) > 0 {
		runtime.Gosched()
	}
	setState(state)
}

type BaseState struct{}

func (s *BaseState) InitState()    {}
func (s *BaseState) UpdateState()  {}
func (s *BaseState) CleanupState() {}

var _ State = (*BaseState)(nil)

func setState(state State) {
	if _state != nil {
		_state.CleanupState()
	}
	for _, view := range _views {
		view.CleanupView()
	}
	_views = make([]View, 0)
	for _, actor := range _actors {
		actor.CleanupActor()
	}
	_actors = make(map[ActorId]Actor)
	runtime.GC()

	_state = state
	_state.InitState()
}
