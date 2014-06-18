package gopher

import (
	al "github.com/dradtke/go-allegro/allegro"
	"github.com/dradtke/gopher/config"
)

// Init() initializes the game by creating the event queue, installing
// input systems, creating the display, and starting the FPS timer.
func Init() {
	var err error

	// Create an event queue. All of the event sources we care about should
	// register themselves to this queue.
	if eventQueue, err = al.CreateEventQueue(); err != nil {
		Fatal(err)
	}

	// Install the Keyboard driver.
	if err = al.InstallKeyboard(); err != nil {
		Fatal(err)
	}
	if keyboard, err := al.KeyboardEventSource(); err != nil {
		Fatal(err)
	} else {
		eventQueue.RegisterEventSource(keyboard)
	}

	// Create a 640x480 window and give it a title.
	al.SetNewDisplayFlags(al.WINDOWED | al.DIRECT3D)
	if display, err = al.CreateDisplay(config.DisplayWidth(), config.DisplayHeight()); err != nil {
		Fatal(err)
	}
	display.SetWindowTitle(config.GameName())
	eventQueue.Register(display)

	// Create the FPS timer.
	if fpsTimer, err = al.CreateTimer(1.0 / float64(FPS)); err != nil {
		Fatal(err)
	}
	eventQueue.Register(fpsTimer)
	fpsTimer.Start()
}

// Cleanup() destroys some common resources.
func Cleanup() {
	if fpsTimer != nil {
		fpsTimer.Destroy()
	}

    /*
    ch := make(chan bool)
    go func() {
        if display != nil {
            display.Destroy()
        }
        ch <- true
    }()

    waiting := true
    for waiting {
        select {
        case <-ch:
            waiting = false
            fmt.Println("got it!")

        case <-time.After(100 * time.Millisecond):
            fmt.Println("waiting...")
        }
    }
    */
    if display != nil {
        display.Destroy()
    }

	if eventQueue != nil {
		eventQueue.Destroy()
	}
}
