package entity

import (
	"math/rand"

	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

// Handler implements the evo.EntityHandler.
type Handler struct {
	width  int
	height int

	animalStats *Stats
	plantStats  *Stats

	collectStats bool
}

// NewHandler returns a new entity handler.
func NewHandler(width, height int) *Handler {
	return &Handler{
		width:  width,
		height: height,

		animalStats: &Stats{
			Lifetime:     NewStatistics(),
			Interactions: NewStatistics(),
			Generation:   NewStatistics(),
			DeathBy:      NewStatistics(),
		},
		plantStats: &Stats{
			Lifetime:     NewStatistics(),
			Interactions: NewStatistics(),
			Generation:   NewStatistics(),
			DeathBy:      NewStatistics(),
		},
		collectStats: true,
	}
}

// InitPopulation initializes a population with a given count.
func (h *Handler) InitPopulation(count int) []*Creature {
	var (
		creatures = make([]*Creature, count)
	)

	for i := range creatures {
		radius := rand.Float64()*rand.Float64()*rand.Float64()*10 + 2.0

		creatures[i] = NewCreature(randomPosition(creatures, h.width, h.height, radius), radius)
	}

	return creatures
}

// UpdatePopulation updates all entities.
// Also adds new child entities and removes dead ones.
func (h *Handler) UpdatePopulation(creatures []*Creature) []*Creature {
	var remove []int
	for i, c := range creatures {
		c.Update()

		if !c.Alive {
			if h.collectStats {
				if c.Brain == nil {
					h.plantStats.Interactions.Add(c.Interactions)
					h.plantStats.Lifetime.Add(int(c.Age))
					h.plantStats.Generation.Add(c.Consts.Generation)
					h.plantStats.DeathBy.Add(int(c.DeathBy))
				} else {
					h.animalStats.Interactions.Add(c.Interactions)
					h.animalStats.Lifetime.Add(int(c.Age))
					h.animalStats.Generation.Add(c.Consts.Generation)
					h.animalStats.DeathBy.Add(int(c.DeathBy))
				}
			}
			remove = append(remove, i)
			continue
		}

		if c.State == StateBreading {
			c.State = StateAdult
			c.LastBread = c.Age
			c.Energy -= c.Radius
			for i := 0; i < rand.Intn(int(1/(c.Radius*c.Radius*c.Radius*c.Radius)*100)+1)+1; i++ {
				child := c.NewChild()
				if c.Energy-child.Energy > 0 {
					c.Energy -= child.Energy
					creatures = append(creatures, child)
				}
			}
		}
	}

	for _, i := range remove {
		creatures = RemoveEntity(creatures, i)
	}

	return creatures
}

func (h *Handler) AnimalStats() *Stats {
	return h.animalStats
}

func (h *Handler) PlantStats() *Stats {
	return h.plantStats
}

// RemoveEntity removes an entity at a given index.
func RemoveEntity(creatures []*Creature, i int) []*Creature {
	if i+1 >= len(creatures) {
		return creatures[:i]
	}

	return append(creatures[:i], creatures[i+1:]...)
}

func randomPosition(creatures []*Creature, width, height int, radius float64) math64.Vec2 {
	pos := math64.Vec2{
		X: rand.Float64()*(float64(width)-(2*radius)) + radius,
		Y: rand.Float64()*(float64(height)-(2*radius)) + radius,
	}

	for _, creature := range creatures {
		if creature == nil {
			continue
		}

		if collision.CircleCircle(&creature.Pos, creature.Radius, &pos, radius) {
			return randomPosition(creatures, width, height, radius)
		}
	}

	return pos
}
