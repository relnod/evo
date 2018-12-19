package entity

import (
	"math"

	"github.com/relnod/evo/pkg/num"
)

type Eye struct {
	Dir   num.Vec2
	Range float32
	FOV   float32

	Count   int
	Biggest float32
}

func NewEye(eyeRange float32) *Eye {
	return &Eye{
		Dir:   num.Vec2{},
		Range: eyeRange,
		FOV:   (80 / eyeRange * 40) * math.Pi / 180.0,
	}
}

func (e *Eye) Sees(c *Creature) {
	e.Count++
	if c.Radius > e.Biggest {
		e.Biggest = c.Radius
	}
}

func (e *Eye) Reset() {
	e.Count = 0
	e.Biggest = 0
}
