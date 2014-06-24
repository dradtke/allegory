package gopher

import (
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/gopher/config"
	"github.com/dradtke/gopher/console"
	"runtime"
)

// Loop() is the main game loop.
func Loop() {
	var (
		running    = true
		ticking    = false
	)

	for running {
		ev := _eventQueue.WaitForEvent(&_event)

		switch e := ev.(type) {
		case al.TimerEvent:
			if e.Source() == _fpsTimer {
				ticking = true
				goto eventHandled
			}

		case al.DisplayCloseEvent:
			running = false
			goto eventHandled
        }

        // Check subsystems
        if console.HandleEvent(ev) {
            goto eventHandled
        }

        // Finally, pass it to the views
        for e := _views.Front(); e != nil; e = e.Next() {
            if handled := e.Value.(View).HandleEvent(ev); handled {
                break
            }
        }

	eventHandled:
		if running && ticking && _eventQueue.IsEmpty() {
			// TODO: add a delta value
			NotifyAll(&tick{0})

			_state.Render()
            console.Render()
			al.FlipDisplay()
			al.ClearToColor(config.BlankColor())

			ticking = false
		}
	}

	al.ClearToColor(config.BlankColor())
	al.FlipDisplay()

	_display.SetWindowTitle("Shutting down...")
	//console.Save()

	// Tell all processes to quit immediately, then wait
	// for them to finish before exiting.
	NotifyAll(&quit{})
	for _processes.Len() > 0 {
        runtime.Gosched()
	}
}
