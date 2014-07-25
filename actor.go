package allegory

import (
	"reflect"
	"sync/atomic"
)

type ActorId uint32

func (id ActorId) Destroy() {
	_actors[id].CleanupActor()
	delete(_actors, id)
}

type Actor interface {
	InitActor()
	UpdateActor()
	CleanupActor()
}

type BaseActor struct {
	Id ActorId

	// X and Y are the coordinates of the actor.
	X, Y float32

	// Xspeed and Yspeed are speed values used to extrapolate the actor's position in times of lag.
	Xspeed, Yspeed float32
}

func (a *BaseActor) InitActor()    {}
func (a *BaseActor) UpdateActor()  {}
func (a *BaseActor) CleanupActor() {}

func NewBaseActor(x, y float32) *BaseActor {
	base := new(BaseActor)
	base.X = x
	base.Y = y
	return base
}

func (a *BaseActor) HandleCommand(msg interface{}) {}
func (a *BaseActor) Move(x, y float32)             { a.X += x; a.Y += y }
func (a *BaseActor) CalculatePos(delta float32) (x, y float32) {
	return a.X + (a.Xspeed * delta), a.Y + (a.Yspeed * delta)
}

// Ensure that BaseActor implements Actor.
var _ Actor = (*BaseActor)(nil)

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
	assignActorIdField(reflect.ValueOf(actor), id)
	actor.InitActor()
	return id
}

func FindActor(id ActorId) Actor {
	return _actors[id]
}

var actorIdType = reflect.TypeOf(ActorId(0))

// If the actor embeds BaseActor, set its Id field.
func assignActorIdField(actorVal reflect.Value, id ActorId) {
	for actorVal.Kind() == reflect.Interface || actorVal.Kind() == reflect.Ptr {
		actorVal = actorVal.Elem()
	}
	if _, ok := actorVal.Type().FieldByName("BaseActor"); ok {
		base := actorVal.FieldByName("BaseActor")
		if _, ok := base.Type().FieldByName("Id"); ok {
			i := base.FieldByName("Id")
			if i.Type() == actorIdType {
				i.Set(reflect.ValueOf(id))
			}
		}
	}
}

/* -- Actor Components -- */

type RenderableActor interface {
	RenderActor(delta float32)
}
