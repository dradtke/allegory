package allegory

import (
	"runtime"
)

// GameState is an interface to the game's current state. Only one game
// state is active at any point in time, and states can be changed
// by using one of the NewState*() functions.
type GameState interface {
	// Perform initialization; this method is called once, when the
	// state becomes the game state.
	InitGameState()

	// Called once per frame to perform any necessary updates.
	UpdateGameState()

	// Perform cleanup; this method is called once, when the state
	// has been replaced by another one.
	CleanupGameState()
}

type RenderableGameState interface {
	GameState

	// Render; this is called (ideally) once per frame, with a delta
	// value calculated based on lag.
	RenderGameState(delta float32)
}

// NewState() changes the state, regardless of the status of currently
// running processes.
func NewState(state GameState) {
	if _state != nil {
		_state.CleanupGameState()
	}
	for _, view := range _views {
		view.CleanupView()
	}
	_views = make([]View, 0)
	for _, actor := range _actors {
		actor.CleanupActor()
	}
	_actors = make([]Actor, 0)
	_actorLayers = make(map[uint][]Actor)

	runtime.GC()

	_state = state
	_state.InitGameState()
}

// NewStateWait() waits for all processes to finish without
// blocking the current goroutine, then changes the game state.
func NewStateWait(state GameState) {
	go func() {
		for len(_processes) > 0 {
			runtime.Gosched()
		}
		NewState(state)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(state GameState) {
	NotifyAllProcesses(&quit{})
	for len(_processes) > 0 {
		runtime.Gosched()
	}
	NewState(state)
}

type BaseGameState struct{}

func (s *BaseGameState) InitGameState()    {}
func (s *BaseGameState) UpdateGameState()  {}
func (s *BaseGameState) CleanupGameState() {}

var _ GameState = (*BaseGameState)(nil)
