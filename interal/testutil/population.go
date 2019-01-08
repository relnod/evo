package testutil

import (
	"math/rand"

	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64"
)

// Population returns a new population with a given size.
// Each population for a size is deterministic, so it can be used for
// benchmarks.
func Population(size int) []*entity.Creature {
	rand.Seed(123734)
	var population []*entity.Creature
	for i := 0; i < size; i++ {
		population = append(population, entity.NewCreature(math64.Vec2{X: rand.Float64() * 10, Y: rand.Float64() * 10}, rand.Float64()*2))
	}
	return population
}
