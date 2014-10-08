package game

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/cache"
	"github.com/dradtke/allegory/examples/util"
	"github.com/dradtke/go-allegro/allegro"
	"math"
	"strconv"
)

const (
	XSTART, YSTART = 200, 200
)

/* -- View -- */

type HeroView struct {
	allegory.BaseKeyView
	hero *Hero

	Left  allegro.KeyCode
	Right allegro.KeyCode
	Jump  allegro.KeyCode
	Pause allegro.KeyCode
}

func (v *HeroView) InitConfig(cfg *allegro.Config) {
	if left, err := cfg.IntValue("Controls", "left"); err == nil {
		v.Left = allegro.KeyCode(left)
	} else {
		panic(err)
	}

	if right, err := cfg.IntValue("Controls", "right"); err == nil {
		v.Right = allegro.KeyCode(right)
	} else {
		panic(err)
	}

	if jump, err := cfg.IntValue("Controls", "jump"); err == nil {
		v.Jump = allegro.KeyCode(jump)
	} else {
		panic(err)
	}

	if pause, err := cfg.IntValue("Controls", "pause"); err == nil {
		v.Pause = allegro.KeyCode(pause)
	} else {
		panic(err)
	}
}

func (v *HeroView) UpdateView() {
	if _, ok := v.hero.State.(*transitioning); ok || v.hero.needsStateValidation {
		v.validateState()
		v.hero.needsStateValidation = false
	}
}

func (v *HeroView) HandleEvent(event interface{}) bool {
	v.BaseKeyView.HandleEvent(event)

	if e, ok := event.(allegro.KeyDownEvent); ok && e.KeyCode() == v.Pause {
		allegory.NewState(new(PausedState))
		return true
	}

	switch state := v.hero.State.(type) {
	case *Standing:
		switch e := event.(type) {
		case allegro.KeyDownEvent:
			switch e.KeyCode() {
			case v.Left:
				v.hero.HandleCommand(&Walk{-1})
				return true

			case v.Right:
				v.hero.HandleCommand(&Walk{1})
				return true

			case v.Jump:
				v.hero.HandleCommand(&Jump{0})
				return true
			}
		}

	case *Walking:
		switch e := event.(type) {
		case allegro.KeyUpEvent:
			switch e.KeyCode() {
			case v.Left:
				if state.dir < 0 {
					v.hero.HandleCommand(&Stand{})
					return true
				}

			case v.Right:
				if state.dir > 0 {
					v.hero.HandleCommand(&Stand{})
					return true
				}
			}

		case allegro.KeyDownEvent:
			switch e.KeyCode() {
			case v.Jump:
				v.hero.HandleCommand(&Jump{state.dir})
				return true
			}
		}
	}

	return false
}

func (v *HeroView) validateState() {
	if _, ok := v.hero.State.(*Jumping); ok {
		// remain jumping
		return
	}

	var (
		left  = v.IsDown[v.Left]
		right = v.IsDown[v.Right]
	)

	if left && !right {
		v.hero.HandleCommand(&Walk{-1})
	} else if right && !left {
		v.hero.HandleCommand(&Walk{1})
	} else {
		v.hero.HandleCommand(&Stand{})
	}
}

type Hero struct {
	allegory.StatefulActor
	GroundY float32
	dir     int8 // 1 for right, -1 for left

	needsStateValidation bool

	// configurable values
	Jumpspeed, Gravity, Walkspeed float32
}

func (h *Hero) InitActor() {
	if h.X == 0 && h.Y == 0 {
		h.X, h.Y = XSTART, YSTART
		h.GroundY = h.Y
	}

	if h.State == nil {
		h.stand()
	}
}

func (h *Hero) HandleCommand(cmd interface{}) {
	switch cmd := cmd.(type) {
	case *Stand:
		h.stand()
	case *Walk:
		h.walk(cmd.dir)
	case *Jump:
		h.jump(cmd.Intertia)
	}
}

func (h *Hero) InitConfig(cfg *allegro.Config) {
	if jumpspeed, err := cfg.Float32Value("Hero", "jumpspeed"); err == nil {
		h.Jumpspeed = jumpspeed
	} else {
		panic(err)
	}

	if gravity, err := cfg.Float32Value("Hero", "gravity"); err == nil {
		h.Gravity = gravity
	} else {
		panic(err)
	}

	if walkspeed, err := cfg.Float32Value("Hero", "walkspeed"); err == nil {
		h.Walkspeed = walkspeed
	} else {
		panic(err)
	}
}

