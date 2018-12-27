package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/system"
	"github.com/relnod/evo/pkg/world"
)

func TestEntitySystemUpdate(t *testing.T) {
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

		world := &world.World{
			Creatures: []*entity.Creature{
				living(),
				c,
				living(),
			},
		}
		entitySystem := system.NewEntity(world)
		entitySystem.Update()

		assert.Equal(tt, 2, len(world.Creatures))
		assert.NotContains(tt, world.Creatures, c)
	})
}
