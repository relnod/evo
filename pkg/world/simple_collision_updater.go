package world

import (
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

// SimpleCollisionHandler takes a simple aproach in detecting the collision of
// creatures with creatures, by brute forcing every combination.
// Implements the evo.CollisionHandler
type SimpleCollisionHandler struct {
	width  int
	height int
}

// NewSimpleCollisionHandler returns a new simpe collisio updater.
func NewSimpleCollisionHandler(width, height int) *SimpleCollisionHandler {
	return &SimpleCollisionHandler{
		width:  width,
		height: height,
	}
}

// DetectCollisions checks the collision for all creatures.
func (s *SimpleCollisionHandler) DetectCollisions(creatures []*entity.Creature) {
	s.updateCreatureCreature(creatures)
	s.updateCreatureEyeCreature(creatures)
	s.updateCreatureEdge(creatures)
}

func (s *SimpleCollisionHandler) updateCreatureCreature(creatures []*entity.Creature) {
	for _, c1 := range creatures {
		if c1.Speed <= 0 {
			continue
		}
		for _, c2 := range creatures {
			if collision.CircleCircle(&c1.Pos, c1.Radius, &c2.Pos, c2.Radius) {
				c1.Collide(c2)
			}
		}
	}
}

func (s *SimpleCollisionHandler) updateCreatureEyeCreature(creatures []*entity.Creature) {
	for _, c1 := range creatures {
		if len(c1.Eyes) == 0 {
			continue
		}
		for _, eye := range c1.Eyes {
			for _, c2 := range creatures {
				d := math64.Vec2{X: c2.Pos.X - c1.Pos.X, Y: c2.Pos.Y - c1.Pos.Y}
				if d.Len() > eye.Range {
					continue
				}

				if math64.Angle(&d, &c1.Dir) > eye.FOV {
					continue
				}

				eye.Sees(c2)
			}
		}
	}
}

func (s *SimpleCollisionHandler) updateCreatureEdge(creatures []*entity.Creature) {
	for _, c := range creatures {
		if c.Pos.X < 0.0 {
			s.handleCreatureEdgeCollision(c, collision.LEFT)
		} else if c.Pos.X > float64(s.width) {
			s.handleCreatureEdgeCollision(c, collision.RIGHT)
		} else if c.Pos.Y < 0.0 {
			s.handleCreatureEdgeCollision(c, collision.TOP)
		} else if c.Pos.Y > float64(s.height) {
			s.handleCreatureEdgeCollision(c, collision.BOT)
		}
	}
}

func (s *SimpleCollisionHandler) handleCreatureEdgeCollision(e *entity.Creature, border int) {
	switch border {
	case collision.LEFT:
		e.Pos.X += float64(s.width)
	case collision.RIGHT:
		e.Pos.X -= float64(s.width)
	case collision.TOP:
		e.Pos.Y += float64(s.height)
	case collision.BOT:
		e.Pos.Y -= float64(s.height)
	}
}
