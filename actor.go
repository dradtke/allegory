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

	actorType = reflect.TypeOf((*Actor)(nil)).Elem()
	viewType  = reflect.TypeOf((*View)(nil)).Elem()
)

/* -- StatefulActor -- */

// A StatefulActor is an extension of BaseActor that makes it easy to update
// the state of an actor.
type StatefulActor struct {
	BaseActor
	State  ActorState
	Parent Actor // the struct embedding StatefulActor
}

func (a *StatefulActor) ChangeState(state ActorState) {
	if state == nil {
		return
	}
	a.State = state
	a.State.InitActorState(a.Parent)
}

func (a *StatefulActor) InitActor() {
	a.BaseActor.InitActor()
	if a.State != nil {
		a.State.InitActorState(a.Parent)
	}
}

func (a *StatefulActor) UpdateActor() {
	a.BaseActor.UpdateActor()
	if a.State != nil {
		newState := a.State.UpdateActorState(a.Parent)
		if newState != nil {
			a.ChangeState(newState)
		}
	}
}

func (a *StatefulActor) RenderActor(delta float32) {
	if a.State != nil {
		a.State.RenderActorState(a.Parent, delta)
	}
}

/* -- ActorState -- */

type ActorState interface {
	InitActorState(Actor)
	UpdateActorState(Actor) ActorState
	RenderActorState(Actor, float32)
	CleanupActorState(Actor)
}

type BaseActorState struct{}

func (a *BaseActorState) InitActorState(actor Actor)                  {}
func (a *BaseActorState) UpdateActorState(actor Actor) ActorState     { return nil }
func (a *BaseActorState) RenderActorState(actor Actor, delta float32) {}
func (a *BaseActorState) CleanupActorState(actor Actor)               {}

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
	// TODO: calling Elem() before this breaks initStatefulActor
	if base := actorVal.Elem().FieldByName("StatefulActor"); base.IsValid() {
		initStatefulActor(base, actorVal)
		actorVal = base
	}
	if base := actorVal.FieldByName("BaseActor"); base.IsValid() {
		initBaseActor(base, id)
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

// If the actor embeds BaseActor, set its Id field.
func initBaseActor(actorVal reflect.Value, id ActorId) {
	if idField := actorVal.FieldByName("Id"); idField.IsValid() && idField.Type() == actorIdType {
		idField.Set(reflect.ValueOf(id))
	}
}
