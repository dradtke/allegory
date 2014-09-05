package allegory

import (
	"github.com/dradtke/go-allegro/allegro"
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
	State State
}

func (v *BaseView) InitView()    {}
func (v *BaseView) UpdateView()  {}
func (v *BaseView) CleanupView() {}

var _ View = (*BaseView)(nil)

// BaseKeyView provides an IsDown field that maps from an Allegro keycode
// to whether or not that key is currently pressed.
type BaseKeyView struct {
	BaseView
	IsDown map[allegro.KeyCode]bool
}

func (v *BaseKeyView) HandleEvent(event interface{}) bool {
	switch e := event.(type) {
	case allegro.KeyDownEvent:
		v.IsDown[e.KeyCode()] = true
	case allegro.KeyUpEvent:
		v.IsDown[e.KeyCode()] = false
	}
	return false
}

var _ PlayerView = (*BaseKeyView)(nil)

// AddView() registers a new view.
func AddView(view View) {
	viewVal := reflect.ValueOf(view)
	for viewVal.Kind() == reflect.Interface || viewVal.Kind() == reflect.Ptr {
		viewVal = viewVal.Elem()
	}
	if base := viewVal.FieldByName("BaseKeyView"); base.IsValid() {
		initBaseKeyView(base)
		viewVal = base
	}
	if base := viewVal.FieldByName("BaseView"); base.IsValid() {
		initBaseView(base, reflect.ValueOf(_state))
	}
	view.InitView()
	_views = append(_views, view)
}

func RemoveView(view View) {
	for i, v := range _views {
		if v == view {
			view.CleanupView()
			_views = append(_views[:i], _views[i+1:]...)
			return
		}
	}
}

func initBaseKeyView(baseViewVal reflect.Value) {
	if isDown := baseViewVal.FieldByName("IsDown"); isDown.IsValid() {
		isDown.Set(reflect.MakeMap(isDown.Type()))
	}
}

var stateType = reflect.TypeOf((*State)(nil)).Elem()

func initBaseView(baseViewVal, stateVal reflect.Value) {
	if s := baseViewVal.FieldByName("State"); s.IsValid() && s.Type().Implements(stateType) {
		s.Set(stateVal)
	}
}
