package allegory

import (
	"errors"
	"fmt"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/go-allegro/allegro"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// loopExiting() provides a more readable stack-trace and shows an error dialog on runtime panic.
func loopExiting() {
	r := recover()
	if r == nil {
		// everything is good
		return
	}

	var failure error
	switch r := r.(type) {
	case error:
		failure = r
	case string:
		failure = errors.New(r)
	default:
		failure = fmt.Errorf("%v", r)
	}

	fmt.Fprintf(os.Stderr, "%s\n", failure.Error())
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	skip := 3 // TODO: figure out something better for this value? Just include all .go file?
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

	Fatal(failure)
}

// loop() is the main game loop.
func loop() {
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

	//defer loopExiting()

	go readStdin()
	Debug("Starting Up...")

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

		// If the event wasn't handled, pass it to the views.
		for _, view := range _views {
			if view, ok := view.(PlayerView); ok {
				if handled := view.HandleEvent(ev); handled {
					break
				}
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
				for _, view := range _views {
					view.UpdateView()
				}
				_state.UpdateState()
				lag -= step
			}

			delta := float32(lag / step)
			if s, ok := _state.(RenderableState); ok {
				s.RenderState(delta)
			}
			//allegro.HoldBitmapDrawing(true) // ???: why does this kill it?
			for i := uint(0); i <= _highestLayer; i++ {
				layer, ok := _actorLayers[i]
				if !ok {
					continue
				}
				for _, actor := range layer {
					if actor, ok := actor.(RenderableActor); ok {
						actor.RenderActor(delta)
					}
				}
			}
			//allegro.HoldBitmapDrawing(false)
			allegro.FlipDisplay()
			allegro.ClearToColor(config.BlankColor())

			ticking = false
		}
	}

	Debug("...Shutting Down")

	allegro.ClearToColor(config.BlankColor())
	allegro.FlipDisplay()

	// Tell all processes to quit immediately, then wait
	// for them to finish before exiting.
	NotifyAllProcesses(&quit{})
	for len(_processes) > 0 {
		runtime.Gosched()
	}
}
