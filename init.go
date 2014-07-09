package gopher

import (
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/dialog"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
	"github.com/dradtke/gopher/config"
	"github.com/dradtke/gopher/console"
)

// Init() initializes the game by creating the event queue, installing
// input systems, creating the display, and starting the FPS timer.
func Init(state GameState, views ...View) {
	var err error

	// Allegro
	if err = al.Install(); err != nil {
		panic(err)
	}
	_atexit = append(_atexit, al.Uninstall)

	// Native Dialogs Addon
	if err = dialog.Install(); err != nil {
		panic(err)
	}
	_atexit = append(_atexit, dialog.Shutdown)

	// Primitives Addon
	if err = prim.Install(); err != nil {
		Fatal(err)
	}
	_atexit = append(_atexit, prim.Uninstall)

	// Image addon
	if err = image.Install(); err != nil {
		Fatal(err)
	}
	_atexit = append(_atexit, image.Uninstall)

	// Font Addon
	font.Install()
	_atexit = append(_atexit, font.Uninstall)

	// Event Queue
	if _eventQueue, err = al.CreateEventQueue(); err != nil {
		Fatal(err)
	}

	// Keyboard Driver
	var keyboard *al.EventSource
	if err = al.InstallKeyboard(); err != nil {
		Fatal(err)
	}
	if keyboard, err = al.KeyboardEventSource(); err != nil {
		Fatal(err)
	} else {
		_eventQueue.RegisterEventSource(keyboard)
	}

	// Display
	al.SetNewDisplayFlags(config.DisplayFlags())
	w, h := config.DisplaySize()
	if _display, err = al.CreateDisplay(w, h); err != nil {
		Fatal(err)
	}
	_display.SetWindowTitle(config.GameName())
	_eventQueue.Register(_display)
	al.ClearToColor(config.BlankColor())
	al.FlipDisplay()

	// FPS Timer
	if _fpsTimer, err = al.CreateTimer(1.0 / float64(config.Fps())); err != nil {
		Fatal(err)
	}
	_eventQueue.Register(_fpsTimer)
	_fpsTimer.Start()

	// Initialize subsystems.
	console.Init(_eventQueue)

	// Set the state.
	if state == nil {
		state = &BlankState{}
	}
	setState(state, views...)
}

// Cleanup() destroys some common resources and runs all necessary
// _atexit functions.
func Cleanup() {
	if _fpsTimer != nil {
		_fpsTimer.Destroy()
	}
	if _display != nil {
		_display.Destroy()
	}
	if _eventQueue != nil {
		_eventQueue.Destroy()
	}

	for _, f := range _atexit {
		f()
	}
}
