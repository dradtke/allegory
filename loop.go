package gopher

import (
	"fmt"
	"github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/gopher/config"
	"github.com/dradtke/gopher/console"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Loop() is the main game loop.
func Loop() {
	var (
		running = true
		ticking = false

		secondsPerFrame = 1 / float64(config.Fps())
		step            = time.Duration(secondsPerFrame * float64(time.Second))
		lastUpdate      = time.Now()
		lag             time.Duration
		now             time.Time
		elapsed         time.Duration
	)

	// Provide a more readable stack-trace on runtime panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%v\n", r)
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			skip := 4 // skip past where we are now + two panic.c locations + os_<something>.c
			for {
				if _, file, line, ok := runtime.Caller(skip); ok && filepath.Ext(file) == ".go" {
					if rel, err := filepath.Rel(cwd, file); err == nil {
						fmt.Fprintf(os.Stderr, "    %s:%d\n", rel, line)
					} else {
						break
					}
					skip += 1
				} else {
					break
				}
			}
		}
	}()

	for running {
		ev := _eventQueue.WaitForEvent(&_event)

		switch e := ev.(type) {
		case allegro.TimerEvent:
			if e.Source() == _fpsTimer {
				ticking = true
				goto eventHandled
			}

		case allegro.DisplayCloseEvent:
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
			now = time.Now()
			elapsed = now.Sub(lastUpdate)
			lastUpdate = now
			lag += elapsed
			for lag >= step {
				NotifyAllProcesses(&tick{})
				for _, actor := range _actors {
					actor.UpdateActor()
				}
				for e := _views.Front(); e != nil; e = e.Next() {
					e.Value.(View).UpdateView()
				}
				lag -= step
			}

			delta := float32(lag / step)
			_state.RenderState(delta) // ???: is this needed with the actors?
			//allegro.HoldBitmapDrawing(true) // ???: why does this kill it?
			for _, actor := range _actors {
				if actor, ok := actor.(RenderableActor); ok {
					actor.RenderActor(delta)
				}
			}
			//allegro.HoldBitmapDrawing(false)
			console.Render()
			allegro.FlipDisplay()
			allegro.ClearToColor(config.BlankColor())

			ticking = false
		}
	}

	allegro.ClearToColor(config.BlankColor())
	allegro.FlipDisplay()

	_display.SetWindowTitle("Shutting down...")
	//console.Save()

	// Tell all processes to quit immediately, then wait
	// for them to finish before exiting.
	NotifyAllProcesses(&quit{})
	for _processes.Len() > 0 {
		runtime.Gosched()
	}
}
