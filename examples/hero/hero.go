package hero

import (
	"github.com/dradtke/allegory"
	"github.com/dradtke/allegory/examples/util"
	"github.com/dradtke/go-allegro/allegro"
	"math"
	"strconv"
)

const (
	XSTART, YSTART = 200, 200
)

var (
	ImgStanding *allegro.Bitmap
	ImgWalking  []*allegro.Bitmap
)

/* -- Standing -- */

type Standing struct {
	allegory.BaseActorState
}

func (s *Standing) RenderActorState(actor allegory.Actor, delta float32) {
	h := actor.(*Hero)
	x, y := h.CalculatePos(delta)
	ImgStanding.Draw(x, y, util.DirToFlags(h.Dir))
}

/* -- Walking -- */

type Walking struct {
	Animation *allegory.AnimationProcess
	Dir       int8
}

func (s *Walking) InitActorState(actor allegory.Actor) {
	h := actor.(*Hero)
	h.Dir = s.Dir
	s.Animation = &allegory.AnimationProcess{Repeat: true, Step: 6, Frames: ImgWalking}
	allegory.RunProcess(s.Animation)
}

func (s *Walking) UpdateActorState(actor allegory.Actor) allegory.ActorState {
	h := actor.(*Hero)
	h.Move(h.Walkspeed*float32(h.Dir), 0)
	return nil
}

func (s *Walking) CleanupActorState(actor allegory.Actor) {
	allegory.Close(s.Animation)
}

func (s *Walking) RenderActorState(actor allegory.Actor, delta float32) {
	h := actor.(*Hero)
	x, y := h.CalculatePos(delta)
	s.Animation.CurrentFrame().Draw(x, y, util.DirToFlags(h.Dir))
}

/* -- Jumping -- */

type Jumping struct {
	allegory.BaseActorState
	Inertia int8

	jumpspeed float32
}

func (s *Jumping) InitActorState(actor allegory.Actor) {
	s.jumpspeed = -actor.(*Hero).Jumpspeed
}

func (s *Jumping) UpdateActorState(actor allegory.Actor) allegory.ActorState {
	h := actor.(*Hero)
	h.Move(h.Walkspeed*float32(s.Inertia), s.jumpspeed)
	if h.Y >= h.GroundY {
		h.Y = h.GroundY
		return new(Standing)
	}
	s.jumpspeed = float32(math.Min(float64(h.Jumpspeed), float64(s.jumpspeed+h.Gravity)))
	return nil
}

func (s *Jumping) RenderActorState(actor allegory.Actor, delta float32) {
	h := actor.(*Hero)
	x, y := h.CalculatePos(delta)
	ImgStanding.Draw(x, y, util.DirToFlags(h.Dir))
}

/* -- Hero -- */

type Hero struct {
	allegory.StatefulActor
	GroundY float32
	Dir     int8

	// configurable values
	Jumpspeed, Gravity, Walkspeed float32
}

func (h *Hero) HandleCommand(cmd interface{}) {
	switch cmd := cmd.(type) {
	case *Stand:
		h.stand()
	case *Walk:
		h.walk(cmd.Dir)
	case *Jump:
		h.jump(cmd.Intertia)
	}
}

func (h *Hero) stand() {
	h.ChangeState(new(Standing))
}

func (h *Hero) walk(dir int8) {
	if _, ok := h.State.(*Walking); ok && h.Dir == dir {
		return
	}
	walking := new(Walking)
	walking.Dir = dir
	h.ChangeState(walking)
}

func (a *Hero) jump(inertia int8) {
	if _, ok := a.State.(*Jumping); ok {
		return
	}
	jumping := new(Jumping)
	jumping.Inertia = inertia
	a.ChangeState(jumping)
}

/* -- Commands -- */

type Stand struct{}

type Walk struct {
	Dir int8
}

type Jump struct {
	Intertia int8 // same as Dir, but 0 if hero was standing
}

var _ allegory.Actor = (*Hero)(nil)

func (h *Hero) InitConfig(cfg *allegro.Config) {
	parseFloat32ConfigValue(cfg, "Hero", "jumpspeed", &h.Jumpspeed)
	parseFloat32ConfigValue(cfg, "Hero", "gravity", &h.Gravity)
	parseFloat32ConfigValue(cfg, "Hero", "walkspeed", &h.Walkspeed)
}

func parseFloat32ConfigValue(cfg *allegro.Config, section, key string, value *float32) error {
	str, err := cfg.Value(section, key)
	if err != nil {
		return err
	}
	val, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return err
	}
	*value = float32(val)
	return nil
}

func (h *Hero) InitActor() {
	h.X, h.Y = XSTART, YSTART
	h.GroundY = h.Y
}
