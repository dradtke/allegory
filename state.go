package allegory

import (
	"container/list"
	"runtime"
)

type StateID string

type gameState struct {
	init        func()
	update      func()
	handleEvent func(event interface{}) bool
	render      func(delta float32)
	cleanup     func()
	// TODO: add pause/resume/others
}

func DefState(id StateID) *gameState {
	s := new(gameState)
	s.init = func() {}
	s.update = func() {}
	s.handleEvent = func(_ interface{}) bool { return false }
	s.render = func(_ float32) {}
	s.cleanup = func() {}
	if _stateMap == nil {
		_stateMap = make(map[StateID]*gameState)
	}
	_stateMap[id] = s
	return s
}

func (s *gameState) Init(f func()) *gameState {
	s.init = f
	return s
}

func (s *gameState) Update(f func()) *gameState {
	s.update = f
	return s
}

func (s *gameState) HandleEvent(f func(event interface{}) bool) *gameState {
	s.handleEvent = f
	return s
}

func (s *gameState) Render(f func(delta float32)) *gameState {
	s.render = f
	return s
}

func (s *gameState) Cleanup(f func()) *gameState {
	s.cleanup = f
	return s
}

// NewState() changes the state, regardless of the status of currently
// running processes.
// TODO: rewrite this so that it's not in terms of PopState() and PushState();
// this causes too many calls to OnPause() and OnResume().
func NewState(stateId StateID) {
	PopState()
	PushState(stateId)
}

// Push a new state to the top of the stack.
func PushState(stateId StateID) {
	state, ok := _stateMap[stateId]
	if !ok {
		Errorf("tried to push invalid state '%s'!", stateId)
		return
	}
	_state.Push(state)
}

func PopState() *gameState {
	return _state.Pop()
}

// NewStateWait() waits for all processes to finish without
// blocking the current goroutine, then changes the game state.
func NewStateWait(stateId StateID) {
	go func() {
		for len(_processes) > 0 {
			runtime.Gosched()
		}
		NewState(stateId)
	}()
}

// NewStateNow() tells all processes to quit,
// waits for them to finish, then changes the game state.
func NewStateNow(stateId StateID) {
	NotifyAllProcesses(&quit{})
	for len(_processes) > 0 {
		runtime.Gosched()
	}
	NewState(stateId)
}

/* -- stateStack -- */

type stateStack struct {
	stack *list.List
}

func (s *stateStack) Empty() bool {
	return s.stack.Len() == 0
}

func (s *stateStack) Current() *gameState {
	front := s.stack.Front()
	if front == nil {
		return nil
	}
	return front.Value.(*gameState)
}

func (s *stateStack) Push(state *gameState) {
	cur := s.Current()
	if cur != nil {
		//cur.OnPause()
	}

	s.stack.PushFront(state)

	if state != nil {
		_processes[state] = make([]interface{}, 0)
		_actors[state] = make([]interface{}, 0)
		_actorLayers[state] = make(map[uint][]interface{})
		state.init()
	}
}

func (s *stateStack) Pop() *gameState {
	oldState := s.stack.Remove(s.stack.Front()).(*gameState)

	if oldState != nil {
		oldState.cleanup()

		if actors, ok := _actors[oldState]; ok {
			for _, actor := range actors {
				if actor, ok := actor.(Cleanupable); ok {
					actor.Cleanup()
				}
			}
			delete(_actors, oldState)
			delete(_actorLayers, oldState)
		}

		runtime.GC()
	}

	if cur := s.Current(); cur != nil {
		//cur.OnResume()
	}

	return oldState
}

func (s *stateStack) Update() {
	if cur := s.Current(); cur != nil {
		cur.update()
	}
}

func (s *stateStack) HandleEvent(event interface{}) bool {
	if cur := s.Current(); cur != nil {
		return cur.handleEvent(event)
	}
	return false
}

func (s *stateStack) Render(delta float32) {
	if cur := s.Current(); cur != nil {
		cur.render(delta)
	}
}

func (s *stateStack) Processes() []interface{} {
	if processes, ok := _processes[s.Current()]; ok && processes != nil {
		return processes
	}
	return make([]interface{}, 0)
}

func (s *stateStack) Actors() []interface{} {
	if actors, ok := _actors[s.Current()]; ok && actors != nil {
		return actors
	}
	return make([]interface{}, 0)
}

func (s *stateStack) ActorLayers() map[uint][]interface{} {
	if layers, ok := _actorLayers[s.Current()]; ok && layers != nil {
		return layers
	}
	return make(map[uint][]interface{})
}
