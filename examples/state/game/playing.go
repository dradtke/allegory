package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
)

type PlayingState struct {
	allegory.BaseGameState

	heroView *HeroView
	hero     *Hero
}

func (s *PlayingState) InitGameState() {
	cfg := cache.Config("game")

	s.hero = new(Hero)
	s.hero.State = NewStanding(s.hero)
	s.hero.InitConfig(cfg)
	allegory.AddActor(1, s.hero)

	s.heroView = new(HeroView)
	s.heroView.hero = s.hero
	s.heroView.InitConfig(cfg)
	allegory.AddView(s.heroView)
}

func (s *PlayingState) OnResume() {
	s.hero.needsStateValidation = true
}
