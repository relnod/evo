package entity_test

import (
	"testing"

	deep "github.com/patrikeh/go-deep"
	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/entity"
)

func TestCreatureUpdate(t *testing.T) {
	living := func() *entity.Creature {
		return &entity.Creature{
			Alive:  true,
			Energy: 2.0,
			Age:    1.0,
			Consts: entity.Constants{LifeExpectancy: 2.0},
		}
	}

	t.Run("lives if energy is greater than 0, and age is less than life expectancy", func(tt *testing.T) {
		c := living()
		assert.Equal(tt, true, c.Alive)
	})

	t.Run("dies if energy is less or equal to 0", func(tt *testing.T) {
		c := living()
		c.Energy = 0.0

		c.Update()
		assert.Equal(tt, false, c.Alive)
	})

	t.Run("dies is age is greater that life expectancy", func(tt *testing.T) {
		c := living()
		c.Age = 2.0
		c.Consts.LifeExpectancy = 1.0

		c.Update()
		assert.Equal(tt, false, c.Alive)
	})
}

func TestCreatureCollide(t *testing.T) {
	t.Run("nothing happens, when both are not moving", func(tt *testing.T) {
		c1 := &entity.Creature{Brain: nil, Alive: true, Radius: 1.0}
		c2 := &entity.Creature{Brain: nil, Alive: true, Radius: 1.0}

		c1.Collide(c2)
		assert.Equal(tt, true, c1.Alive)
		assert.Equal(tt, true, c2.Alive)
	})

	t.Run("c2 dies, if c1 is moving, but c1 not", func(tt *testing.T) {
		c1 := &entity.Creature{Brain: &deep.Neural{}, Alive: true, Radius: 1.0}
		c2 := &entity.Creature{Brain: nil, Alive: true, Radius: 1.0}

		c1.Collide(c2)
		assert.Equal(tt, true, c1.Alive)
		assert.Equal(tt, false, c2.Alive)
	})

	t.Run("c2 dies, if c1 is bigger", func(tt *testing.T) {
		c1 := &entity.Creature{Brain: &deep.Neural{}, Alive: true, Radius: 2.0}
		c2 := &entity.Creature{Brain: &deep.Neural{}, Alive: true, Radius: 1.0}

		c1.Collide(c2)
		assert.Equal(tt, true, c1.Alive)
		assert.Equal(tt, false, c2.Alive)
	})

	t.Run("c2 lives, if c1 is bigger, but is same species", func(tt *testing.T) {
		c1 := &entity.Creature{Brain: &deep.Neural{}, Alive: true, Radius: 1.01}
		c2 := &entity.Creature{Brain: &deep.Neural{}, Alive: true, Radius: 1.0}

		c1.Collide(c2)
		assert.Equal(tt, true, c1.Alive)
		assert.Equal(tt, true, c2.Alive)
	})
}

func TestNewMutaedBrain(t *testing.T) {
	t.Run("doesn't crash, when adding a new input layer", func(tt *testing.T) {
		brain := entity.NewBrain(2)
		entity.NewMutatedBrain(brain, 4)
	})
}
