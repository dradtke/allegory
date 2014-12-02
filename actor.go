package allegory

type Actor struct {
	X, Y, xspeed, yspeed float32
	Width, Height int
}

func (a *Actor) Move(x, y float32) { a.X += x; a.Y += y; a.xspeed, a.yspeed = x, y }
func (a *Actor) CalculatePos(delta float32) (x, y float32) {
	return a.X + (a.xspeed * delta), a.Y + (a.yspeed * delta)
}

/* -- Related methods -- */

func AddActor(layer uint, actor, state interface{}) {
	cur := _state.Current()
	if cur == nil {
		return
	}
	_actors[cur] = append(_actors[cur], actor)
	if l, ok := _actorLayers[cur][layer]; ok {
		l = append(l, actor)
	} else {
		_actorLayers[cur][layer] = []interface{}{actor}
	}
	if layer > _highestLayer {
		_highestLayer = layer
	}
	if state != nil {
		_actorStates[actor] = state
		if state, ok := state.(Initializable); ok {
			state.Init()
		}
	}
	if actor, ok := actor.(Initializable); ok {
		actor.Init()
	}
}

func SetActorState(actor, state interface{}) {
	if state == nil {
		delete(_actorStates, actor)
	} else {
		if oldState, ok := _actorStates[actor]; ok && oldState != nil {
			if oldState, ok := oldState.(Cleanupable); ok {
				oldState.Cleanup()
			}
		}
		_actorStates[actor] = state
		if state, ok := state.(Initializable); ok {
			state.Init()
		}
	}
}

func ActorState(actor interface{}) interface{} {
	return _actorStates[actor]
}

func DestroyActor(actor interface{}) {
	cur := _state.Current()
	if cur == nil {
		return
	}
	actors := _state.Actors()
	for i, a := range actors {
		if a == actor {
			actors = append(actors[:i], actors[i+1:]...)
			break
		}
	}
	for i := uint(0); i < _highestLayer; i++ {
		layer, ok := _actorLayers[cur][i]
		if !ok {
			continue
		}
		for j, a := range layer {
			if a == actor {
				layer = append(layer[:j], layer[j+1:]...)
			}
		}
		_actorLayers[cur][i] = layer
	}
	if actor, ok := actor.(Cleanupable); ok {
		actor.Cleanup()
	}
}
