package entity

import "github.com/relnod/evo/num"

type Eye struct {
	Dir   num.Vec2
	Range float32
	FOV   float32

	Saw float32
}

func (e *Eye) Sees(c *Creature) {
	e.Saw = c.Radius
}

func (e *Eye) Reset() {
	e.Saw = 0
}
