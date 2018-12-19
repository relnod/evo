package system

import (
	"math/rand"

	"github.com/relnod/evo/pkg/collision"
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math32"
	"github.com/relnod/evo/pkg/world"
)

type Entity struct {
	world *world.World
}

func NewEntity(world *world.World) *Entity {
	return &Entity{world: world}
}

func (s *Entity) Init() {
	for i := range s.world.Creatures {
		var radius float32 = 3.0
		if s.world.StartMode == world.StartModeRandom {
			radius = rand.Float32()*rand.Float32()*rand.Float32()*10 + 2.0
		}

		s.world.Creatures[i] = entity.NewCreature(s.randomPosition(radius), radius)
	}
}

func (s *Entity) Update() {
	for i, c := range s.world.Creatures {
		c.Update()

		if !c.Alive {
			s.world.RemoveCreature(i)
			continue
		}

		if c.State == entity.StateBreading {
			c.State = entity.StateAdult
			c.LastBread = c.Age
			// log.Printf("Genration: %d, Population: %d\n", c.Generation+1, len(s.world.Creatures))
			c.Energy -= c.Radius
			for i := 0; i < rand.Intn(2)+1; i++ {
				child := c.GetChild()
				if c.Energy-child.Energy > 0 {
					c.Energy -= child.Energy
					s.world.Creatures = append(s.world.Creatures, child)
				}
			}
		}
	}
}

func (s *Entity) randomPosition(radius float32) math32.Vec2 {
	pos := math32.Vec2{
		X: rand.Float32()*(s.world.Width-(2*radius)) + radius,
		Y: rand.Float32()*(s.world.Height-(2*radius)) + radius,
	}

	for _, creature := range s.world.Creatures {
		if creature == nil {
			continue
		}

		if collision.CircleCircle(&creature.Pos, creature.Radius, &pos, radius) {
			return s.randomPosition(radius)
		}
	}

	return pos
}
