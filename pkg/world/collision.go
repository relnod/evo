package world

import (
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64/collision"
)

// Collision defines an interface for a 2D collision, that can be resolved.
type Collision interface {
	// Resolve resolves the collision.
	Resolve()
}

// CollisionDetector detects collisions in the world.
type CollisionDetector interface {
	DetectCollisions(creatures []*entity.Creature) []Collision
}

// ResolveAllCollisions resolves all given collisions.
func ResolveAllCollisions(collisions []Collision) {
	for _, c := range collisions {
		c.Resolve()
	}
}

type creatureCreatureCollision struct {
	creature1 *entity.Creature
	creature2 *entity.Creature
}

func (c *creatureCreatureCollision) Resolve() {
	c.creature1.Collide(c.creature2)
}

type eyeCreatureCollision struct {
	eye      *entity.Eye
	creature *entity.Creature
}

func (c *eyeCreatureCollision) Resolve() {
	c.eye.Sees(c.creature)
}

type creatureBorderCollision struct {
	creature *entity.Creature
	border   int

	// TODO: it's pretty stupid that we have to store the world size here. And
	// quite expensive...
	width  int
	height int
}

func (c *creatureBorderCollision) Resolve() {
	switch c.border {
	case collision.LEFT:
		c.creature.Pos.X += float64(c.width)
	case collision.RIGHT:
		c.creature.Pos.X -= float64(c.width)
	case collision.TOP:
		c.creature.Pos.Y += float64(c.height)
	case collision.BOT:
		c.creature.Pos.Y -= float64(c.height)
	}
}
