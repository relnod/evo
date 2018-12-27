package entity

import (
	"math"

	"github.com/relnod/evo/pkg/math64"
)

type Eye struct {
	Dir   math64.Vec2 `json:"dir"`
	Range float64     `json:"range"`
	FOV   float64     `json:"fov"`

	Count    int
	Biggest  float64
	Distance float64
}

func NewEye(eyeRange float64) *Eye {
	return &Eye{
		Dir:   math64.Vec2{},
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
