package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/entity"
)

func TestEntityHandlerUpdatePopulation(t *testing.T) {
	living := func() *entity.Creature {
		return &entity.Creature{
			Alive:  true,
			Energy: 2.0,
			Age:    1.0,
			Consts: entity.Constants{LifeExpectancy: 2.0},
		}
	}

	t.Run("removes dead creatures", func(tt *testing.T) {
		c := living()
		c.Alive = false

		population := []*entity.Creature{
			living(),
			c,
			living(),
		}
		entityHandler := &entity.Handler{}
		populationAfterUpdate := entityHandler.UpdatePopulation(population)

		assert.Equal(tt, 2, len(populationAfterUpdate))
		assert.NotContains(tt, populationAfterUpdate, c)
	})
}
