package allegory

import (
	"errors"
	"github.com/dradtke/go-allegro/allegro"
	"math"
)

var (
	NoFrames = errors.New("no frames were provided for this animation")
)

/* -- DelayProcess -- */

// DelayProcess is a process that waits a set amount of time,
// then takes action.
type DelayProcess struct {
	BaseProcess

	timer uint

	// Delay is the number of ticks
	// before activating.
	Delay uint

	// Activate is the function called once time runs out.
	Activate func()

	// Successor is the process to kick off after Activate
	// is called.
	Successor Process
}

func (p *DelayProcess) Tick() (bool, error) {
	p.timer++
	if p.timer >= p.Delay {
		if p.Activate != nil {
			p.Activate()
		}
		return false, nil
	}
	return true, nil
}

// Next() returns a reference to the process to run once
// the delay for this one is up.
func (p *DelayProcess) Next() Process {
	return p.Successor
}

/* -- AnimationProcess -- */

type AnimationProcess struct {
	BaseProcess

	timer        uint
	currentFrame *allegro.Bitmap

	// Step is the number of ticks between frames.
	Step uint

	// Frames is a list of bitmaps, one for each frame.
	Frames []*allegro.Bitmap

	Repeat, Paused, Reversed bool
}

func (p *AnimationProcess) InitProcess() error {
	if p.Frames == nil || len(p.Frames) == 0 {
		return NoFrames
	}
	if !p.Repeat {
		p.currentFrame = p.Frames[0]
	} else {
		index := len(p.Frames) - 1
		p.timer = uint(index) * p.Step
		p.currentFrame = p.Frames[index]
	}
	return nil
}

func (p *AnimationProcess) Tick() (bool, error) {
	if p.Paused {
		return true, nil
	}
	if !p.Reversed {
		p.timer++
	} else {
		p.timer--
	}
	index := int(math.Floor(float64(p.timer) / float64(p.Step)))
	if index >= len(p.Frames) {
		if p.Repeat {
			p.timer, index = 0, 0
		} else {
			return false, nil
		}
	} else if index < 0 {
		if p.Repeat {
			index = len(p.Frames) - 1
			p.timer = uint(index) * p.Step
		} else {
			return false, nil
		}
	}
	p.currentFrame = p.Frames[index]
	return true, nil
}

func (p *AnimationProcess) HandleMessage(msg interface{}) error {
	switch msg.(type) {
	case *PauseAnimation:
		p.Paused = true
	case *ResumeAnimation:
		p.Paused = false
	case *ResetAnimation:
		p.timer = 0
	}
	return nil
}

func (p *AnimationProcess) CurrentFrame() *allegro.Bitmap {
	return p.currentFrame
}

type (
	PauseAnimation  struct{}
	ResumeAnimation struct{}
	ResetAnimation  struct{}
)
