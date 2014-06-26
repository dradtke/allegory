package gopher

// View is an interface for a receiver of events. It can be thought
// of as a perspective on the state that can listen for events.
type View interface {
    // Initialize the state. This is useful for setting up listeners
    // on the bus.
    InitView()

    // Handle an Allegro event, such as keyboard input.
    HandleEvent(msg interface{}) bool

    // Clean up the state. This is useful for cleaning up listeners
    // on the bus.
    CleanupView()
}
