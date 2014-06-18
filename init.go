package gopher

import (
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/go-allegro/allegro/dialog"
	"github.com/dradtke/go-allegro/allegro/font"
	prim "github.com/dradtke/go-allegro/allegro/primitives"
	"github.com/dradtke/gopher/config"
)

// Init() initializes the game by creating the event queue, installing
// input systems, creating the display, and starting the FPS timer.
func Init(s GameState) {
	var err error

    // Allegro
    if err = al.Install(); err != nil {
        panic(err)
    }
    atexit = append(atexit, al.Uninstall)

    // Native Dialogs Addon
    if err = dialog.Install(); err != nil {
        panic(err)
    }
    atexit = append(atexit, dialog.Shutdown)

    // Primitives Addon
    if err = prim.Install(); err != nil {
        Fatal(err)
    }
    atexit = append(atexit, prim.Uninstall)

    // Font Addon
    font.Install()
    atexit = append(atexit, font.Uninstall)

	// Initialize subsystems.
	// TODO: figure out why this breaks things when debugging...
	//console.Init(gopher.EventQueue())

    // Event Queue
	if eventQueue, err = al.CreateEventQueue(); err != nil {
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
		eventQueue.RegisterEventSource(keyboard)
	}

    // Display
	al.SetNewDisplayFlags(config.DisplayFlags())
    w, h := config.DisplaySize()
	if display, err = al.CreateDisplay(w, h); err != nil {
		Fatal(err)
	}
	display.SetWindowTitle(config.GameName())
	eventQueue.Register(display)
    al.ClearToColor(config.BlankColor())
    al.FlipDisplay()

    // FPS Timer
	if fpsTimer, err = al.CreateTimer(1.0 / float64(config.Fps())); err != nil {
		Fatal(err)
	}
	eventQueue.Register(fpsTimer)
	fpsTimer.Start()

    // Set the state.
    if s == nil {
        s = &BlankState{}
    }
    newState(s)
}

// Cleanup() destroys some common resources and runs all necessary
// atexit functions.
func Cleanup() {
	if fpsTimer != nil {
		fpsTimer.Destroy()
	}
	if display != nil {
		display.Destroy()
	}
	if eventQueue != nil {
		eventQueue.Destroy()
	}

    for _, f := range atexit {
        f()
    }
}
