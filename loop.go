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

	for running {
		var (
			event   = _eventQueue.WaitForEvent(&_event)
			handled = false
		)

		switch e := event.(type) {
		case allegro.TimerEvent:
			if e.Source() == _fpsTimer {
				ticking, handled = true, true
			}

		case allegro.DisplayCloseEvent:
			running, handled = false, true

		case allegro.KeyDownEvent:
			_pressedKeys[e.KeyCode()] = true

		case allegro.KeyUpEvent:
			_pressedKeys[e.KeyCode()] = false
		}

		if !handled {
			_state.HandleEvent(event)
		}

		if running && ticking && _eventQueue.IsEmpty() {
			now = time.Now()
			elapsed = now.Sub(lastUpdate)
			lastUpdate = now
			lag += elapsed
			for lag >= step {
				NotifyAllProcesses(&tick{})
				for _, actor := range _state.Actors() {
					var updated bool
					if state, ok := _actorStates[actor]; ok {
						if state, ok := state.(UpdateableStatefully); ok {
							if newState := state.Update(); newState != nil {
								SetActorState(actor, newState)
							}
							updated = true
						} else if state, ok := state.(Updateable); ok {
							state.Update()
							updated = true
						}
					}
					if !updated {
						if actor, ok := actor.(Updateable); ok {
							actor.Update()
						}
					}
				}
				_state.Update()
				lag -= step
			}

			allegro.ClearToColor(config.BlankColor())

			// Render

			delta := float32(lag / step)
			_state.Render(delta)

			//allegro.HoldBitmapDrawing(true) // ???: why does this kill it?
			actorLayers := _state.ActorLayers()
			for i := uint(0); i <= _highestLayer; i++ {
				layer, ok := actorLayers[i]
				if !ok {
					continue
				}
				for _, actor := range layer {
					var rendered bool
					if state, ok := _actorStates[actor]; ok {
						if state, ok := state.(Renderable); ok {
							state.Render(delta)
							rendered = true
						}
					}
					if !rendered {
						if actor, ok := actor.(Renderable); ok {
							actor.Render(delta)
						}
					}
				}
			}
			//allegro.HoldBitmapDrawing(false)
			allegro.FlipDisplay()

			ticking = false
		}
	}

	allegro.ClearToColor(config.BlankColor())
	allegro.FlipDisplay()

	// Tell all processes to quit immediately, then wait
	// for them to finish before exiting.
	for !_state.Empty() {
		cur := _state.Current()
		if cur == nil {
			_state.Pop()
		} else {
			NotifyAllProcesses(&quit{})
			for len(_processes[cur]) > 0 {
				runtime.Gosched()
			}
			_state.Pop()
			delete(_processes, cur)
		}
	}
}
