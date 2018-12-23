package entity

import (
	"math"

	"github.com/relnod/evo/pkg/math32"
)

type Eye struct {
	Dir   math32.Vec2 `json:"dir"`
	Range float32     `json:"range"`
	FOV   float32     `json:"fov"`

	Count    int
	Biggest  float32
	Distance float32
}

func NewEye(eyeRange float32) *Eye {
	return &Eye{
		Dir:   math32.Vec2{},
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
