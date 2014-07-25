package allegory

import (
	"container/list"
	"github.com/dradtke/go-allegro/allegro"
)

// Frames per second.
const FPS int = 60

var (
	_display    *allegro.Display
	_eventQueue *allegro.EventQueue
	_fpsTimer   *allegro.Timer

	_processes  list.List
	_views      list.List
	_messengers = make(map[Process]chan interface{})

	_state  GameState
	_event  allegro.Event
	_atexit = make([]func(), 0)

	_lastActorId  ActorId
	_actors       = make(map[ActorId]Actor)
	_actorLayers  = make(map[uint][]Actor)
	_highestLayer uint
)

// Display() returns a reference to the game's display.
func Display() *allegro.Display {
	return _display
}

// EventQueue() returns a reference to the game's event queue.
func EventQueue() *allegro.EventQueue {
	return _eventQueue
}

// State() returns a reference to the game's current state.
func State() GameState {
	return _state
}
