package actors

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/example/util"
	"github.com/dradtke/go-allegro/allegro"
	"strconv"
)

type Hero struct {
	allegory.Actor
	Walkspeed uint
}

func (h *Hero) Init() {
	allegory.Debugf("hero walkspeed: %d", h.Walkspeed)
}


func (h *Hero) Standing(dir int8) *heroStanding {
	return &heroStanding{h, dir}
}

func (h *Hero) Walking(dir int8) *heroWalking {
	return &heroWalking{h, dir, nil}
}

/* -- Standing -- */

type heroStanding struct {
	hero *Hero
	dir int8
}

func (h *heroStanding) Render(delta float32) {
	cache.Image("standing.png").Draw(h.hero.X, h.hero.Y, util.DirToFlags(h.dir))
}

func (h *heroStanding) HandleEvent(event interface{}) interface{} {
	switch ev := event.(type) {
	case allegro.KeyDownEvent:
		switch ev.KeyCode() {
		case allegro.KEY_LEFT:
			return &heroWalking{hero: h.hero, dir: -1}
		case allegro.KEY_RIGHT:
			return &heroWalking{hero: h.hero, dir: 1}
		}
	}
	return nil
}

/* -- Walking -- */

type heroWalking struct {
	hero *Hero
	dir int8
	animation *allegory.AnimationProcess
}

func (h *heroWalking) Init() {
	images := make([]*allegro.Bitmap, 0)
	var (
		err error
		frame *allegro.Bitmap
	)
	for i := 1; err == nil; i++ {
		if frame, err = cache.FindImage("walking-" + strconv.Itoa(i) + ".png"); err == nil {
			images = append(images, frame)
		}
	}
	h.animation = &allegory.AnimationProcess{Repeat: true, Step: 6, Frames: images}
	allegory.RunProcess(h.animation)
}

func (h *heroWalking) Render(delta float32) {
	x, y := h.hero.CalculatePos(delta)
	h.animation.CurrentFrame().Draw(x, y, util.DirToFlags(h.dir))
}

func (h *heroWalking) Update() interface{} {
	h.hero.Move(float32(h.hero.Walkspeed) * float32(h.dir), 0)
	if !allegory.KeyDown(allegro.KEY_LEFT) && !allegory.KeyDown(allegro.KEY_RIGHT) {
		allegory.Close(h.animation)
		return h.hero.Standing(h.dir)
	}
	return nil
}
