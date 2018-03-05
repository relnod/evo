package system

import (
	"log"
	"math/rand"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

type Entity struct {
	system *System
}

func NewEntity(system *System) *Entity {
	return &Entity{system: system}
}

func (s *Entity) Init() {
	for i := range s.system.creatures {
		s.system.creatures[i] = s.NewCreature()
	}

	for i := range s.system.food {
		s.system.food[i] = s.NewFood()
	}
}

func (s *Entity) Update() {
	for i, c := range s.system.creatures {
		c.Update()

		if !c.Alive {
			s.system.creatures = append(s.system.creatures[:i], s.system.creatures[i+1:]...)
		}

		if c.Saturation > 10 && (c.Age-c.LastBread) > 40 {
			log.Printf("Genration: %d, Population: %d\n", c.Generation+1, len(s.system.creatures))
			c.LastBread = c.Age

			child1 := c.GetChild()
			s.system.creatures = append(s.system.creatures, child1)

			child2 := c.GetChild()
			s.system.creatures = append(s.system.creatures, child2)
		}
	}
}

func (s *Entity) NewCreature() *entity.Creature {
	c := entity.NewCreature(s.randomPosition(8.0))

	return c
}

func (s *Entity) NewFood() *entity.Food {
	food := &entity.Food{Pos: s.randomPosition(4.0), Radius: 4.0}
	s.system.food = append(s.system.food, food)

	return food
}

func (s *Entity) ResetFood(food *entity.Food) {
	food.Pos = s.randomPosition(food.Radius)
}

func (s *Entity) randomPosition(radius float32) num.Vec2 {
	pos := num.Vec2{
		X: rand.Float32()*(s.system.Width-(2*radius)) + radius,
		Y: rand.Float32()*(s.system.Height-(2*radius)) + radius,
	}

	for _, creature := range s.system.creatures {
		if creature == nil {
			continue
		}

		if collision.CircleCircle(&creature.Pos, creature.Radius, &pos, radius) {
			return s.randomPosition(radius)
		}
	}

	for _, food := range s.system.food {
		if food == nil {
			continue
		}

		if collision.CircleCircle(&food.Pos, food.Radius, &pos, radius) {
			return s.randomPosition(radius)
		}
	}

	return pos
}
