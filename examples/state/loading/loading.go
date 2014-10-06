package loading

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
)

const (
	DATA_DIR    = "data"
	IMG_DIR     = DATA_DIR + "/images"
	GAME_CONFIG = DATA_DIR + "/game.cfg"
)

/* -- Process -- */

func loadImages() {
	err := cache.LoadImages(IMG_DIR)
	if err != nil {
		panic(err)
	}
}

func loadConfig() {
	err := cache.LoadConfig(GAME_CONFIG, "game")
	if err != nil {
		panic(err)
	}
}

/* -- State -- */

// LoadingState is a game state for loading assets.
type LoadingState struct {
	allegory.BaseState
	OnLoad func()
}

func (s *LoadingState) InitState() {
	allegory.After([]func(){loadImages, loadConfig}, s.OnLoad)
}
