package allegory

import (
	"sync"
)

// After() takes a list of functions and kicks each one off in its own goroutine,
// then calls the callback once they've all finished. Everything is run
// in a separate goroutine, so After() returns almost immediately.
func After(routines []func(), callback func()) {
	var wg sync.WaitGroup
	wg.Add(len(routines))
	for _, routine := range routines {
		go func(f func()) {
			f()
			wg.Done()
		}(routine)
	}
	go func() {
		wg.Wait()
		callback()
	}()
}
