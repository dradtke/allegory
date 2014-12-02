package allegory

// Initializable is an interface for values that support initialization.
// This includes game states and actors.
type Initializable interface {
	Init()
}

// InitializableWithFailure is an interface for values that support initialization
// that can fail. This includes processes.
type InitializableWithFailure interface {
	Init() error
}

// Private variant of InitializableWithFailure for internal process definitions.
type privatelyInitializableWithFailure interface {
	init() error
}

// Updateable is an interface for values that support frame-by-frame updates.
// This includes game states and actors.
type Updateable interface {
	Update()
}

// UpdateableStatefully is an interface for values that support frame-by-frame updates
// that modify an internal state. This includes actors.
type UpdateableStatefully interface {
	Update() interface{}
}

// Renderable is an interface for values that support rendering. This includes
// game states and actors.
type Renderable interface {
	Render(delta float32)
}

// Cleanupable is an interface for values that support end-of-life cleanup. This
// includes game states and actors.
type Cleanupable interface {
	Cleanup()
}

// EventHandler is an interface for values that can receive Allegro events.
// This includes game states.
type EventHandler interface {
	HandleEvent(event interface{}) bool
}

// EventHandler is an interface for values that can receive Allegro events and
// also need to modify internal state. This includes actors.
type StatefulEventHandler interface {
	HandleEvent(event interface{}) interface{}
}

// Messagable is an interface for processes that can handle
// messages from the system. Quit and Tick signals are handled
// automatically and will never be passed to HandleMessage().
type Messagable interface {
	HandleMessage(msg interface{}) error
}

// Private variant of Messagable for internal process definitions.
type privatelyMessagable interface {
	handleMessage(msg interface{}) error
}

// Tickable is an interface for processes that need to do something
// on each frame.
type Tickable interface {
	Tick() (bool, error)
}

// Private variant of Tickable for internal process definitions.
type privatelyTickable interface {
	tick() (bool, error)
}

// Continuable is an interface for processes that need to kick off
// another one when this one finishes.
type Continuable interface {
	Next() interface{}
}
