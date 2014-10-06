package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
)

var (
	_playingState *PlayingState
)

type PlayingState struct {
	allegory.BaseState

	heroView *HeroView
	hero     *Hero
}

func (s *PlayingState) InitState() {
	cfg := cache.Config("game")

	if s.hero == nil {
		s.hero = new(Hero)
		s.hero.State = NewStanding(s.hero)
		s.hero.InitConfig(cfg)
	} else {
		if _, ok := s.hero.State.(*Jumping); !ok {
			s.hero.State = &forceValidation{}
		}
	}
	allegory.AddActor(1, s.hero)

	if s.heroView == nil {
		s.heroView = new(HeroView)
		s.heroView.hero = s.hero
		s.heroView.InitConfig(cfg)
	}
	allegory.AddView(s.heroView)

	_playingState = s
}
