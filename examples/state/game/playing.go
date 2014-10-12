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

	if s.hero == nil {
		s.hero = new(Hero)
		s.hero.State = NewStanding(s.hero)
		s.hero.InitConfig(cfg)
	} else {
		if _, ok := s.hero.State.(*Jumping); !ok {
			s.hero.needsStateValidation = true
		}
	}
	allegory.AddActor(1, s.hero)

	if s.heroView == nil {
		s.heroView = new(HeroView)
		s.heroView.hero = s.hero
		s.heroView.InitConfig(cfg)
	}
	allegory.AddView(s.heroView)
}
