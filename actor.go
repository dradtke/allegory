package gopher

import (
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
	// X and Y are the coordinates of the actor.
	X, Y float32

	// Xspeed and Yspeed are speed values used to extrapolate the actor's position in times of lag.
	Xspeed, Yspeed float32
}

func (a *BaseActor) InitActor()    {}
func (a *BaseActor) UpdateActor()  {}
func (a *BaseActor) CleanupActor() {}

func NewBaseActor(x, y float32) BaseActor {
	return BaseActor{x, y, 0, 0}
}

func (a *BaseActor) HandleCommand(msg interface{}) {}
func (a *BaseActor) Move(x, y float32)             { a.X += x; a.Y += y }
func (a *BaseActor) CalculatePos(delta float32) (x, y float32) {
	return a.X + (a.Xspeed * delta), a.Y + (a.Yspeed * delta)
}

// Ensure that BaseActor implements Actor.
var _ Actor = (*BaseActor)(nil)

func AddActor(actor Actor) ActorId {
	id := ActorId(atomic.AddUint32((*uint32)(&_lastActorId), 1))
	_actors[id] = actor
	actor.InitActor()
	return id
}

func FindActor(id ActorId) Actor {
	return _actors[id]
}

/* -- Actor Components -- */

type RenderableActor interface {
	RenderActor(delta float32)
}
