package system

import (
	"math/rand"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

const (
	ModeRandom = iota
	ModeFixed
)

type Mode int

type Entity struct {
	system *System
	mode   Mode
}

func NewEntity(system *System, mode Mode) *Entity {
	return &Entity{system: system, mode: mode}
}

func (s *Entity) Init() {
	for i := range s.system.creatures {

		var radius float32 = 3.0
		if s.mode == ModeRandom {
			radius = rand.Float32()*rand.Float32()*rand.Float32()*10 + 2.0
		}

		s.system.creatures[i] = entity.NewCreature(s.randomPosition(radius), radius)
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
			c.LastBread = c.Age
			// log.Printf("Genration: %d, Population: %d\n", c.Generation+1, len(s.system.creatures))
			c.Energy -= c.Radius
			for i := 0; i < rand.Intn(2)+1; i++ {
				child := c.GetChild()
				if c.Energy-child.Energy > 0 {
					c.Energy -= child.Energy
					s.system.creatures = append(s.system.creatures, child)
				}
			}
		}
	}
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
