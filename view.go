package allegory

import (
	"reflect"
)

// View is an interface for a receiver of events. It can be thought
// of as a perspective on the state that can listen for events.
type View interface {
	// Initialize the state. This is useful for setting up listeners
	// on the bus.
	InitView()

	// Called once per frame to perform any necessary updates.
	UpdateView()

	// Clean up the state. This is useful for cleaning up listeners
	// on the bus.
	CleanupView()
}

// PlayerView is an interface for Views that need to receive
// events from Allegro.
type PlayerView interface {
	View

	// Handle an Allegro event. Returns true or false indicating whether
	// or not the event was consumed.
	HandleEvent(event interface{}) bool
}

// BaseView provides a default View implementation with a State field
// for referencing the active state instance.
type BaseView struct {
	State GameState
}

func (v *BaseView) InitView()    {}
func (v *BaseView) UpdateView()  {}
func (v *BaseView) CleanupView() {}

var _ View = (*BaseView)(nil)

// AddView() registers a new view.
func AddView(view View) {
	viewVal := reflect.ValueOf(view)
	for viewVal.Kind() == reflect.Interface || viewVal.Kind() == reflect.Ptr {
		viewVal = viewVal.Elem()
	}
    cur := _state.Current()
	if base := viewVal.FieldByName("BaseView"); base.IsValid() {
		if s := base.FieldByName("State"); s.IsValid() && s.Type().Implements(stateType) {
			s.Set(reflect.ValueOf(cur))
		}
	}
	view.InitView()
	_views[cur] = append(_views[cur], view)
}

func RemoveView(view View) {
    cur := _state.Current()
	for i, v := range _views[cur] {
		if v == view {
			view.CleanupView()
			_views[cur] = append(_views[cur][:i], _views[cur][i+1:]...)
			return
		}
	}
}

var stateType = reflect.TypeOf((*GameState)(nil)).Elem()
