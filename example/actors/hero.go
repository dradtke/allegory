package actors

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/bus"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/example/signals"
	"github.com/dradtke/go-allegro/allegro"
	"strconv"
)

type Hero struct {
	allegory.Actor
	Gravity   float32
	GroundY   float32
	Jumpspeed float32
	Walkspeed uint
}

func (h *Hero) Init() {
	h.GroundY = h.Y
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
	dir  int8
}

func (h *heroStanding) Render(delta float32) {
	cache.Image("standing.png").Draw(h.hero.X, h.hero.Y, dirToFlags(h.dir))
}

func (h *heroStanding) HandleEvent(event interface{}) interface{} {
	switch ev := event.(type) {
	case allegro.KeyDownEvent:
		switch ev.KeyCode() {
		case allegro.KEY_LEFT:
			return &heroWalking{hero: h.hero, dir: -1}
		case allegro.KEY_RIGHT:
			return &heroWalking{hero: h.hero, dir: 1}
		case allegro.KEY_SPACE:
			return &heroJumping{hero: h.hero, dir: h.dir, jumpspeed: -h.hero.Jumpspeed}
		}

	case allegro.KeyUpEvent:
		switch ev.KeyCode() {
		case allegro.KEY_LEFT:
			if allegory.KeyDown(allegro.KEY_RIGHT) {
				return &heroWalking{hero: h.hero, dir: 1}
			}
		case allegro.KEY_RIGHT:
			if allegory.KeyDown(allegro.KEY_LEFT) {
				return &heroWalking{hero: h.hero, dir: -1}
			}
		}
	}
	return nil
}

/* -- Walking -- */

type heroWalking struct {
	hero      *Hero
	dir       int8
	animation *allegory.AnimationProcess
}

func (h *heroWalking) Init() {
	images := make([]*allegro.Bitmap, 0)
	var (
		err   error
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
	h.animation.CurrentFrame().Draw(x, y, dirToFlags(h.dir))
}

func (h *heroWalking) Update() interface{} {
	left, right := allegory.KeyDown(allegro.KEY_LEFT), allegory.KeyDown(allegro.KEY_RIGHT)
	if left == right {
		allegory.Close(h.animation)
		return h.hero.Standing(h.dir)
	}
	if h.dir > 0 && left {
		h.dir = -1
	} else if h.dir < 0 && right {
		h.dir = 1
	}
	h.hero.Move(float32(h.hero.Walkspeed)*float32(h.dir), 0)
	return nil
}

func (h *heroWalking) HandleEvent(event interface{}) interface{} {
	switch ev := event.(type) {
	case allegro.KeyDownEvent:
		switch ev.KeyCode() {
		case allegro.KEY_SPACE:
			return &heroJumping{hero: h.hero, dir: h.dir, jumpspeed: -h.hero.Jumpspeed, velocity: h.hero.Walkspeed}
		}
	}
	return nil
}

/* -- Jumping -- */

type heroJumping struct {
	hero      *Hero
	dir       int8
	velocity  uint
	jumpspeed float32
}

func (h *heroJumping) Update() interface{} {
	h.hero.Move(float32(h.dir)*float32(h.velocity), h.jumpspeed)
	h.jumpspeed += h.hero.Gravity
	if h.hero.Y >= h.hero.GroundY {
		h.hero.Y = h.hero.GroundY
		bus.Signal(signals.HERO_LANDED)
		left, right := allegory.KeyDown(allegro.KEY_LEFT), allegory.KeyDown(allegro.KEY_RIGHT)
		if left == right {
			return h.hero.Standing(h.dir)
		} else {
			return h.hero.Walking(h.dir)
		}
	}
	return nil
}

func (h *heroJumping) Render(delta float32) {
	cache.Image("standing.png").Draw(h.hero.X, h.hero.Y, dirToFlags(h.dir))
}
