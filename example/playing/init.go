package playing

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/example/actors"
	"github.com/dradtke/allegory/example/g"
	"github.com/dradtke/go-allegro/allegro"
)

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
