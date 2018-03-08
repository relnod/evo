package system

import (
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
}

func (s *Entity) Update() {
	for i, c := range s.system.creatures {
		c.Update()

		if !c.Alive {
			s.system.creatures = append(s.system.creatures[:i], s.system.creatures[i+1:]...)
		}

		if c.State == entity.StateBreading {
			c.State = entity.StateAdult
			// log.Printf("Genration: %d, Population: %d\n", c.Generation+1, len(s.system.creatures))
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

	return pos
}
