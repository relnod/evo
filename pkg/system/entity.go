package system

import (
	"math/rand"

	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
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
		var radius = 3.0
		if s.world.Opts.StartMode == world.StartModeRandom {
			radius = rand.Float64()*rand.Float64()*rand.Float64()*10 + 2.0
		}

		s.world.Creatures[i] = entity.NewCreature(s.randomPosition(radius), radius)
	}
}

func (s *Entity) Update() {
	var remove []int
	for i, c := range s.world.Creatures {
		c.Update()

		if !c.Alive {
			remove = append(remove, i)
			continue
		}

		if c.State == entity.StateBreading {
			c.State = entity.StateAdult
			c.LastBread = c.Age
			c.Energy -= c.Radius
			for i := 0; i < rand.Intn(2)+1; i++ {
				child := c.NewChild()
				if c.Energy-child.Energy > 0 {
					c.Energy -= child.Energy
					s.world.Creatures = append(s.world.Creatures, child)
				}
			}
		}
	}

	for _, i := range remove {
		s.world.RemoveEntity(i)
	}
}

func (s *Entity) randomPosition(radius float64) math64.Vec2 {
	pos := math64.Vec2{
		X: rand.Float64()*(float64(s.world.Width)-(2*radius)) + radius,
		Y: rand.Float64()*(float64(s.world.Height)-(2*radius)) + radius,
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
