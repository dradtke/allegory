package allegory

import (
	"github.com/dradtke/go-allegro/allegro"
	"sync"
)

var (
	_display      *allegro.Display    // the display window
	_displayIcons []*allegro.Bitmap   // icons used in the display
	_eventQueue   *allegro.EventQueue // the global event queue
	_fpsTimer     *allegro.Timer      // the FPS timer; each tick signals a new frame

	_processes    []Process                            // an internal list of running processes
	_processMutex sync.Mutex                           // a mutex used to protect _processes
	_views        []View                               // an internal list of game views
	_messengers   = make(map[Process]chan interface{}) // an internal map from process to message channel

	_state  GameState
	_event  allegro.Event
	_atexit = make([]func(), 0)

	_actors       = make([]Actor, 0)
	_actorLayers  = make(map[uint][]Actor)
	_actorsMutex  sync.Mutex
	_highestLayer uint

	_stdin = make(chan string) // channel of data read from stdin
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

func Stdin() <-chan string {
	return _stdin
}
