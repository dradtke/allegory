package allegory

import (
	"github.com/dradtke/go-allegro/allegro"
	"reflect"
	"sync/atomic"
)

type ActorId uint32

func (id ActorId) Destroy() {
	_actors[id].CleanupActor()
	delete(_actors, id)
}

/* -- Actor Interfaces -- */

type Actor interface {
	InitActor()
	UpdateActor()
	CleanupActor()
}

type RenderableActor interface {
	Actor
	RenderActor(delta float32)
}

type animation struct {
	images  []*allegro.Bitmap
	step    int // the number of frames it takes to advance one image
	counter int // incremented once each frame
}

/* -- BaseActor -- */

type BaseActor struct {
	Id ActorId

	// X and Y are the coordinates of the actor.
	X, Y float32

	// Width and Height are the... well... the width and the height.
	Width, Height int

	// Xspeed and Yspeed are speed values used to extrapolate the actor's position in times of lag.
	Xspeed, Yspeed float32
}

func (a *BaseActor) InitActor()    {}
func (a *BaseActor) CleanupActor() {}
func (a *BaseActor) UpdateActor()  {}

func (a *BaseActor) HandleCommand(cmd interface{}) {}
func (a *BaseActor) Move(x, y float32)             { a.X += x; a.Y += y }
func (a *BaseActor) CalculatePos(delta float32) (x, y float32) {
	return a.X + (a.Xspeed * delta), a.Y + (a.Yspeed * delta)
}

var (
	// Ensure that BaseActor implements Actor.
	_ Actor = (*BaseActor)(nil)

	baseActorType      = reflect.TypeOf((*BaseActor)(nil)).Elem()
	baseActorStateType = reflect.TypeOf((*BaseActorState)(nil)).Elem()
	statefulActorType  = reflect.TypeOf((*StatefulActor)(nil)).Elem()
	actorType          = reflect.TypeOf((*Actor)(nil)).Elem()
	viewType           = reflect.TypeOf((*View)(nil)).Elem()
)

/* -- StatefulActor -- */

// A StatefulActor is an extension of BaseActor that makes it easy to update
// the state of an actor.
type StatefulActor struct {
	BaseActor
	State  ActorState // read-only
	Parent Actor      // read-only; the struct embedding StatefulActor
}

func (a *StatefulActor) ChangeState(state ActorState) {
	if state == nil {
		return
	}
	a.State = state
	a.State.InitActorState()
}

func (a *StatefulActor) InitActor() {
	a.BaseActor.InitActor()
	if a.State != nil {
		a.State.InitActorState()
	}
}

func (a *StatefulActor) UpdateActor() {
	a.BaseActor.UpdateActor()
	if a.State != nil {
		newState := a.State.UpdateActorState()
		if newState != nil {
			a.ChangeState(newState)
		}
	}
}

func (a *StatefulActor) RenderActor(delta float32) {
	if a.State != nil {
		a.State.RenderActorState(delta)
	}
}

/* -- ActorState -- */

type ActorState interface {
	InitActorState()
	UpdateActorState() ActorState
	RenderActorState(delta float32)
	CleanupActorState()
}

type BaseActorState struct {
	Actor Actor
}

func (a *BaseActorState) InitActorState()                {}
func (a *BaseActorState) UpdateActorState() ActorState   { return nil }
func (a *BaseActorState) RenderActorState(delta float32) {}
func (a *BaseActorState) CleanupActorState()             {}

var _ ActorState = (*BaseActorState)(nil)

/* -- Related methods -- */

func AddActor(layer uint, actor Actor) ActorId {
	id := ActorId(atomic.AddUint32((*uint32)(&_lastActorId), 1))
	_actors[id] = actor
	if l, ok := _actorLayers[layer]; ok {
		l = append(l, actor)
	} else {
		_actorLayers[layer] = []Actor{actor}
	}
	if layer > _highestLayer {
		_highestLayer = layer
	}
	actorVal := reflect.ValueOf(actor)
	for actorVal.Kind() == reflect.Ptr || actorVal.Kind() == reflect.Interface {
		actorVal = actorVal.Elem()
	}
	actor.InitActor()
	return id
}

func FindActor(id ActorId) Actor {
	return _actors[id]
}

func initStatefulActor(stateful, parent reflect.Value) {
	if parent.Type().Implements(actorType) {
		stateful.FieldByName("Parent").Set(parent)
	}
}

var actorIdType = reflect.TypeOf(ActorId(0))
