package allegory

import (
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/allegory/console"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/dialog"
	"github.com/dradtke/go-allegro/allegro/font"
	"github.com/dradtke/go-allegro/allegro/image"
	"github.com/dradtke/go-allegro/allegro/primitives"
	"os"
	"path/filepath"
	"runtime"
)

// Init() initializes the game by creating the event queue, installing
// input systems, creating the display, and starting the FPS timer. It also
// changes the working directory to the package root relative to GOPATH,
// if one was specified.
func Init(state GameState, views ...View) {
	runtime.LockOSThread()
	var err error

	if pkg_root := config.PackageRoot(); pkg_root != "" {
		for _, dir := range filepath.SplitList(os.Getenv("GOPATH")) {
			p := filepath.Join(dir, "src", pkg_root)
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				os.Chdir(p)
				break
			}
		}
	}

	// Allegro
	if err = allegro.Install(); err != nil {
		panic(err)
	}
	_atexit = append(_atexit, allegro.Uninstall)

	// Native Dialogs Addon
	if err = dialog.Install(); err != nil {
		panic(err)
	}
	_atexit = append(_atexit, dialog.Shutdown)

	// Primitives Addon
	if err = primitives.Install(); err != nil {
		Fatal(err)
	}
	_atexit = append(_atexit, primitives.Uninstall)

	// Image addon
	if err = image.Install(); err != nil {
		Fatal(err)
	}
	_atexit = append(_atexit, image.Uninstall)

	// Font Addon
	font.Install()
	_atexit = append(_atexit, font.Uninstall)

	// Event Queue
	if _eventQueue, err = allegro.CreateEventQueue(); err != nil {
		Fatal(err)
	}

	// Keyboard Driver
	var keyboard *allegro.EventSource
	if err = allegro.InstallKeyboard(); err != nil {
		Fatal(err)
	}
	if keyboard, err = allegro.KeyboardEventSource(); err != nil {
		Fatal(err)
	} else {
		_eventQueue.RegisterEventSource(keyboard)
	}

	// Display
	allegro.SetNewDisplayFlags(config.DisplayFlags())
	w, h := config.DisplaySize()
	if _display, err = allegro.CreateDisplay(w, h); err != nil {
		Fatal(err)
	}
	_display.SetWindowTitle(config.Title())
	_eventQueue.Register(_display)
	allegro.ClearToColor(config.BlankColor())
	allegro.FlipDisplay()

	// FPS Timer
	if _fpsTimer, err = allegro.CreateTimer(1.0 / float64(config.Fps())); err != nil {
		Fatal(err)
	}
	_eventQueue.Register(_fpsTimer)
	_fpsTimer.Start()

	// Initialize subsystems.
	console.Init(_eventQueue)

	// Set the state.
	if state == nil {
		state = &BaseState{}
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

	runtime.UnlockOSThread()
}
