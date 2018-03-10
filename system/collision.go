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

	creatures []*entity.Creature
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
		cell.creatures = cell.creatures[:0]
	}

	for _, c := range s.cells {
		for _, creature := range s.system.creatures {
			if collision.CircleCircle(&c.center, c.radius, &creature.Pos, creature.Radius) {
				c.creatures = append(c.creatures, creature)
			}
		}
	}
}

func (s *Collision) CreatureCreature() {
	for _, c := range s.cells {
		for i := 0; i < len(c.creatures); i++ {
			for j := 0; j < len(c.creatures); j++ {
				if collision.CircleCircle(&c.creatures[i].Pos, c.creatures[i].Radius, &c.creatures[j].Pos, c.creatures[j].Radius) {
					c.creatures[i].Collide(c.creatures[j])
					c.creatures[j].Collide(c.creatures[i])
				}
			}
		}
	}
}

func (s *Collision) CreatureEyeCreature() {
	for _, c := range s.cells {
		for _, creature1 := range c.creatures {
			if creature1.Eye == nil {
				continue
			}

			for _, creature2 := range c.creatures {
				d := num.Vec2{X: creature2.Pos.X - creature1.Pos.X, Y: creature2.Pos.Y - creature1.Pos.Y}
				if d.Len() > 50 {
					continue
				}

				if num.Angle(&d, &creature1.Dir) > creature1.Eye.FOV {
					continue
				}

				creature1.Eye.Sees(creature2)
			}
		}
	}
}

func (s *Collision) CreatureBorder() {
	for _, c := range s.system.creatures {
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

func (s *Collision) SetCreatureBorderCB(cb func(*entity.Creature, int)) {
	s.cbCreatureBorder = cb
}
