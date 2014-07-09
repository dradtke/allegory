package gopher

import (
	"fmt"
	"os"
)

type Process interface {
    // Do initialization on this process.
    InitProcess()

	// Handle incoming messages. If a non-nil error value is
    // returned, then the process immediately quits and no
    // successors (if any) will be run.
	HandleMessage(msg interface{}) error

	// Tell the process to step forward one frame.
	// Returning a non-nil error value will cause the process
	// to log the error and quit. Otherwise, the boolean
	// return value should be `true` to indicate that it needs
	// to continue processing, or `false` to indicate a
	// successful termination.
	// It is essentially a special case of HandleMessage().
	Tick() (bool, error)

    // Do cleanup after this process exits, but before the
    // next one (if any) is kicked off.
	CleanupProcess()
}

type ProcessContinuer interface {
	// NextProcess() is an optional method for processes that
	// need to kick off another process once they're done.
	NextProcess() Process
}

// NotifyProcess() sends an arbitrary message to a process.
func NotifyProcess(p Process, msg interface{}) {
    defer func() {
        if r := recover(); r != nil {
            // the channel was closed, but don't let it kill
            // the program
        }
    }()
	if ch, ok := _messengers[p]; ok {
		ch <- msg
	}
}

// NotifyAllProcesses() sends an arbitrary message to all running
// _processes.
func NotifyAllProcesses(msg interface{}) {
	for e := _processes.Front(); e != nil; e = e.Next() {
		NotifyProcess(e.Value.(Process), msg)
	}
}

// Close() sends a Quit message to a process.
func Close(p Process) {
	if ch, ok := _messengers[p]; ok {
		ch <- &quit{}
	}
}

// RunProcess() takes a Process and kicks it off in a new
// goroutine. That goroutine continually listens for messages
// on its internal channel and dispatches them to the defined
// handler, with two special cases:
//
//    1. Quit messages, which cause the process to quit and
//       clean up without kicking off additional _processes
//
//    2. Tick messages, which simply tell the process to
//       process one frame.
//
func RunProcess(p Process) {
	ch := make(chan interface{})
	_messengers[p] = ch
	e := _processes.PushBack(p)

	go func() {
		var (
			alive   bool  = true // is the process running?
			carryOn bool  = true // should the process kick off its successor, if any?
			err     error = nil
		)

        p.InitProcess()

		for alive {
			switch msg := <-ch; msg.(type) {
			case *quit:
				alive = false
				carryOn = false

			case *tick:
				if alive, err = p.Tick(); err != nil {
					alive = false
					carryOn = false
					fmt.Fprintf(os.Stderr, "Process exited with error message '%s'\n", err.Error())
				}

			default:
				if err := p.HandleMessage(msg); err != nil {
                    alive = false
                    carryOn = false
					fmt.Fprintf(os.Stderr, "Process handled %v with error message '%s'\n", err.Error())
				}
			}
		}

        p.CleanupProcess()

		if p, ok := p.(ProcessContinuer); carryOn && ok {
			if next := p.NextProcess(); next != nil {
				RunProcess(next)
			}
		}

		_processes.Remove(e)
		delete(_messengers, p)
		close(ch)
	}()
}

// DelayProcess is a process that waits a set amount of time,
// then takes action.
type DelayProcess struct {
	timer float32

	// Delay is the number of ticks (TODO: change to seconds)
	// before activating.
	Delay float32

	// Activate is the function called once time runs out.
	Activate func()

	// Successor is the process to kick off after Activate
	// is called.
	Successor Process
}

func (p *DelayProcess) HandleMessage(msg interface{}) error {
	return nil
}

func (p *DelayProcess) Tick() (bool, error) {
	p.timer++
	if p.timer >= p.Delay {
		if p.Activate != nil {
			p.Activate()
		}
		return false, nil
	}
	return true, nil
}

// Next() returns a reference to the process to run once
// the delay for this one is up.
func (p *DelayProcess) Next() Process {
	return p.Successor
}

type tick struct{}

type quit struct{}
