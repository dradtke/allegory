package allegory

import (
	"fmt"
	"os"
)

type Process interface {
	// Do initialization on this process. This is done synchronously
	// before anything else, and if an error is returned, then
	// the process isn't kicked off. It should never block.
	InitProcess() error

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

type BaseProcess struct{}

func (p *BaseProcess) InitProcess() error                  { return nil }
func (p *BaseProcess) HandleMessage(msg interface{}) error { return nil }
func (p *BaseProcess) Tick() (bool, error)                 { return false, nil }
func (p *BaseProcess) CleanupProcess()                     {}

var _ Process = (*BaseProcess)(nil)

type ProcessContinuer interface {
	// NextProcess() is an optional method for processes that
	// need to kick off another process once they're done.
	NextProcess() Process
}

// NotifyProcess() sends an arbitrary message to a process.
func NotifyProcess(p Process, msg interface{}) {
	defer func() {
		// don't let closed channels kill the program
		recover()
	}()
	if ch, ok := _messengers[p]; ok {
		ch <- msg
	}
}

// NotifyAllProcesses() sends an arbitrary message to all running
// _processes.
func NotifyAllProcesses(msg interface{}) {
	for _, process := range _processes {
		NotifyProcess(process, msg)
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
	if err := p.InitProcess(); err != nil {
		fmt.Fprintf(os.Stderr, "error during process initialization: %s\n", err.Error())
		return
	}

	ch := make(chan interface{})
	_messengers[p] = ch
	_processMutex.Lock()
	_processes = append(_processes, p)
	_processMutex.Unlock()

	go func() {
		defer func() {
			_processMutex.Lock()
			for i, process := range _processes {
				if process == p {
					_processes = append(_processes[:i], _processes[i+1:]...)
					break
				}
			}
			_processMutex.Unlock()
			delete(_messengers, p)
			close(ch)
		}()

		var (
			alive   bool  = true // is the process running?
			carryOn bool  = true // should the process kick off its successor, if any?
			err     error = nil
		)

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
	}()
}

type tick struct{}

type quit struct{}
