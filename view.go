package allegory

// View is an interface for a receiver of events. It can be thought
// of as a perspective on the state that can listen for events.
type View interface {
	// Initialize the state. This is useful for setting up listeners
	// on the bus.
	InitView()

	// Handle an Allegro event, such as keyboard input.
	HandleEvent(event interface{}) bool

	// Called once per frame to perform any necessary updates.
	UpdateView()

	// Clean up the state. This is useful for cleaning up listeners
	// on the bus.
	CleanupView()
}

type BaseView struct{}

func (v *BaseView) InitView()                          {}
func (v *BaseView) HandleEvent(event interface{}) bool { return false }
func (v *BaseView) UpdateView()                        {}
func (v *BaseView) CleanupView()                       {}

var _ View = (*BaseView)(nil)
