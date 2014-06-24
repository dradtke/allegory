package gopher

import (
	"fmt"
	"os"
)

type Process interface {
	// HandleMessage() handles incoming messages.
	HandleMessage(msg interface{}) error

	// Tick() tells the process to step forward one frame.
	// Returning a non-nil error value will cause the process
	// to log the error and quit. Otherwise, the boolean
	// return value should be `true` to indicate that it needs
	// to continue processing, or `false` to indicate a
	// successful termination.
    // It is essentially a special case of HandleMessage().
	Tick(delta float32) (bool, error)
}

type ProcessCloser interface {
	// Cleanup() is an optional method for _processes that
	// need to do some cleanup once they're done.
	Cleanup()
}

type ProcessParent interface {
	// Next() is an optional method for _processes that
	// need to kick off another process once they're done.
	Next() Process
}

// Notify() sends an arbitrary message to a process.
func Notify(p Process, msg interface{}) {
	if ch, ok := _messengers[p]; ok {
		ch <- msg
	}
}

// NotifyAll() sends an arbitrary message to all running
// _processes.
func NotifyAll(msg interface{}) {
	for e := _processes.Front(); e != nil; e = e.Next() {
		Notify(e.Value.(Process), msg)
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
			alive   bool = true
			carryOn bool = true
			err     error
		)

		for alive {
			msg := <-ch
			switch m := msg.(type) {
			case *quit:
				alive = false
				carryOn = false

			case *tick:
				alive, err = p.Tick(m.delta)
				if err != nil {
					alive = false
					carryOn = false
					fmt.Fprintf(os.Stderr, "Process exited with error message '%s'\n", err.Error())
				}

			default:
                if err := p.HandleMessage(msg); err != nil {
                    fmt.Fprintf(os.Stderr, "Process handled %v with error message '%s'\n", err.Error())
                }
			}
		}

		if c, ok := p.(ProcessCloser); ok {
			c.Cleanup()
		}

		if n, ok := p.(ProcessParent); carryOn && ok {
			next := n.Next()
			if next != nil {
				RunProcess(next)
			}
		}

		_processes.Remove(e)
		close(ch)
		delete(_messengers, p)
	}()
}

// tick is a Tick message.
type tick struct {
	delta float32
}

// quit is a Quit message.
type quit struct{}

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

func (p *DelayProcess) Tick(delta float32) (bool, error) {
	p.timer++
	if p.timer >= p.Delay {
		if p.Activate != nil {
			p.Activate()
		}
		return false, nil
	}
	return true, nil
}

func (p *DelayProcess) Next() Process {
	return p.Successor
}
