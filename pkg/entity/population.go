package entity

import (
	"math/rand"

	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

// InitPopulation initializes a population with a given count and a world size.
func InitPopulation(count, width, height int) []*Creature {
	creatures := make([]*Creature, count)

	for i := range creatures {
		radius := rand.Float64()*rand.Float64()*rand.Float64()*10 + 2.0

		creatures[i] = NewCreature(randomPosition(creatures, width, height, radius), radius)
	}

	return creatures
}

// randomPosition returns a new random position in the world that is free.
// A position is free, if it won't collide with any other creature.
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

// PopulationUpdater implements the evo.EntityUpdater.
type PopulationUpdater struct {
	animalStats *DeathStats
	plantStats  *DeathStats

	collectStats bool
}

// NewPopulationUpdater returns a new population updater.
func NewPopulationUpdater() *PopulationUpdater {
	return &PopulationUpdater{
		animalStats:  &DeathStats{},
		plantStats:   &DeathStats{},
		collectStats: true,
	}
}

// UpdatePopulation updates all entities.
// Also adds new child entities and removes dead ones.
func (p *PopulationUpdater) UpdatePopulation(creatures []*Creature) []*Creature {
	var remove []int
	for i, c := range creatures {
		c.Update()

		if !c.Alive {
			if p.collectStats {
				if c.Brain == nil {
					p.plantStats.Add(c)
				} else {
					p.animalStats.Add(c)
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

// AnimalStats returns the death stats for animals.
func (p *PopulationUpdater) AnimalStats() *DeathStats {
	return p.animalStats
}

// PlantStats returns the death stats for plants.
func (p *PopulationUpdater) PlantStats() *DeathStats {
	return p.plantStats
}

// ClearStats resets the internal stats counters.
func (p *PopulationUpdater) ClearStats() {
	p.plantStats.Clear()
	p.animalStats.Clear()
}

// RemoveEntity removes an entity at a given index.
func RemoveEntity(creatures []*Creature, i int) []*Creature {
	if i+1 >= len(creatures) {
		return creatures[:i]
	}

	return append(creatures[:i], creatures[i+1:]...)
}