/* -- Walking -- */

type Walking struct {
	allegory.BaseActorState
	hero *Hero

	images    []*allegro.Bitmap
	animation *allegory.AnimationProcess
	dir       int8
}

func NewWalking(h *Hero) *Walking {
	s := new(Walking)
	s.hero = h
	s.images = make([]*allegro.Bitmap, 0)

	var (
		frame *allegro.Bitmap
		err   error = nil
	)

	for i := 1; err == nil; i++ {
		frame, err = cache.FindImage("walking-" + strconv.Itoa(i) + ".png")
		if err == nil {
			s.images = append(s.images, frame)
		}
	}

	return s
}

func (s *Walking) InitActorState() {
	s.hero.dir = s.dir
	s.animation = &allegory.AnimationProcess{Repeat: true, Step: 6, Frames: s.images}
	allegory.RunProcess(s.animation)
}

func (s *Walking) UpdateActorState() allegory.ActorState {
	s.hero.Move(s.hero.Walkspeed*float32(s.hero.dir), 0)
	return nil
}

func (s *Walking) CleanupActorState() {
	allegory.Close(s.animation)
}

func (s *Walking) RenderActorState(delta float32) {
	x, y := s.hero.CalculatePos(delta)
	s.animation.CurrentFrame().Draw(x, y, util.DirToFlags(s.hero.dir))
}

func (h *Hero) walk(dir int8) {
	if _, ok := h.State.(*Walking); ok && h.dir == dir {
		return
	}
	walking := NewWalking(h)
	walking.dir = dir
	h.ChangeState(walking)
}

/* -- Standing -- */

type Standing struct {
	allegory.BaseActorState
	hero *Hero

	image *allegro.Bitmap
}

func NewStanding(h *Hero) *Standing {
	if h == nil {
		panic("NewStanding() called with nil Hero!")
	}
	s := new(Standing)
	s.hero = h
	s.image = cache.Image("standing.png")
	return s
}

func (s *Standing) RenderActorState(delta float32) {
	if s.hero == nil {
		panic("RenderActorState() called with nil s.hero!")
	}
	x, y := s.hero.CalculatePos(delta)
	s.image.Draw(x, y, util.DirToFlags(s.hero.dir))
}

func (h *Hero) stand() {
	if _, ok := h.State.(*Standing); !ok {
		h.ChangeState(NewStanding(h))
	}
}

/* -- Jumping -- */

type Jumping struct {
	allegory.BaseActorState
	hero *Hero

	image     *allegro.Bitmap
	inertia   int8
	jumpspeed float32
}

func NewJumping(h *Hero) *Jumping {
	s := new(Jumping)
	s.hero = h
	s.image = cache.Image("standing.png")
	return s
}

func (s *Jumping) InitActorState() {
	s.jumpspeed = -s.hero.Jumpspeed
}

func (s *Jumping) UpdateActorState() allegory.ActorState {
	s.hero.Move(s.hero.Walkspeed*float32(s.inertia), s.jumpspeed)
	if s.hero.Y >= s.hero.GroundY {
		s.hero.Y = s.hero.GroundY
		return new(transitioning)
	}
	s.jumpspeed = float32(math.Min(float64(s.hero.Jumpspeed), float64(s.jumpspeed+s.hero.Gravity)))
	return nil
}

func (s *Jumping) RenderActorState(delta float32) {
	x, y := s.hero.CalculatePos(delta)
	s.image.Draw(x, y, util.DirToFlags(s.hero.dir))
}

func (a *Hero) jump(inertia int8) {
	if _, ok := a.State.(*Jumping); ok {
		return
	}
	jumping := NewJumping(a)
	jumping.inertia = inertia
	a.ChangeState(jumping)
}

type transitioning struct {
	allegory.BaseActorState
}

/* -- Commands -- */

type Stand struct{}

type Walk struct {
	dir int8
}

type Jump struct {
	Intertia int8 // same as dir, but 0 if hero was standing
}

var _ allegory.Actor = (*Hero)(nil)
