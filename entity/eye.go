package entity

import "github.com/relnod/evo/num"

type Eye struct {
	Dir   num.Vec2
	Range float32
	FOV   float32

	Count   int
	Biggest float32
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
