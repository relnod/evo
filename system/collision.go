package system

import (
	"math"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

type Collision struct {
	system *System

	cbCreatureBorder func(*entity.Creature, int)

	cells []*cell
}

type cell struct {
	topLeft  num.Vec2
	botRight num.Vec2

	center num.Vec2
	radius float32

	static []*entity.Creature
	moving []*entity.Creature
}

func NewCollision(s *System, numCells int) *Collision {
	cellsPerRow := int(math.Sqrt(float64(numCells)))

	cellWidth := s.Width / float32(numCells/cellsPerRow)
	cellHeight := s.Height / float32(numCells/cellsPerRow)

	radius := (&num.Vec2{X: cellWidth, Y: cellHeight}).Len()

	cells := make([]*cell, numCells, numCells)
	for row := 0; row < cellsPerRow; row++ {
		for col := 0; col < cellsPerRow; col++ {
			cells[row*cellsPerRow+col] = &cell{
				topLeft:  num.Vec2{X: cellWidth * float32(row), Y: cellHeight * float32(col)},
				botRight: num.Vec2{X: cellWidth * float32(row+1), Y: cellHeight * float32(col+1)},
			}

			cells[row*cellsPerRow+col].center = num.Vec2{X: cells[row*cellsPerRow+col].topLeft.X + cellWidth, Y: cells[row*cellsPerRow+col].topLeft.Y + cellHeight}
			cells[row*cellsPerRow+col].radius = radius
		}
	}

	return &Collision{
		system: s,

		cells: cells,
	}
}

func (s *Collision) Update() {
	s.ResetCells()

	s.CreatureCreature()
	s.CreatureEyeCreature()
	s.CreatureBorder()

}

func (s *Collision) ResetCells() {
	for _, cell := range s.cells {
		cell.static = cell.static[:0]
		cell.moving = cell.moving[:0]
	}

	for _, c := range s.cells {
		for _, creature := range s.system.creatures {
			if collision.CircleCircle(&c.center, c.radius, &creature.Pos, creature.Radius) {
				if creature.Speed == 0 && creature.State != entity.StateChild {
					c.static = append(c.static, creature)
				} else {
					c.moving = append(c.moving, creature)
				}
			}
		}
	}
}

func (s *Collision) CreatureCreature() {
	for _, c := range s.cells {
		for _, creature := range c.moving {
			for _, moving := range c.moving {
				if collision.CircleCircle(&creature.Pos, creature.Radius, &moving.Pos, moving.Radius) {
					creature.Collide(moving)
				}
			}
			for _, static := range c.static {
				if collision.CircleCircle(&creature.Pos, creature.Radius, &static.Pos, static.Radius) {
					creature.Collide(static)
				}
			}
		}
	}
}

func (s *Collision) CreatureEyeCreature() {
	for _, c := range s.cells {
		for _, creature := range c.moving {
			if creature.Eye == nil {
				continue
			}

			for _, moving := range c.moving {
				d := num.Vec2{X: moving.Pos.X - creature.Pos.X, Y: moving.Pos.Y - creature.Pos.Y}
				if d.Len() > 50 {
					continue
				}

				if num.Angle(&d, &creature.Dir) > creature.Eye.FOV {
					continue
				}

				creature.Eye.Sees(moving)
			}

			for _, static := range c.static {
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

func (s *Collision) CreatureBorder() {
	for _, cell := range s.cells {
		for _, c := range cell.moving {
			if c.Pos.X < 0.0 {
				s.cbCreatureBorder(c, collision.LEFT)
			} else if c.Pos.X > s.system.Width {
				s.cbCreatureBorder(c, collision.RIGHT)
			} else if c.Pos.Y < 0.0 {
				s.cbCreatureBorder(c, collision.TOP)
			} else if c.Pos.Y > s.system.Height {
				s.cbCreatureBorder(c, collision.BOT)
			}
		}
	}
}

func (s *Collision) SetCreatureBorderCB(cb func(*entity.Creature, int)) {
	s.cbCreatureBorder = cb
}
