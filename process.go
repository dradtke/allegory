package allegory

import (
	"fmt"
	"os"
)

// NotifyProcess() sends an arbitrary message to a process.
func NotifyProcess(proc interface{}, msg interface{}) {
	defer func() {
		// don't let closed channels kill the program
		recover()
	}()
	if ch, ok := _messengers[proc]; ok {
		ch <- msg
	}
}

// NotifyAllProcesses() sends an arbitrary message to all running
// processes.
func NotifyAllProcesses(msg interface{}) {
	for _, process := range _processes[_state.Current()] {
		NotifyProcess(process, msg)
	}
}

// NotifyWhere() sends an arbitrary message to each running process
// that matches the filter criteria.
func NotifyWhere(msg interface{}, filter func(interface{}) bool) {
	for _, process := range _processes[_state.Current()] {
		if filter(process) {
			NotifyProcess(process, msg)
		}
	}
}

// Close() sends a Quit message to a process.
func Close(proc interface{}) {
	NotifyProcess(proc, &quit{})
}

// RunProcess() takes a Process and kicks it off in a new
// goroutine. That goroutine continually listens for messages
// on its internal channel and dispatches them to the defined
// handler, with two special cases:
//
//    1. Quit messages, which cause the process to quit and
//       clean up without kicking off additional processes.
//
//    2. Tick messages, which simply tell the process to
//       process one frame.
//
func RunProcess(proc interface{}) {
	var initFn func() error = nil

	if proc, ok := proc.(privatelyInitializableWithFailure); ok {
		initFn = proc.init
	} else if proc, ok := proc.(InitializableWithFailure); ok {
		initFn = proc.Init
	}

	if initFn != nil {
		if err := initFn(); err != nil {
			fmt.Fprintf(os.Stderr, "error during process initialization: %s\n", err.Error())
			return
		}
	}

	cur := _state.Current()
	ch := make(chan interface{})
	_messengers[proc] = ch
	_processMutex.Lock()
	_processes[cur] = append(_processes[cur], proc)
	_processMutex.Unlock()

	go func(cur *gameState) {
		defer func() {
			_processMutex.Lock()
			for i, process := range _processes[cur] {
				if process == proc {
					_processes[cur] = append(_processes[cur][:i], _processes[cur][i+1:]...)
					break
				}
			}
			_processMutex.Unlock()
			delete(_messengers, proc)
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
				var tickFn func() (bool, error) = nil

				if proc, ok := proc.(privatelyTickable); ok {
					tickFn = proc.tick
				} else if proc, ok := proc.(Tickable); ok {
					tickFn = proc.Tick
				}

				if tickFn != nil {
					if alive, err = tickFn(); err != nil {
						alive = false
						carryOn = false
						fmt.Fprintf(os.Stderr, "Process exited with error message '%s'\n", err.Error())
					}
				}

			default:
				var handleMessageFn func(msg interface{}) error = nil

				if proc, ok := proc.(privatelyMessagable); ok {
					handleMessageFn = proc.handleMessage
				} else if proc, ok := proc.(Messagable); ok {
					handleMessageFn = proc.HandleMessage
				}

				if handleMessageFn != nil {
					if err := handleMessageFn(msg); err != nil {
						alive = false
						carryOn = false
						fmt.Fprintf(os.Stderr, "Process handled %v with error message '%s'\n", err.Error())
					}
				}
			}
		}

		if proc, ok := proc.(Cleanupable); ok {
			proc.Cleanup()
		}

		if proc, ok := proc.(Continuable); carryOn && ok {
			if next := proc.Next(); next != nil {
				RunProcess(next)
			}
		}
	}(cur)
}

type tick struct{}

type quit struct{}
