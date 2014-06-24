package gopher

import (
	"container/list"
	al "github.com/dradtke/go-allegro/allegro"
)

// Frames per second.
const FPS int = 60

var (
	_display    *al.Display
	_eventQueue *al.EventQueue
	_fpsTimer   *al.Timer

	_processes list.List
	_views     list.List
	_messengers = make(map[Process]chan interface{})

	_state GameState
	_event al.Event
	_atexit     = make([]func(), 0)
)

// Display() returns a reference to the game's display.
func Display() *al.Display {
	return _display
}

// EventQueue() returns a reference to the game's event queue.
func EventQueue() *al.EventQueue {
	return _eventQueue
}

// State() returns a reference to the game's current state.
func State() GameState {
	return _state
}
