package allegory

import (
	"container/list"
	"github.com/dradtke/allegory/config"
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
func initialize(state *gameState) {
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
	_display.SetWindowTitle(config.WindowTitle())
	if icons := config.WindowIcons(); icons != nil {
		_displayIcons = make([]*allegro.Bitmap, 0)
		for _, icon := range icons {
			bmp, err := allegro.LoadBitmap(icon)
			if err != nil {
				Error(err)
				continue
			}
			_displayIcons = append(_displayIcons, bmp)
		}
		_display.SetDisplayIcons(_displayIcons)
	}
	_eventQueue.Register(_display)
	allegro.ClearToColor(config.BlankColor())
	allegro.FlipDisplay()

	// FPS Timer
	if _fpsTimer, err = allegro.CreateTimer(1.0 / float64(config.Fps())); err != nil {
		Fatal(err)
	}
	_eventQueue.Register(_fpsTimer)
	_fpsTimer.Start()

	_state = stateStack{list.New()}
	_processes = make(map[*gameState][]interface{})
	_actors = make(map[*gameState][]interface{})
	_actorLayers = make(map[*gameState]map[uint][]interface{})
	_actorStates = make(map[interface{}]interface{})
	_messengers = make(map[interface{}]chan interface{})
	_pressedKeys = make(map[allegro.KeyCode]bool)
}

// cleanup() destroys some common resources and runs all necessary
// _atexit functions.
func cleanup() {
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

// Run() initializes Allegro and Allegory and kicks off the main game loop.
// It won't return until the game ends.
func Run(initialState StateID) {
	allegro.Run(func() {
		defer cleanup()
		state, ok := _stateMap[initialState]
		if !ok {
			Errorf("allegory.Run() called with invalid state id: %s", initialState)
			return
		}
		initialize(state)
		PushState(initialState)
		loop()
	})
}
