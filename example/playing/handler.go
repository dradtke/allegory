package playing

import (
	"github.com/dradtke/allegory"
)

func HandleEvent(event interface{}) bool {
	state := allegory.ActorState(_hero)
	if state, ok := state.(allegory.StatefulEventHandler); ok {
		if newState := state.HandleEvent(event); newState != nil {
			allegory.SetActorState(_hero, newState)
		}
	}
	return false
}
