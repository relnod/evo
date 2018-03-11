package system

import (
	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
	"github.com/relnod/evo/world"
)

type Collision struct {
	world *world.World

	cbCreatureBorder func(*entity.Creature, int)
}

func NewCollision(s *world.World) *Collision {
	return &Collision{
		world: s,
	}
}

func (s *Collision) Update() {
	s.CreatureCreature()
	s.CreatureEyeCreature()
	s.CreatureEdge()
}

func (s *Collision) CreatureCreature() {
	for _, cell := range s.world.Cells {
		for _, creature := range cell.Dynamic {
			for _, dynamic := range cell.Dynamic {
				if collision.CircleCircle(&creature.Pos, creature.Radius, &dynamic.Pos, dynamic.Radius) {
					creature.Collide(dynamic)
				}
			}
			for _, static := range cell.Static {
				if collision.CircleCircle(&creature.Pos, creature.Radius, &static.Pos, static.Radius) {
					creature.Collide(static)
				}
			}
		}
	}
}

func (s *Collision) CreatureEyeCreature() {
	for _, cell := range s.world.Cells {
		for _, creature := range cell.Dynamic {
			if creature.Eye == nil {
				continue
			}

			for _, dynamic := range cell.Dynamic {
				d := num.Vec2{X: dynamic.Pos.X - creature.Pos.X, Y: dynamic.Pos.Y - creature.Pos.Y}
				if d.Len() > 50 {
					continue
				}

				if num.Angle(&d, &creature.Dir) > creature.Eye.FOV {
					continue
				}

				creature.Eye.Sees(dynamic)
			}

			for _, static := range cell.Static {
				d := num.Vec2{X: static.Pos.X - creature.Pos.X, Y: static.Pos.Y - creature.Pos.Y}
				if d.Len() > 50 {
					continue
				}

				if num.Angle(&d, &creature.Dir) > creature.Eye.FOV {
					continue
				}

				creature.Eye.Sees(static)
			}
		}
	}
}

func (s *Collision) CreatureEdge() {
	for _, cell := range s.world.Cells {
		for _, c := range cell.Dynamic {
			if c.Pos.X < 0.0 {
				s.handleCreatureEdgeCollision(c, collision.LEFT)
			} else if c.Pos.X > s.world.Width {
				s.handleCreatureEdgeCollision(c, collision.RIGHT)
			} else if c.Pos.Y < 0.0 {
				s.handleCreatureEdgeCollision(c, collision.TOP)
			} else if c.Pos.Y > s.world.Height {
				s.handleCreatureEdgeCollision(c, collision.BOT)
			}
		}
	}
}

func (s *Collision) handleCreatureEdgeCollision(e *entity.Creature, border int) {
	if s.world.EdgeMode == world.EdgeModeLoop {
		switch border {
		case collision.LEFT:
			e.Pos.X += s.world.Width
		case collision.RIGHT:
			e.Pos.X -= s.world.Width
		case collision.TOP:
			e.Pos.Y += s.world.Height
		case collision.BOT:
			e.Pos.Y -= s.world.Height
		}
	}
}

func (s *Collision) FindCreature(pos *num.Vec2) *entity.Creature {
	cell := s.world.FindCell(pos)

	if cell == nil {
		return nil
	}

	for _, c := range cell.Dynamic {
		if collision.CirclePoint(&c.Pos, c.Radius+5, pos) {
			return c
		}
	}

	for _, c := range cell.Static {
		if collision.CirclePoint(&c.Pos, c.Radius+5, pos) {
			return c
		}
	}

	return nil
}
