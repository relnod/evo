package system

import (
	"math"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

type Collision struct {
	system *System

	cbCreatureCreature func(*entity.Creature, *entity.Creature)
	cbCreatureEyeFood  func(*entity.Creature, *entity.Food)
	cbCreatureFood     func(*entity.Creature, *entity.Food)
	cbCreatureBorder   func(*entity.Creature, int)

	cells []*cell
}

type cell struct {
	topLeft  num.Vec2
	botRight num.Vec2

	center num.Vec2
	radius float32

	creatures []*entity.Creature
	food      []*entity.Food
}

func NewCollision(s *System, numCells int) *Collision {
	cellsPerRow := int(math.Sqrt(float64(numCells)))

	cellWidth := s.Width / float32(numCells/cellsPerRow)
	cellHeight := s.Height / float32(numCells/cellsPerRow)

	radius := cellWidth
	if cellWidth < cellHeight {
		radius = cellHeight
	}

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
	s.CreatureEyeFood()
	s.CreatureFood()
	s.CreatureBorder()

}

func (s *Collision) ResetCells() {
	for _, cell := range s.cells {
		cell.creatures = cell.creatures[:0]
		cell.food = cell.food[:0]
	}

	for _, c := range s.cells {
		for _, creature := range s.system.creatures {
			if collision.CircleCircle(&c.center, c.radius, &creature.Pos, creature.Radius) {
				c.creatures = append(c.creatures, creature)
			}
		}

		for _, food := range s.system.food {
			if collision.CircleCircle(&c.center, c.radius, &food.Pos, food.Radius) {
				c.food = append(c.food, food)
			}
		}
	}
}

func (s *Collision) CreatureCreature() {
	for _, c := range s.cells {
		for i := 0; i < len(c.creatures); i++ {
			for j := i + 1; j < len(c.creatures); j++ {
				if collision.CircleCircle(&c.creatures[i].Pos, c.creatures[i].Radius, &c.creatures[j].Pos, c.creatures[j].Radius) {
					s.cbCreatureCreature(c.creatures[i], c.creatures[j])
				}
			}
		}
	}
}

func (s *Collision) CreatureEyeFood() {
	for _, c := range s.cells {
		for _, creature := range c.creatures {
			for _, food := range c.food {
				d := num.Vec2{X: food.Pos.X - creature.Pos.X, Y: food.Pos.Y - creature.Pos.Y}
				if d.Len() > 50 {
					continue
				}
				x := creature.Dir.X / d.X
				y := creature.Dir.Y / d.Y
				if x < 0.001 && y < 0.001 {
					s.cbCreatureEyeFood(creature, food)
				}
			}
		}
	}
}

func (s *Collision) CreatureFood() {
	for _, c := range s.cells {
		for _, creature := range c.creatures {
			for _, food := range c.food {
				if collision.CircleCircle(&creature.Pos, creature.Radius, &food.Pos, food.Radius) {
					s.cbCreatureFood(creature, food)
				}
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

func (s *Collision) SetCreatureCreatureCB(cb func(*entity.Creature, *entity.Creature)) {
	s.cbCreatureCreature = cb
}

func (s *Collision) SetCreatureEyeFoodCB(cb func(*entity.Creature, *entity.Food)) {
	s.cbCreatureEyeFood = cb
}

func (s *Collision) SetCreatureFoodCB(cb func(*entity.Creature, *entity.Food)) {
	s.cbCreatureFood = cb
}

func (s *Collision) SetCreatureBorderCB(cb func(*entity.Creature, int)) {
	s.cbCreatureBorder = cb
}
