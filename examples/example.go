// +build ignore

package main

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/config"
	"github.com/dradtke/allegory/examples/hero"
	"github.com/dradtke/go-allegro/allegro"
	"os"
	"path/filepath"
	"strconv"
)

const (
	DATA_DIR    = "data"
	IMG_DIR     = DATA_DIR + "/images"
	GAME_CONFIG = DATA_DIR + "/game.cfg"
)

var (
	Hero *hero.Hero
)

/* -- Process -- */

func loadImages() {
	filepath.Walk(IMG_DIR, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		return cache.LoadImage(path, path[len(IMG_DIR)+1:])
	})
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
}

func (s *LoadingState) InitState() {
	allegory.After([]func(){loadImages, loadConfig}, func() {
		allegory.NewState(new(GameState))
	})
}

type GameState struct {
	allegory.BaseState
}

func (s *GameState) InitState() {
	var (
		frame *allegro.Bitmap
		err   error = nil
	)
	hero.ImgStanding = cache.Image("standing.png") // TODO: find a way to move this to InitActor()?
	hero.ImgWalking = make([]*allegro.Bitmap, 0)
	for i := 1; err == nil; i++ {
		frame, err = cache.FindImage("walking-" + strconv.Itoa(i) + ".png")
		if err == nil {
			hero.ImgWalking = append(hero.ImgWalking, frame)
		}
	}

	h := new(hero.Hero)
	h.State = new(hero.Standing)
	h.InitConfig(cache.Config("game"))
	allegory.AddActor(1, h)

	heroView := new(HeroView)
	heroView.Hero = h
	allegory.AddView(heroView)

	Hero = h
}

/* -- View -- */

type HeroView struct {
	allegory.BaseKeyView
	Hero *hero.Hero
}

func (v *HeroView) UpdateView() {
	switch state := v.Hero.State.(type) {
	case *hero.Standing:
		if v.IsDown[allegro.KEY_LEFT] {
			v.Hero.HandleCommand(&hero.Walk{-1})
		} else if v.IsDown[allegro.KEY_RIGHT] {
			v.Hero.HandleCommand(&hero.Walk{1})
		} else if v.IsDown[allegro.KEY_SPACE] {
			v.Hero.HandleCommand(&hero.Jump{0})
		}

	case *hero.Walking:
		if state.Dir > 0 && !v.IsDown[allegro.KEY_RIGHT] {
			v.Hero.HandleCommand(&hero.Stand{})
		} else if state.Dir < 0 && !v.IsDown[allegro.KEY_LEFT] {
			v.Hero.HandleCommand(&hero.Stand{})
		} else if v.IsDown[allegro.KEY_SPACE] {
			v.Hero.HandleCommand(&hero.Jump{state.Dir})
		}
	}
}

/* -- Main -- */

func ReadInput() {
	for {
		line := <-allegory.Stdin()
		if _, _, err := allegory.ParseAssignment(line); err == nil {
			// Parsed!
		}
	}
}

func main() {
	config.SetWindowTitle("Let's Go!")
	config.SetWindowIcons(IMG_DIR + "/standing.png")

	allegory.Run(new(LoadingState))
}
