// +build ignore

package main

import (
	"github.com/dradtke/allegory"
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
	Images map[string]*allegro.Bitmap
	Config *allegro.Config
	Hero   *hero.Hero
)

/* -- Process -- */

// Init is a process struct for loading game files.
type Init struct {
	allegory.BaseProcess
	Finished func()

	done chan struct{}
}

// The loading process kicks off a goroutine that loads the necessary game
// assets.
func (p *Init) InitProcess() error {
	p.done = make(chan struct{}, 1)
	Images = make(map[string]*allegro.Bitmap)
	go p.loadAll()
	return nil
}

func (p *Init) loadImages() {
	filepath.Walk(IMG_DIR, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		if bmp, err := allegro.LoadBitmap(path); err == nil {
			Images[path[len(IMG_DIR)+1:]] = bmp
		}
		return nil
	})
}

func (p *Init) loadConfig() {
	cfg, err := allegro.LoadConfig(GAME_CONFIG)
	if err != nil {
		panic(err)
	}
	Config = cfg
}

func (p *Init) loadAll() {
	p.loadImages()
	p.loadConfig()
	p.done <- struct{}{}
}

// Since all the heavy lifting is done in the goroutine started by
// InitProcess(), Tick() just reports whether or not everything has finished
// loading by checking the contents of the `done` channel.
func (p *Init) Tick() (bool, error) {
	select {
	case <-p.done:
		return false, nil
	default:
		return true, nil
	}
}

// If a callback was supplied, call it to indicate that everything
// is loaded and ready to go.
func (p *Init) CleanupProcess() {
	if p.Finished != nil {
		p.Finished()
	}
}

/* -- State -- */

// LoadingState is a game state for loading assets.
type LoadingState struct {
	allegory.BaseState
}

// When the state is initialized, kick off an Init process. When it completes,
// we should change to the regular game state.
func (s *LoadingState) InitState() {
	init := new(Init)
	init.Finished = func() {
		allegory.NewState(new(GameState))
	}
	allegory.RunProcess(init)
}

type GameState struct {
	allegory.BaseState
}

func (s *GameState) InitState() {
	var (
		frame *allegro.Bitmap
		ok    bool = true
	)
	hero.ImgStanding = Images["standing.png"] // TODO: find a way to move this to InitActor()?
	hero.ImgWalking = make([]*allegro.Bitmap, 0)
	for i := 1; ok; i++ {
		frame, ok = Images["walking-"+strconv.Itoa(i)+".png"]
		if ok {
			hero.ImgWalking = append(hero.ImgWalking, frame)
		}
	}

	h := new(hero.Hero)
	h.State = new(hero.Standing)
	h.InitConfig(Config)
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
	allegory.Init(new(LoadingState))
	defer allegory.Cleanup()
	go ReadInput()
	allegory.Loop()
}
