package allegory

import (
	"container/list"
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

	// Called when a state is pushed over this one on the stack.
	// Processes are automatically paused, but this can be used
	// to take care of other tasks that are process-independent.
	OnPause()

	// Called when the state overriding this one was popped off the
	// stack. Useful for validating the game state when using
	// an event-based approach.
	OnResume()

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
// TODO: rewrite this so that it's not in terms of PopState() and PushState();
// this causes too many calls to OnPause() and OnResume().
func NewState(state GameState) {
	PopState()
	PushState(state)
}

// Push a new state to the top of the stack.
func PushState(state GameState) {
	_state.Push(state)
}

func PopState() GameState {
	return _state.Pop()
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
func (s *BaseGameState) OnPause()          {}
func (s *BaseGameState) OnResume()         {}
func (s *BaseGameState) CleanupGameState() {}

var _ GameState = (*BaseGameState)(nil)

/* -- stateStack -- */

type stateStack struct {
	stack *list.List
}

func (s *stateStack) Empty() bool {
	return s.stack.Len() == 0
}

func (s *stateStack) Current() GameState {
	front := s.stack.Front()
	if front == nil {
		return nil
	}
	return front.Value.(GameState)
}

func (s *stateStack) Push(state GameState) {
	cur := s.Current()
	if cur != nil {
		cur.OnPause()
	}

	s.stack.PushFront(state)

	if state != nil {
		_processes[state] = make([]Process, 0)
		_views[state] = make([]View, 0)
		_actors[state] = make([]Actor, 0)
		_actorLayers[state] = make(map[uint][]Actor)
		state.InitGameState()
	}
}

func (s *stateStack) Pop() GameState {
	oldState := s.stack.Remove(s.stack.Front()).(GameState)

	if oldState != nil {
		oldState.CleanupGameState()

		if views, ok := _views[oldState]; ok {
			for _, view := range views {
				view.CleanupView()
			}
			delete(_views, oldState)
		}

		if actors, ok := _actors[oldState]; ok {
			for _, actor := range actors {
				actor.CleanupActor()
			}
			delete(_actors, oldState)
			delete(_actorLayers, oldState)
		}

		runtime.GC()
	}

	if cur := s.Current(); cur != nil {
		cur.OnResume()
	}

	return oldState
}

func (s *stateStack) Update() {
	cur := s.Current()
	if cur != nil {
		cur.UpdateGameState()
	}
}

func (s *stateStack) Render(delta float32) {
	cur := s.Current()
	if cur != nil {
		if cur, ok := cur.(RenderableGameState); ok {
			cur.RenderGameState(delta)
		}
	}
}

func (s *stateStack) Processes() []Process {
	if processes, ok := _processes[s.Current()]; ok && processes != nil {
		return processes
	}
	return make([]Process, 0)
}

func (s *stateStack) Views() []View {
	if views, ok := _views[s.Current()]; ok && views != nil {
		return views
	}
	return make([]View, 0)
}

func (s *stateStack) Actors() []Actor {
	if actors, ok := _actors[s.Current()]; ok && actors != nil {
		return actors
	}
	return make([]Actor, 0)
}

func (s *stateStack) ActorLayers() map[uint][]Actor {
	if layers, ok := _actorLayers[s.Current()]; ok && layers != nil {
		return layers
	}
	return make(map[uint][]Actor)
}
