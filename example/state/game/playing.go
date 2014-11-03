package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/example/actors/hero"
)

type PlayingState struct {
	allegory.BaseGameState

	hero     *hero.Hero
}

func (s *PlayingState) InitGameState() {
	cfg := cache.Config("game")

	s.hero = hero.NewHero().InitConfig(cfg)
	allegory.AddActor(1, s.hero)
	allegory.AddView(hero.NewHeroView(s.hero, s.PauseGame).InitConfig(cfg))
}

func (s *PlayingState) OnResume() {
	s.hero.NeedsStateValidation = true
}

func (s *PlayingState) PauseGame() {
	allegory.PushState(new(PausedState))
}
