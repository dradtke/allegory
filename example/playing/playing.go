package playing

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/example/actors"
	"github.com/dradtke/allegory/example/g"
	"github.com/dradtke/go-allegro/allegro"
)

var (
	_hero *actors.Hero
)

func Register() {
	allegory.DefState("playing").
		Init(Init).
		Update(Update).
		HandleEvent(HandleEvent)
}

func Init() {
	var (
		err error
		cfg *allegro.Config
	)

	err = cache.LoadImages(g.IMG_DIR)
	if err != nil {
		allegory.Fatal(err)
	}

	cfg, err = allegro.LoadConfig(g.GAME_CONFIG)
	if err != nil {
		allegory.Fatal(err)
	}

	_hero = new(actors.Hero)
	_hero.X, _hero.Y = 200, 200
	allegory.ReadConfig(cfg, "Hero", _hero)
	allegory.AddActor(1, _hero, _hero.Standing(1))
}

func HandleEvent(event interface{}) bool {
	if event, ok := event.(allegro.KeyDownEvent); ok {
		if event.KeyCode() == allegro.KEY_ENTER {
			allegory.PushState("playing/paused")
			return true
		}
	}

	heroState := allegory.ActorState(_hero)
	if heroState, ok := heroState.(allegory.StatefulEventHandler); ok {
		if newState := heroState.HandleEvent(event); newState != nil {
			allegory.SetActorState(_hero, newState)
		}
	}

	return false
}

func Update() {
	// TODO: write updates here
}
