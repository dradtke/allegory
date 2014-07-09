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
	Pos() (x, y float32)
	SetPos(x, y float32)
	CleanupActor()
}

func AddActor(actor Actor) ActorId {
	id := ActorId(atomic.AddUint32((*uint32)(&_lastActorId), 1))
	_actors[id] = actor
	return id
}

func FindActor(id ActorId) Actor {
	return _actors[id]
}

/* -- Actor Components -- */

type RenderableActor interface {
	RenderActor(delta float32)
}
