// Package bus manages the engine's event system.
//
// While computer events, such as keyboard input, are handled by
// Allegro's built-in event system, this package provides a
// flexible API for defining custom events. Each event should be
// given a unique event id. While some events are defined by the
// Gopher engine, their id's are very high, so if you start at 1
// and work your way up, there shouldn't be any problems with overlap.
// Creating a custom event is simply a matter of settling on
// an id for it and using it to register handlers, then signal
// it whenever it happens. Handlers are run synchronously, so if
// they need to perform some long computation, then they should kick
// off their own goroutine.
package bus

import (
	"container/list"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
)

var (
	_bus     = make(map[uint]*list.List)
	_curried = make(map[*list.Element][]reflect.Value)
)

// Signal() calls all of the registered listeners for a given
// event type, as long as the parameters exactly match the ones
// that were passed into this function. For example, this works:
//
//      const MyEventId uint = 1
//
//      func onMyEventTrigger(val string) {
//          fmt.Printf("my event trigger received: %s\n", val)
//      }
//
//      func main() {
//          bus.AddListener(MyEventId, onEventTrigger)
//          bus.Signal(MyEventId, "hello signals!")
//      }
//
// ...but if onMyEventTrigger() took anything except exactly
// one string parameter, then it would not be called and a warning
// would be issued out to standard error.
//
// As long as the parameters line up, listeners can take any number
// of parameters, including 0.
//
func Signal(eventType uint, params ...interface{}) {
	listeners, ok := _bus[eventType]
	if !ok || listeners.Len() == 0 {
		return
	}
	paramValues := make([]reflect.Value, len(params))
	for i, param := range params {
		paramValues[i] = reflect.ValueOf(param)
	}
	numParams := len(paramValues)
l:
	for e := listeners.Front(); e != nil; e = e.Next() {
		curriedValues := _curried[e]
		numCurried := len(curriedValues)
		n := numParams + numCurried
		f := reflect.ValueOf(e.Value)
		t := f.Type()
		if t.NumIn() != n {
			fmt.Fprintf(os.Stderr, "invalid callback registerd for event type %d: "+
				"need %d parameters, but have %d\n", eventType, n, t.NumIn())
			continue l
		}
		allValues := make([]reflect.Value, n)
		for i := 0; i < n; i++ {
			in := t.In(i)
			failed := false
			if i < numCurried {
				if in != paramValues[i].Type() {
					failed = true
				} else {
					allValues[i] = curriedValues[i]
				}
			} else {
				if in != paramValues[i-numCurried].Type() {
					failed = true
				} else {
					allValues[i] = paramValues[i-numCurried]
				}
			}
			if failed {
				// TODO: if it's convertible to the desired type, then convert it
				fmt.Fprintf(os.Stderr, "invalid callback registered for event type %d: "+
					"need %s parameter, but have %s\n",
					eventType, paramValues[i].Type().Name(), t.In(i).Name())
				continue l
			}
		}
		f.Call(allValues)
	}
}

// AddListener() registers a handler for a given event type.
func AddListener(eventType uint, f interface{}, curry ...interface{}) error {
	if reflect.ValueOf(f).Kind() != reflect.Func {
		return errors.New("cannot register non-func callback")
	}
	eventBus, ok := _bus[eventType]
	if !ok {
		eventBus = new(list.List)
		_bus[eventType] = eventBus
	}
	e := eventBus.PushBack(f)
	curriedValues := make([]reflect.Value, len(curry))
	for i, x := range curry {
		curriedValues[i] = reflect.ValueOf(x)
	}
	_curried[e] = curriedValues
	return nil
}

// RemoveListener() unregisters a handler for a given event type.
func RemoveListener(eventType uint, f interface{}) error {
	listeners := _bus[eventType]
	for e := listeners.Front(); e != nil; e = e.Next() {
		if &e.Value == &f {
			listeners.Remove(e)
			delete(_curried, e)
			return nil
		}
	}
	return errors.New("event handler not found")
}

// Clear() unregisters all handlers on the bus for all event types,
// then immediately runs a garbage collection.
func Clear() {
	for eventType, listeners := range _bus {
		for e := listeners.Front(); e != nil; e = e.Next() {
			delete(_curried, e)
		}
		listeners.Init()
		delete(_bus, eventType)
	}
	runtime.GC()
}
