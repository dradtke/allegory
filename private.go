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
	_state        stateStack
	_stateMap     map[StateID]*gameState

	_processes   map[*gameState][]interface{} // an internal list of running processes
	_actors      map[*gameState][]interface{}
	_actorLayers map[*gameState]map[uint][]interface{}
	_actorStates map[interface{}]interface{}

	_messengers map[interface{}]chan interface{} // an internal map from process to message channel
	_atexit     []func()

	_actorsMutex  sync.Mutex
	_processMutex sync.Mutex // a mutex used to protect _processes

	_event        allegro.Event
	_pressedKeys  map[allegro.KeyCode]bool
	_highestLayer uint
	_stdin        = make(chan string) // channel of data read from stdin
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
func State() *gameState {
	return _state.Current()
}

func Stdin() <-chan string {
	return _stdin
}
